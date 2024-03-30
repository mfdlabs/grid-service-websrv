package versioncompatibility

import (
	"net/http"
	"strings"

	"github.com/mfdlabs/grid-service-websrv/flags"
	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

type allowedMD5HashesResponse struct {
	Data []string `json:"data"`
}

func getAllowedMd5Hashes(w http.ResponseWriter, _ *http.Request) {
	response := allowedMD5HashesResponse{
		Data: make([]string, 0),
	}

	// Parse allowed MD5 hashes.
	allowedMD5Hashes := *flags.VersionCompatibilityAllowedClientMD5Hashes
	if allowedMD5Hashes == "" {
		httphelpers.WriteJSON(w, response)
		return
	}

	// Split allowed MD5 hashes.
	allowedMD5HashesSplit := strings.Split(allowedMD5Hashes, ",")
	for _, allowedMD5Hash := range allowedMD5HashesSplit {
		if allowedMD5Hash != "" {
			response.Data = append(response.Data, allowedMD5Hash)
		}
	}

	httphelpers.WriteJSON(w, response)
}
