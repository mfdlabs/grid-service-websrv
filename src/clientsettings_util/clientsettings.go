package clientsettingsutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/hashicorp/vault/api"
	"github.com/mfdlabs/grid-service-websrv/metrics"
)

// ClientSettingsProvider provides the client settings in a thread-safe way.
type ClientSettingsProvider struct {
	rwMutex *sync.RWMutex

	cachedGroupSettings                    map[string]map[string]interface{}
	cachedGroupDepends                     map[string][]string
	currentRefreshedGroups                 map[string]bool
	allowedGroupsFromClientSettingsService map[string]bool

	vaultClient *api.Client
	mountPath   string
	path        string
}

// NewClientSettingsProvider creates a new ClientSettingsProvider.
func NewClientSettingsProvider(vaultClient *api.Client, mountPath, path string, refreshInterval int) *ClientSettingsProvider {
	csp := &ClientSettingsProvider{
		rwMutex:                                &sync.RWMutex{},
		cachedGroupSettings:                    make(map[string]map[string]interface{}),
		cachedGroupDepends:                     make(map[string][]string),
		currentRefreshedGroups:                 make(map[string]bool),
		allowedGroupsFromClientSettingsService: make(map[string]bool),
		vaultClient:                            vaultClient,
		mountPath:                              mountPath,
		path:                                   path,
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

func mergeArrays(a, b []string) []string {
	for _, v := range b {
		if stringInSlice(v, a) {
			continue
		}

		a = append(a, v)
	}

	return a
}

func (csp *ClientSettingsProvider) resolveWithDependencies(group string) (map[string]interface{}, []string, bool) {
	metrics.GroupReads.WithLabelValues(group).Inc()

	var (
		cs  map[string]interface{}
		dps []string
		ok  bool
	)
	if cs, ok = csp.cachedGroupSettings[group]; !ok {
		return nil, nil, false
	}

	if depends, ok := csp.cachedGroupDepends[group]; ok {
		oldCs := cs
		cs = make(map[string]interface{})
		dps = csp.cachedGroupDepends[group]

		for _, depend := range depends {
			if dependSettings, dependDependencies, ok := csp.resolveWithDependencies(depend); ok {
				cs = mergeMaps(cs, dependSettings)
				dps = mergeArrays(dps, dependDependencies)
			}
		}

		cs = mergeMaps(cs, oldCs) // Overwrite the old settings with the new ones.
	}

	return cs, dps, true
}

// GetGroup returns the client settings.
func (csp *ClientSettingsProvider) GetGroup(group string) (map[string]interface{}, []string, bool, bool) {
	csp.rwMutex.RLock()
	defer csp.rwMutex.RUnlock()

	cs, depends, ok := csp.resolveWithDependencies(group)
	if !ok {
		return nil, nil, false, false
	}

	return cs, depends, csp.allowedGroupsFromClientSettingsService[group], true
}

// Get returns a client setting.
func (csp *ClientSettingsProvider) Get(group, name string) (interface{}, bool) {
	csp.rwMutex.RLock()
	defer csp.rwMutex.RUnlock()

	cs, _, ok := csp.resolveWithDependencies(group)
	if !ok {
		return nil, false
	}

	value, ok := cs[name]
	return value, ok
}

// Set sets the client settings.
func (csp *ClientSettingsProvider) Set(group, name string, value interface{}) (bool, error) {
	csp.rwMutex.Lock()

	existsingCachedGroupSettings, ok := csp.cachedGroupSettings[group]
	if !ok {
		existsingCachedGroupSettings = make(map[string]interface{})
		csp.allowedGroupsFromClientSettingsService[group] = true
	} else {
		if existsingCachedGroupSettings[name] == value {
			csp.rwMutex.Unlock()
			return false, nil
		}
	}

	existsingCachedGroupSettings[name] = value

	csp.rwMutex.Unlock()

	return true, csp.Import(group, existsingCachedGroupSettings, csp.cachedGroupDepends[group], csp.allowedGroupsFromClientSettingsService[group])
}

// Import imports the client settings to Vault.
func (csp *ClientSettingsProvider) Import(group string, data map[string]interface{}, depends []string, isAllowedFromClientSettingsService bool) error {
	csp.rwMutex.Lock()
	defer csp.rwMutex.Unlock()

	metrics.GroupWrites.WithLabelValues(group).Inc()

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
			case int32, int64, uint32, uint64, float32, float64, int, uint, json.Number:
				metadata[key] = "int"
			}
		}

		switch v := value.(type) {
		case bool:
			vaultData[key] = strconv.FormatBool(v)
		case int:
			vaultData[key] = strconv.FormatInt(int64(v), 10)
		case int32:
			vaultData[key] = strconv.FormatInt(int64(v), 10)
		case int64:
			vaultData[key] = strconv.FormatInt(v, 10)
		case uint:
			vaultData[key] = strconv.FormatUint(uint64(v), 10)
		case uint32:
			vaultData[key] = strconv.FormatUint(uint64(v), 10)
		case uint64:
			vaultData[key] = strconv.FormatUint(v, 10)
		case float32:
			vaultData[key] = strconv.FormatFloat(float64(v), 'f', -1, 32)
		case float64:
			vaultData[key] = strconv.FormatFloat(v, 'f', -1, 64)
		case json.Number:
			vaultData[key] = v.String()
		default:
			vaultData[key] = v
		}
	}

	if len(depends) > 0 {
		for _, depend := range depends {
			depend = strings.TrimSpace(depend)

			if _, ok := csp.cachedGroupSettings[depend]; !ok {
				return fmt.Errorf("unknown dependency %s", depend)
			}
		}

		csp.cachedGroupDepends[group] = depends
		metadata["$dependencies"] = strings.Join(depends, ",")
	}

	metadata["$allowed"] = strconv.FormatBool(isAllowedFromClientSettingsService)

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

// ImportReference imports the client settings to Vault.
func (csp *ClientSettingsProvider) ImportReference(group, reference string, isAllowedFromClientSettingsService bool) error {
	csp.rwMutex.Lock()
	defer csp.rwMutex.Unlock()

	_, err := csp.vaultClient.KVv2(csp.mountPath).Put(
		context.Background(),
		csp.path+"/"+group,
		map[string]interface{}{},
	)
	if err != nil {
		return err
	}

	err = csp.vaultClient.KVv2(csp.mountPath).PutMetadata(
		context.Background(),
		csp.path+"/"+group,
		api.KVMetadataPutInput{
			CustomMetadata: map[string]interface{}{
				"$ref":     reference,
				"$allowed": strconv.FormatBool(isAllowedFromClientSettingsService),
			},
		},
	)
	if err != nil {
		return err
	}

	csp.allowedGroupsFromClientSettingsService[group] = isAllowedFromClientSettingsService

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
		// Check if allowed from client setting service
		if allowed, ok := metadata.CustomMetadata["$allowed"]; ok {
			if allowed.(string) != "true" {
				csp.allowedGroupsFromClientSettingsService[group] = false
			} else {
				csp.allowedGroupsFromClientSettingsService[group] = true
			}
		} else {
			csp.allowedGroupsFromClientSettingsService[group] = true
		}

		// Check for $ref
		if ref, ok := metadata.CustomMetadata["$ref"]; ok {
			csp.getSettingsForGroupAndCache(ref.(string))
			csp.cachedGroupSettings[group], _, _ = csp.resolveWithDependencies(ref.(string))
			csp.currentRefreshedGroups[group] = true

			return csp.cachedGroupSettings[group]
		}

		// Check for $dependencies
		if dependencies, ok := metadata.CustomMetadata["$dependencies"]; ok {
			split := strings.Split(dependencies.(string), ",")
			depends := make([]string, 0)

			for _, dependency := range split {
				dependency = strings.TrimSpace(dependency)

				if dependency == "" {
					continue
				}

				depends = append(depends, dependency)

				csp.getSettingsForGroupAndCache(dependency) // ahead of time
			}

			csp.cachedGroupDepends[group] = depends
		}
	} else {
		csp.allowedGroupsFromClientSettingsService[group] = true
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.EqualFold(b, a) {
			return true
		}
	}

	return false
}

func (csp *ClientSettingsProvider) checkForDeletedGroups(vaultGroups []string) {
	for group := range csp.cachedGroupSettings {
		if !stringInSlice(group, vaultGroups) {
			delete(csp.cachedGroupSettings, group)
			delete(csp.cachedGroupDepends, group)
			delete(csp.allowedGroupsFromClientSettingsService, group)
		}
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

	stringGroups := make([]string, len(groups.Data["keys"].([]interface{})))
	for _, group := range groups.Data["keys"].([]interface{}) {
		metrics.SettingsRefreshes.WithLabelValues(group.(string)).Inc()
		csp.getSettingsForGroupAndCache(group.(string))
		stringGroups = append(stringGroups, group.(string))
	}

	csp.checkForDeletedGroups(stringGroups)

	return nil
}
