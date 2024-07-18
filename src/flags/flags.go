package flags

import (
	"flag"
	"time"
)

var (
	// BindAddressIpv4 is the address to bind the prometheus metrics server to.
	BindAddressIpv4 = flag.String("bind-address-ipv4", ":8080", "Address to bind the HTTP server to.")

	// HelpFlag prints the usage.
	HelpFlag = flag.Bool("help", false, "Print usage.")

	///////////////////////////////////////////////////////////////
	// Vault configuration
	///////////////////////////////////////////////////////////////

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

	///////////////////////////////////////////////////////////////
	// Client Settings configuration
	///////////////////////////////////////////////////////////////

	// ClientSettingsVaultMountPath is the mount path of the client settings in Vault.
	ClientSettingsVaultMountPath = flag.String("client-settings-vault-mount-path", "client-settings", "Mount path of the client settings in Vault. (environment variable: CLIENT_SETTINGS_VAULT_MOUNT_PATH)")

	// ClientSettingsVaultPath is the path of the client settings in Vault.
	ClientSettingsVaultPath = flag.String("client-settings-vault-path", "", "Path of the client settings in Vault. (environment variable: CLIENT_SETTINGS_VAULT_PATH)")

	// ClientSettingsProviderRefreshInterval is the interval at which the client settings provider refreshes the client settings.
	ClientSettingsProviderRefreshInterval = flag.Int("client-settings-provider-refresh-interval", 30, "Interval at which the client settings provider refreshes the client settings. (environment variable: CLIENT_SETTINGS_PROVIDER_REFRESH_INTERVAL)")

	// ClientSettingsApiKeys is the comma-separated list of API keys to use to authenticate with the client settings provider.
	ClientSettingsApiKeys = flag.String("client-settings-api-keys", "", "Comma-separated list of API keys to use to authenticate with the client settings provider. (environment variable: CLIENT_SETTINGS_API_KEYS)")

	// ClientSettingsSecuredSettingsList is a comma-seperated list of applications that can only be accessed via game servers or through client settings secured settings API.
	ClientSettingsSecuredSettingsList = flag.String("client-settings-secured-settings-list", "", "Comma-separated list of applications that can only be accessed via game servers or through client settings secured settings API. (environment variable: CLIENT_SETTINGS_SECURED_SETTINGS_LIST)")

	///////////////////////////////////////////////////////////////
	// Avatar Configuration
	///////////////////////////////////////////////////////////////

	// AvatarApiUrl is the URL of the Avatar API.
	AvatarApiUrl = flag.String("avatar-api-url", "https://avatar.roblox.com", "URL of the Avatar API. (environment variable: AVATAR_API_URL)")

	// AvatarFetchCacheInvalidationInterval is the interval at which the avatar fetch cache is invalidated.
	AvatarFetchCacheInvalidationInterval = flag.Duration("avatar-fetch-cache-invalidation-interval", 5*time.Minute, "Interval at which the avatar fetch cache is invalidated. (environment variable: AVATAR_FETCH_CACHE_INVALIDATION_INTERVAL)")

	// AvatarFetchCacheItemExpiration is the expiration time of an item in the avatar fetch cache.
	AvatarFetchCacheItemExpiration = flag.Duration("avatar-fetch-cache-item-expiration", 5*time.Minute, "Expiration time of an item in the avatar fetch cache. (environment variable: AVATAR_FETCH_CACHE_ITEM_EXPIRATION)")

	// AvatarApiShouldDowngradeBodyColorsFormat is whether or not to downgrade the body colors format.
	AvatarApiShouldDowngradeBodyColorsFormat = flag.Bool("avatar-api-should-downgrade-body-colors-format", true, "Whether or not to downgrade the body colors format. (environment variable: AVATAR_API_SHOULD_DOWNGRADE_BODY_COLORS_FORMAT)")

	///////////////////////////////////////////////////////////////
	// Ephemeral Counters Configuration
	///////////////////////////////////////////////////////////////

	// EphemeralCountersInfluxDbReportingEnabled is whether or not to enable InfluxDB reporting for ephemeral counters.
	EphemeralCountersInfluxDbReportingEnabled = flag.Bool("ephemeral-counters-influxdb-reporting-enabled", false, "Whether or not to enable InfluxDB reporting for ephemeral counters. (environment variable: EPHEMERAL_COUNTERS_INFLUXDB_REPORTING_ENABLED)")

	// EphemeralCountersInfluxDbMaxBatchSize is the maximum batch size to use when reporting ephemeral counters to InfluxDB.
	EphemeralCountersInfluxDbMaxBatchSize = flag.Int("ephemeral-counters-influxdb-max-batch-size", 5000, "Maximum batch size to use when reporting ephemeral counters to InfluxDB. (environment variable: EPHEMERAL_COUNTERS_INFLUXDB_MAX_BATCH_SIZE)")

	// EphemeralCountersInfluxDbFlushInterval is the interval at which to flush the ephemeral counters to InfluxDB.
	EphemeralCountersInfluxDbFlushInterval = flag.Duration("ephemeral-counters-influxdb-flush-interval", 5*time.Second, "Interval at which to flush the ephemeral counters to InfluxDB. (environment variable: EPHEMERAL_COUNTERS_INFLUXDB_FLUSH_INTERVAL)")

	// EphemeralCountersInfluxDbReportingUrl is the URL of the InfluxDB server to report ephemeral counters to.
	EphemeralCountersInfluxDbReportingUrl = flag.String("ephemeral-counters-influxdb-reporting-url", "http://localhost:8086", "URL of the InfluxDB server to report ephemeral counters to. (environment variable: EPHEMERAL_COUNTERS_INFLUXDB_REPORTING_URL)")

	// EphemeralCountersInfluxDbReportingDatabase is the database to report ephemeral counters to.
	EphemeralCountersInfluxDbReportingDatabase = flag.String("ephemeral-counters-influxdb-reporting-database", "ephemeral_counters", "Database to report ephemeral counters to. (environment variable: EPHEMERAL_COUNTERS_INFLUXDB_REPORTING_DATABASE)")

	// EphemeralCountersInfluxDbReportingToken is the token to use to authenticate with InfluxDB.
	EphemeralCountersInfluxDbReportingToken = flag.String("ephemeral-counters-influxdb-reporting-token", "", "Token to use to authenticate with InfluxDB. (environment variable: EPHEMERAL_COUNTERS_INFLUXDB_REPORTING_TOKEN)")

	// EphemeralCountersInfluxDbReportingOrganization is the organization to report ephemeral counters to.
	EphemeralCountersInfluxDbReportingOrganization = flag.String("ephemeral-counters-influxdb-reporting-organization", "mfdlabs", "Organization to report ephemeral counters to. (environment variable: EPHEMERAL_COUNTERS_INFLUXDB_REPORTING_ORGANIZATION)")

	///////////////////////////////////////////////////////////////
	// Version Compatibility Configuration
	///////////////////////////////////////////////////////////////

	// VersionCompatibilityAllowedClientMD5Hashes is the comma-separated list of allowed client MD5 hashes.
	VersionCompatibilityAllowedClientMD5Hashes = flag.String("version-compatibility-allowed-client-md5-hashes", "", "Comma-separated list of allowed client MD5 hashes. (environment variable: VERSION_COMPATIBILITY_ALLOWED_CLIENT_MD5_HASHES)")

	///////////////////////////////////////////////////////////////
	// Asset Delivery API Configuration
	///////////////////////////////////////////////////////////////

	// AssetDeliveryApiUrl is the URL of the Asset Delivery API.
	AssetDeliveryApiUrl = flag.String("asset-delivery-api-url", "https://assetdelivery.roblox.com", "URL of the Asset Delivery API. (environment variable: ASSET_DELIVERY_API_URL)")
)

const FlagsUsageString string = `
	[-h|--help] [--bind-address-ipv4[=:8080]] 
	[--vault-address=http://localhost:8200] [--vault-credential[=]] [--vault-authentication-type[=token]] [--vault-namespace[=]] [--vault-skip-tls-verify[=false]]
	[--client-settings-vault-mount-path[=client-settings]] [--client-settings-vault-path[=]] [--client-settings-provider-refresh-interval[=30]] [--client-settings-api-keys[=]]
	[--avatar-api-url[=https://avatar.roblox.com]] [--avatar-fetch-cache-invalidation-interval[=5m0s]] [--avatar-fetch-cache-item-expiration[=5m0s]] [--avatar-api-should-downgrade-body-colors-format[=true]]
	[--ephemeral-counters-influxdb-reporting-enabled[=false]] [--ephemeral-counters-influxdb-reporting-interval[=10s]] [--ephemeral-counters-influxdb-reporting-url[=http://localhost:8086]] [--ephemeral-counters-influxdb-reporting-database[=ephemeral_counters]] [--ephemeral-counters-influxdb-reporting-token[=]] [--ephemeral-counters-influxdb-reporting-organization[=mfdlabs]]
	[--version-compatibility-allowed-client-md5-hashes[=]]
	[--asset-delivery-api-url[=https://assetdelivery.roblox.com]]`
