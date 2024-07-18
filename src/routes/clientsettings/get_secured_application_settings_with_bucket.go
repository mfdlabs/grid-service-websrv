package clientsettings

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

func getSecuredApplicationSettingsWithBucket(w http.ResponseWriter, r *http.Request) {
	if !apiKeyIsValidForThisRequest(r) {
		w.WriteHeader(http.StatusUnauthorized)
		httphelpers.WriteRobloxJSONError(w, "The api key is invalid..")
		return
	}

	vars := mux.Vars(r)

	applicationName := vars["applicationName"]
	if applicationName == "" {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The application name is invalid.")
		return
	}

	if !isApplicationSecured(applicationName) {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The application name is invalid.")
		return
	}

	bucketName := vars["bucketName"]
	if bucketName == "" {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The bucket name is invalid.")
		return
	}

	applicationSettings, depends, _, ok := clientSettingsProvider.GetGroupWithBucket(applicationName, bucketName)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
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
