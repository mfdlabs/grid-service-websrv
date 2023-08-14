package vault

import (
	"fmt"
	"strings"
)

type vaultAuthenticationType int

const (
	vaultAuthenticationTypeToken vaultAuthenticationType = iota
	vaultAuthenticationTypeAppRole
	vaultAuthenticationTypeLDAP
)

func toVaultAuthenticationType(s string) (vaultAuthenticationType, error) {
	switch strings.ToLower(s) {
	case "token":
		return vaultAuthenticationTypeToken, nil
	case "approle":
		return vaultAuthenticationTypeAppRole, nil
	case "ldap":
		return vaultAuthenticationTypeLDAP, nil
	default:
		return vaultAuthenticationTypeToken, fmt.Errorf("unknown VaultAuthenticationType %s", s)
	}
}

func (v vaultAuthenticationType) string() string {
	switch v {
	case vaultAuthenticationTypeToken:
		return "Token"
	case vaultAuthenticationTypeAppRole:
		return "AppRole"
	case vaultAuthenticationTypeLDAP:
		return "LDAP"
	default:
		return "Unknown"
	}
}
