package clientsettings

import (
	"net/http"

	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

func refreshAllApplicationSettings(w http.ResponseWriter, r *http.Request) {
	if !apiKeyIsValidForThisRequest(r) {
		w.WriteHeader(http.StatusUnauthorized)
		httphelpers.WriteRobloxJSONError(w, "Invalid API key.")
		return
	}

	if err := clientSettingsProvider.Refresh(); err != nil {
		httphelpers.WriteRobloxJSONErr(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
