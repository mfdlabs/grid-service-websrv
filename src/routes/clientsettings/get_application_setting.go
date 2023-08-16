package clientsettings

import (
	"net/http"

	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

func getApplicationSetting(w http.ResponseWriter, r *http.Request) {
	applicationName := r.URL.Query().Get("applicationName")
	if applicationName == "" {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The application name is invalid.")
		return
	}

	settingName := r.URL.Query().Get("settingName")
	if settingName == "" {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The setting name is invalid.")
		return
	}

	value, ok := clientSettingsProvider.Get(applicationName, settingName)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		httphelpers.WriteRobloxJSONError(w, "The setting was not found.")
		return
	}

	response := &getApplicationSettingResponse{
		settingName,
		value,
	}

	httphelpers.WriteJSON(w, response)
}
