package flags

func applyEnvironmentVariableFlags() {
	getEnvironmentVariableOrFlag("BIND_ADDRESS_IPv4", BindAddressIpv4)

	getEnvironmentVariableOrFlag("VAULT_ADDR", VaultAddress)
	getEnvironmentVariableOrFlag("VAULT_CREDENTIAL", VaultCredential)
	getEnvironmentVariableOrFlag("VAULT_AUTHENTICATION_TYPE", VaultAuthenticationType)
	getEnvironmentVariableOrFlag("VAULT_NAMESPACE", VaultNamespace)
	getEnvironmentVariableOrFlag("VAULT_SKIP_TLS_VERIFY", VaultSkipTlsVerify)

	getEnvironmentVariableOrFlag("CLIENT_SETTINGS_VAULT_MOUNT_PATH", ClientSettingsVaultMountPath)
	getEnvironmentVariableOrFlag("CLIENT_SETTINGS_VAULT_PATH", ClientSettingsVaultPath)
	getEnvironmentVariableOrFlag("CLIENT_SETTINGS_PROVIDER_REFRESH_INTERVAL", ClientSettingsProviderRefreshInterval)
	getEnvironmentVariableOrFlag("CLIENT_SETTINGS_API_KEYS", ClientSettingsApiKeys)
}
