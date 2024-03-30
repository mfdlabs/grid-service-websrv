package assetdelivery

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mfdlabs/grid-service-websrv/flags"
)

// It at least needs to have one of the following query params.
var requiredQueryParams = []string{
	"id",
	"userAssetId",
	"assetVersionId",
	"hash",
	"marAssetHash", // Moderation Agnostic Request Asset Hash
}

func readAndCleanQueryParams(r *http.Request) string {
	queryParams := make([]string, len(r.URL.Query()))

	query := r.URL.Query()
	for key, value := range query {
		// Ensure that the query param is not empty and url encoded.
		if len(value) > 0 {
			queryParams = append(queryParams, url.QueryEscape(key)+"="+url.QueryEscape(value[0]))
		}
	}

	return strings.Join(queryParams, "&")
}

func buildAssetDeliveryUrl() string {
	return *flags.AssetDeliveryApiUrl + "/v1/asset/?"
}

func applyHeadersToClientRequest(r *http.Request, clientRequest *http.Request) {
	// Copy headers from the original request to the client request.
	for key, values := range r.Header {
		for _, value := range values {
			// Do not pass the "Host" header to the client request.
			if strings.ToLower(key) == "host" {
				continue
			}

			clientRequest.Header.Add(key, value)
		}
	}
}

func buildRawAssetHttpRequest(r *http.Request, assetUrl string) (*http.Request, error) {
	// Build the client request.
	clientRequest, err := http.NewRequest(r.Method, assetUrl, nil)
	if err != nil {
		return nil, err
	}

	// Copy headers from the original request to the client request.
	applyHeadersToClientRequest(r, clientRequest)

	return clientRequest, nil
}

func getAsset(w http.ResponseWriter, r *http.Request) {
	// No query, no request.
	if len(r.URL.Query()) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the request has at least one of the required query params.
	hasRequiredQueryParam := false
	for _, requiredQueryParam := range requiredQueryParams {
		if _, ok := r.URL.Query()[requiredQueryParam]; ok {
			hasRequiredQueryParam = true
			break
		}
	}

	if !hasRequiredQueryParam {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Read and clean the query params.
	queryParams := readAndCleanQueryParams(r)

	// Build the asset delivery URL.
	assetUrl := buildAssetDeliveryUrl() + queryParams

	// Build the client request.
	clientRequest, err := buildRawAssetHttpRequest(r, assetUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the client request.
	clientResponse, err := rawAssetsClient.Do(clientRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check for redirect, it will be redirecting to the CDN.
	if clientResponse.StatusCode == http.StatusFound {
		// Follow the redirect on our end.
		// This is to ensure that we can add the correct headers to the response.
		// and to ensure that the response is in the correct format.

		w.Header().Set("Location", clientResponse.Header.Get("Location"))
		w.WriteHeader(http.StatusFound)
		return
	}

	// Copy headers from the client response to the response.
	for key, values := range clientResponse.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Read body and write it to the response.
	reader := io.NopCloser(clientResponse.Body)
	defer reader.Close()

	w.WriteHeader(clientResponse.StatusCode)
	_, err = io.Copy(w, reader)
	if err != nil {
		return
	}
}
