package clientsettings

import (
	"net/http"

	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

func getApplicationSettings(w http.ResponseWriter, r *http.Request) {
	applicationName := r.URL.Query().Get("applicationName")
	if applicationName == "" {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The application name is invalid.")
		return
	}

	applicationSettings, ok := clientSettingsProvider.GetGroup(applicationName)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The application name is invalid.")
		return
	}

	response := &getApplicationSettingsResponse{
		ApplicationSettings: applicationSettings,
	}

	httphelpers.WriteJSON(w, response)
}
