package clientsettings

import (
	"net/http"
	"strings"

	"github.com/mfdlabs/grid-service-websrv/flags"
)

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
