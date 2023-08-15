package clientsettingsutil

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/hashicorp/vault/api"
)

// ClientSettingsProvider provides the client settings in a thread-safe way.
type ClientSettingsProvider struct {
	rwMutex *sync.RWMutex

	cachedGroupSettings    map[string]map[string]interface{}
	cachedGroupDepends     map[string][]string
	currentRefreshedGroups map[string]bool

	vaultClient *api.Client
	mountPath   string
	path        string
}

// NewClientSettingsProvider creates a new ClientSettingsProvider.
func NewClientSettingsProvider(vaultClient *api.Client, mountPath, path string, refreshInterval int) *ClientSettingsProvider {
	csp := &ClientSettingsProvider{
		rwMutex:             &sync.RWMutex{},
		cachedGroupSettings: make(map[string]map[string]interface{}),
		cachedGroupDepends:  make(map[string][]string),
		vaultClient:         vaultClient,
		mountPath:           mountPath,
		path:                path,
	}

	go csp.updateThread(refreshInterval)

	return csp
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	for k, v := range b {
		a[k] = v
	}

	return a
}

func (csp *ClientSettingsProvider) resolveWithDependencies(group string) (map[string]interface{}, bool) {
	var (
		cs map[string]interface{}
		ok bool
	)
	if cs, ok = csp.cachedGroupSettings[group]; !ok {
		return nil, false
	}

	if depends, ok := csp.cachedGroupDepends[group]; ok {
		oldCs := cs
		cs = make(map[string]interface{})

		for _, depend := range depends {
			if dependSettings, ok := csp.resolveWithDependencies(depend); ok {
				cs = mergeMaps(cs, dependSettings)
			}
		}

		cs = mergeMaps(cs, oldCs) // Overwrite the old settings with the new ones.
	}

	return cs, true
}

// Get returns the client settings.
func (csp *ClientSettingsProvider) Get(group string) (map[string]interface{}, bool) {
	csp.rwMutex.RLock()
	defer csp.rwMutex.RUnlock()

	cs, ok := csp.resolveWithDependencies(group)

	return cs, ok
}

// Set sets the client settings.
func (csp *ClientSettingsProvider) Set(group, name string, value interface{}) {
	csp.rwMutex.Lock()
	defer csp.rwMutex.Unlock()

	csp.cachedGroupSettings[group][name] = value
}

// Import imports the client settings to Vault.
func (csp *ClientSettingsProvider) Import(group string, data map[string]interface{}, depends []string) error {
	csp.rwMutex.Lock()
	defer csp.rwMutex.Unlock()

	// It is in the format of rbx-client-settings.
	// We need to convert it to the format that Vault expects.
	// All settings that aren't prefixed and aren't already strings will have their
	// types put in the metadata.

	metadata := make(map[string]interface{})
	vaultData := make(map[string]interface{})

	for key, value := range data {
		if !regexp.MustCompile(`^([DS])?F(Flag|Int|String)`).MatchString(key) {
			switch value.(type) {
			case bool:
				metadata[key] = "bool"
			case float64:
				metadata[key] = "int"
			}
		}

		switch v := value.(type) {
		case bool:
			vaultData[key] = strconv.FormatBool(v)
		case float64:
			vaultData[key] = strconv.FormatFloat(v, 'f', -1, 64)
		default:
			vaultData[key] = v
		}
	}

	if len(depends) > 0 {
		for _, depend := range depends {
			if _, ok := csp.cachedGroupSettings[depend]; !ok {
				return fmt.Errorf("unknown dependency %s", depend)
			}
		}

		csp.cachedGroupDepends[group] = depends
		metadata["$dependencies"] = strings.Join(depends, ",")
	}

	_, err := csp.vaultClient.KVv2(csp.mountPath).Put(
		context.Background(),
		csp.path+"/"+group,
		vaultData,
	)
	if err != nil {
		return err
	}

	err = csp.vaultClient.KVv2(csp.mountPath).PutMetadata(
		context.Background(),
		csp.path+"/"+group,
		api.KVMetadataPutInput{
			CustomMetadata: metadata,
		},
	)
	if err != nil {
		return err
	}

	csp.cachedGroupSettings[group] = data

	return nil
}

// Refresh refreshes the client settings from Vault.
func (csp *ClientSettingsProvider) Refresh() error {
	return csp.doRefresh()
}

func parseData(data, metadata map[string]interface{}) map[string]interface{} {
	// rbx-client-settings are like this:
	// FFlag is a bool
	// FInt is an int
	// FString is a string
	//
	// Settings prefixed with no prefix (such as FFlagTest) are static settings.
	// Settings prefixed with D (such as DFFlagTest) are dynamic settings.
	// Settings prefixed with S (such as SFFlagTest) are server synchronized settings.
	//
	// Vault secrets values are all strings, so we need to convert them to the correct type.
	// We also need to remove the prefix.

	settings := make(map[string]interface{})
	for key, value := range data {
		// Check if the metadata has the key, if it does then that's the type.
		if metadata != nil {
			if s, ok := metadata[key]; ok {
				switch strings.ToLower(s.(string)) {
				case "bool":
					b, err := strconv.ParseBool(value.(string))
					if err != nil {
						b = false
					}

					settings[key] = b
				case "int":
					i, err := strconv.Atoi(value.(string))
					if err != nil {
						i = 0
					}

					settings[key] = i
				case "string":
					settings[key] = value.(string)
				}

				continue
			}
		}

		// Check the regex for the prefix.
		if !regexp.MustCompile(`^([DS])?F(Flag|Int|String)`).MatchString(key) {
			settings[key] = value
			continue
		}

		// If it doesn't start with F, remove the prefix.
		k := key
		if key[0] != 'F' {
			k = key[1:]
		}

		switch k[1] {
		case 'F':
			b, err := strconv.ParseBool(value.(string))
			if err != nil {
				b = false
			}

			settings[key] = b
		case 'I':
			i, err := strconv.ParseInt(value.(string), 10, 64)
			if err != nil {
				i = 0
			}

			settings[key] = i
		case 'S':
			settings[key] = value.(string)
		}
	}

	return settings
}

func (csp *ClientSettingsProvider) getSettingsForGroupAndCache(group string) map[string]interface{} {
	if exists := csp.currentRefreshedGroups[group]; exists {
		return csp.cachedGroupSettings[group]
	}

	metadata, err := csp.vaultClient.KVv2(csp.mountPath).GetMetadata(context.Background(), csp.path+"/"+group)
	if err != nil {
		glog.Warningf("Skipping group %s: %s", group, err.Error())

		return nil // DON'T PANIC
	}

	if metadata != nil {
		// Check for $ref
		if ref, ok := metadata.CustomMetadata["$ref"]; ok {
			csp.getSettingsForGroupAndCache(ref.(string))
			csp.cachedGroupSettings[group], _ = csp.resolveWithDependencies(ref.(string))
			csp.currentRefreshedGroups[group] = true

			return csp.cachedGroupSettings[group]
		}

		// Check for $dependencies
		if dependencies, ok := metadata.CustomMetadata["$dependencies"]; ok {
			split := strings.Split(dependencies.(string), ",")
			depends := make([]string, 0)

			for _, dependency := range split {
				if dependency == "" {
					continue
				}

				depends = append(depends, dependency)

				csp.getSettingsForGroupAndCache(dependency) // ahead of time
			}

			csp.cachedGroupDepends[group] = depends
		}
	}

	secret, err := csp.vaultClient.KVv2(csp.mountPath).Get(context.Background(), csp.path+"/"+group)
	if err != nil {
		glog.Warningf("Skipping group %s: %s", group, err.Error())

		return nil // DON'T PANIC
	}

	csp.cachedGroupSettings[group] = parseData(secret.Data, metadata.CustomMetadata)
	csp.currentRefreshedGroups[group] = true

	return csp.cachedGroupSettings[group]
}

func (csp *ClientSettingsProvider) updateThread(refreshInterval int) {
	for {
		if err := csp.doRefresh(); err != nil {
			glog.Errorf("Failed to refresh settings: %s", err.Error())
		}

		time.Sleep(time.Duration(refreshInterval) * time.Second)
	}
}

func (csp *ClientSettingsProvider) doRefresh() error {
	csp.rwMutex.Lock()
	defer csp.rwMutex.Unlock()

	csp.currentRefreshedGroups = make(map[string]bool)

	groups, err := csp.vaultClient.Logical().List(csp.mountPath + "/metadata/" + csp.path)
	if err != nil {
		return err
	}

	if groups == nil || groups.Data == nil {
		return errors.New("no groups found")
	}

	for _, group := range groups.Data["keys"].([]interface{}) {
		csp.getSettingsForGroupAndCache(group.(string))
	}

	return nil
}
