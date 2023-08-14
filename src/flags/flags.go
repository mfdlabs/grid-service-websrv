package flags

import (
	"flag"
)

var (
	// BindAddressIpv4 is the address to bind the prometheus metrics server to.
	BindAddressIpv4 = flag.String("bind-address-ipv4", ":8080", "Address to bind the HTTP server to.")

	// HelpFlag prints the usage.
	HelpFlag = flag.Bool("help", false, "Print usage.")

	///////////////////////////
	// Vault configuration

	// VaultAddress is the address of the Vault server.
	VaultAddress = flag.String("vault-address", "http://localhost:8200", "Address of the Vault server. (environment variable: VAULT_ADDR)")

	// VaultToken is the token to use to authenticate with Vault.
	VaultCredential = flag.String("vault-credential", "", "Token to use to authenticate with Vault. (environment variable: VAULT_CREDENTIAL)")

	// VaultAuthenticationType is the authentication type to use to authenticate with Vault.
	VaultAuthenticationType = flag.String("vault-authentication-type", "token", "Authentication type to use to authenticate with Vault. (environment variable: VAULT_AUTHENTICATION_TYPE)")

	// VaultNamespace is the namespace to use to authenticate with Vault.
	VaultNamespace = flag.String("vault-namespace", "", "Namespace to use to authenticate with Vault. (environment variable: VAULT_NAMESPACE)")

	// VaultSkipTlsVerify skips TLS verification for Vault.
	VaultSkipTlsVerify = flag.Bool("vault-skip-tls-verify", false, "Skip TLS verification for Vault. (environment variable: VAULT_SKIP_TLS_VERIFY)")

	///////////////////////////
	// Client Settings configuration

	// ClientSettingsVaultMountPath is the mount path of the client settings in Vault.
	ClientSettingsVaultMountPath = flag.String("client-settings-vault-mount-path", "client-settings", "Mount path of the client settings in Vault. (environment variable: CLIENT_SETTINGS_VAULT_MOUNT_PATH)")

	// ClientSettingsVaultPath is the path of the client settings in Vault.
	ClientSettingsVaultPath = flag.String("client-settings-vault-path", "", "Path of the client settings in Vault. (environment variable: CLIENT_SETTINGS_VAULT_PATH)")

	// ClientSettingsProviderRefreshInterval is the interval at which the client settings provider refreshes the client settings.
	ClientSettingsProviderRefreshInterval = flag.Int("client-settings-provider-refresh-interval", 30, "Interval at which the client settings provider refreshes the client settings. (environment variable: CLIENT_SETTINGS_PROVIDER_REFRESH_INTERVAL)")

	// ClientSettingsApiKeys is the comma-separated list of API keys to use to authenticate with the client settings provider.
	ClientSettingsApiKeys = flag.String("client-settings-api-keys", "", "Comma-separated list of API keys to use to authenticate with the client settings provider. (environment variable: CLIENT_SETTINGS_API_KEYS)")
)

const FlagsUsageString string = `
	[-h|--help] [--bind-address-ipv4[=:8080]] 
	[--vault-address=http://localhost:8200] [--vault-credential[=]] [--vault-authentication-type[=token]]
	[--vault-namespace[=]] [--vault-skip-tls-verify[=false]]`
