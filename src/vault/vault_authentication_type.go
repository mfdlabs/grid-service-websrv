package vault

import (
	"fmt"
	"strings"
)

type VaultAuthenticationType int

const (
	Token VaultAuthenticationType = iota
	AppRole
	LDAP
)

func ToVaultAuthenticationType(s string) (VaultAuthenticationType, error) {
	switch strings.ToLower(s) {
	case "token":
		return Token, nil
	case "approle":
		return AppRole, nil
	case "ldap":
		return LDAP, nil
	default:
		return Token, fmt.Errorf("unknown VaultAuthenticationType %s", s)
	}
}

func (v VaultAuthenticationType) String() string {
	switch v {
	case Token:
		return "Token"
	case AppRole:
		return "AppRole"
	case LDAP:
		return "LDAP"
	default:
		return "Unknown"
	}
}
