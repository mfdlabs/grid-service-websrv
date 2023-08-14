package clientsettings

import (
	"net/http"
	"strings"

	"github.com/mfdlabs/grid-service-websrv/flags"
	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

type importApplicationSettingsRequest struct {
	ApplicationName     string                 `json:"applicationName"`
	ApplicationSettings map[string]interface{} `json:"applicationSettings"`
	Dependendies        []string               `json:"dependencies"`
}

func apiKeyIsValidForThisRequest(r *http.Request) bool {
	apiKey := r.Header.Get("X-Api-Key")
	apiKeys := strings.Split(*flags.ClientSettingsApiKeys, ",")

	if len(apiKeys) == 1 && apiKeys[0] == "" {
		return true
	}

	if apiKey == "" {
		return false
	}

	for _, validApiKey := range apiKeys {
		if apiKey == validApiKey {
			return true
		}
	}

	return false
}

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
