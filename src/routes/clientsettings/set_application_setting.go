package clientsettings

import (
	"net/http"

	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

func setApplicationSetting(w http.ResponseWriter, r *http.Request) {
	if !apiKeyIsValidForThisRequest(r) {
		w.WriteHeader(http.StatusUnauthorized)
		httphelpers.WriteRobloxJSONError(w, "The API key is invalid.")
		return
	}

	var request setApplicationSettingRequest
	if !httphelpers.ReadJSON(w, r, &request) {
		return
	}

	if request.ApplicationName == "" {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The application name is invalid.")
		return
	}

	if request.SettingName == "" {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The setting name is invalid.")
		return
	}

	oldValue, didExist := clientSettingsProvider.Get(request.ApplicationName, request.SettingName)

	didUpdate, err := clientSettingsProvider.Set(request.ApplicationName, request.SettingName, request.Value)
	if err != nil {
		httphelpers.WriteRobloxJSONErr(w, err)
		return
	}

	response := &setApplicationSettingResponse{
		DidUpdate: didUpdate,
	}

	if didExist {
		response.OldValue = &getApplicationSettingResponse{
			Name:  request.SettingName,
			Value: oldValue,
		}
	}

	response.NewValue = &getApplicationSettingResponse{
		Name:  request.SettingName,
		Value: request.Value,
	}

	w.WriteHeader(http.StatusOK)
	httphelpers.WriteJSON(w, response)
}
