package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/auth/ldap"
	"github.com/mfdlabs/grid-service-websrv/flags"
)

type VaultClient struct {
	Client *api.Client
	t      *time.Ticker
	quit   chan bool
}

var globalClient *VaultClient
var createOnce sync.Once

func initializeAppRoleAuth(roleId, secretId string) (*approle.AppRoleAuth, error) {
	approleSecretId := &approle.SecretID{
		FromString: secretId,
	}

	return approle.NewAppRoleAuth(
		roleId,
		approleSecretId,
	)
}

func initializeLdapAuth(username, password string) (*ldap.LDAPAuth, error) {
	ldapPassword := &ldap.Password{
		FromString: password,
	}

	return ldap.NewLDAPAuth(
		username,
		ldapPassword,
	)
}

func newVaultClient(ctx context.Context, address string, authType vaultAuthenticationType, credential string) error {
	var err error = nil

	createOnce.Do(func() {
		globalClient = &VaultClient{}
		config := api.DefaultConfig()
		config.Address = address
		config.ConfigureTLS(&api.TLSConfig{
			Insecure: *flags.VaultSkipTlsVerify,
		})

		globalClient.Client, err = api.NewClient(config)
		if err != nil {
			return
		}

		if *flags.VaultNamespace != "" {
			globalClient.Client.SetNamespace(*flags.VaultNamespace)
		}

		switch authType {
		case vaultAuthenticationTypeToken:
			globalClient.Client.SetToken(credential)
		case vaultAuthenticationTypeAppRole:
			split := strings.Split(credential, ":")

			if len(split) != 2 {
				err = errors.New("the vaultCredential key when using AppRole authType is expected in the following format: 'roleId:secretId'")
				return
			}

			roleId := split[0]
			secretId := split[1]

			var approleAuth *approle.AppRoleAuth

			approleAuth, err = initializeAppRoleAuth(roleId, secretId)
			if err != nil {
				return
			}

			_, err = globalClient.Client.Auth().Login(ctx, approleAuth)
		case vaultAuthenticationTypeLDAP:
			split := strings.Split(credential, ":")

			if len(split) != 2 {
				err = errors.New("the vaultCredential key when using LDAP authType is expected in the following format: 'username:password'")
				return
			}

			username := split[0]
			password := split[1]

			var ldapAuth *ldap.LDAPAuth

			ldapAuth, err = initializeLdapAuth(username, password)
			if err != nil {
				return
			}

			_, err = globalClient.Client.Auth().Login(ctx, ldapAuth)
		default:
			err = fmt.Errorf("unknown vault authentication type %s", authType.string())
			return
		}

		globalClient.t = time.NewTicker(30 * time.Minute)
		globalClient.quit = make(chan bool)
		go func() {
			for {
				select {
				case <-globalClient.t.C:
					globalClient.Client.Auth().Token().RenewSelf(int(time.Hour))
				case <-globalClient.quit:
					globalClient.t.Stop()
					return
				}
			}
		}()
	})

	return err
}

func (c *VaultClient) StopRefreshingToken() {
	close(c.quit)
}

func GetGlobalVaultClient(ctx context.Context) (*VaultClient, error) {
	if globalClient != nil {
		return globalClient, nil
	}

	if *flags.VaultAddress == "" {
		return nil, errors.New("cannot setup vault client because vault address is empty")
	}
	if *flags.VaultAuthenticationType == "" {
		return nil, errors.New("cannot setup vault client because vault authentication type is empty")
	}
	if *flags.VaultCredential == "" {
		return nil, errors.New("cannot setup vault client because vault credential is empty")
	}

	authType, err := toVaultAuthenticationType(*flags.VaultAuthenticationType)
	if err != nil {
		return nil, err
	}

	err = newVaultClient(ctx, *flags.VaultAddress, authType, *flags.VaultCredential)

	return globalClient, err
}
