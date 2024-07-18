package clientsettings

import (
	"context"
	"sync"

	"github.com/gorilla/mux"

	clientsettingsutil "github.com/mfdlabs/grid-service-websrv/clientsettings_util"
	"github.com/mfdlabs/grid-service-websrv/flags"
	"github.com/mfdlabs/grid-service-websrv/vault"
)

var (
	registerOnce sync.Once

	clientSettingsProvider *clientsettingsutil.ClientSettingsProvider
)

// RegisterRoutes registers the client settings routes.
func RegisterRoutes(r *mux.Router) {
	registerOnce.Do(func() {
		client, err := vault.GetGlobalVaultClient(context.Background())
		if err != nil {
			panic(err)
		}

		clientSettingsProvider = clientsettingsutil.NewClientSettingsProvider(
			client.Client,
			*flags.ClientSettingsVaultMountPath,
			*flags.ClientSettingsVaultPath,
			*flags.ClientSettingsProviderRefreshInterval,
		)

		r.HandleFunc("/{route:(?:(?i)v1\\/settings\\/application)\\/?}", getApplicationSettings).Methods("GET")
		r.HandleFunc("/{route:(?:(?i)v1\\/settings\\/application)\\/?}", importApplicationSettings).Methods("POST")
		r.HandleFunc("/{route:(?:(?i)v1/settings)\\/?}", refreshAllApplicationSettings).Methods("POST")

		r.HandleFunc("/{route:(?:(?i)v1\\/settings\\/application\\/setting)\\/?}", getApplicationSetting).Methods("GET")
		r.HandleFunc("/{route:(?:(?i)v1\\/settings\\/application\\/setting)\\/?}", setApplicationSetting).Methods("POST")

		r.HandleFunc("/v2/settings/application/{applicationName}", getApplicationSettingsV2).Methods("GET")
		r.HandleFunc("/v2/settings/application/{applicationName}/bucket/{bucketName}", getApplicationSettingsWithBucket).Methods("GET")

		r.HandleFunc("/v2/settings/secured-settings/{applicationName}", getSecuredApplicationSettings).Methods("GET")
		r.HandleFunc("/v2/settings/secured-settings/{applicationName}/bucket/{bucketName}", getSecuredApplicationSettingsWithBucket).Methods("GET")
	})
}
