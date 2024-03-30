package clientsettings

import (
	"net/http"
	"strings"

	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

func getApplicationSettings(w http.ResponseWriter, r *http.Request) {
	applicationName := r.URL.Query().Get("applicationName")
	if applicationName == "" {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The application name is invalid.")
		return
	}

	applicationSettings, depends, allowedFromCsApi, ok := clientSettingsProvider.GetGroup(applicationName)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The application name is invalid.")
		return
	}

	if !allowedFromCsApi && !apiKeyIsValidForThisRequest(r) {
		w.WriteHeader(http.StatusUnauthorized)
		httphelpers.WriteRobloxJSONError(w, "The application name is invalid.")
		return
	}

	if len(depends) > 0 {
		w.Header().Set("X-Depends-On", strings.Join(depends, ", "))
	}

	response := &getApplicationSettingsResponse{
		ApplicationSettings: applicationSettings,
	}

	httphelpers.WriteJSON(w, response)
}
