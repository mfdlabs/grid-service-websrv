package clientsettings

import (
	"net/http"

	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

func importApplicationSettings(w http.ResponseWriter, r *http.Request) {
	if !apiKeyIsValidForThisRequest(r) {
		w.WriteHeader(http.StatusUnauthorized)
		httphelpers.WriteRobloxJSONError(w, "The API key is invalid.")
		return
	}

	var request importApplicationSettingsRequest
	if !httphelpers.ReadJSON(w, r, &request) {
		return
	}

	if request.ApplicationName == "" {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The application name is invalid.")
		return
	}

	err := clientSettingsProvider.Import(
		request.ApplicationName,
		request.ApplicationSettings,
		request.Dependendies,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		httphelpers.WriteRobloxJSONErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
