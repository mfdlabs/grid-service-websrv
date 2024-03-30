package assetdelivery

import (
	"net/http"
	"strconv"

	assetdeliveryv1 "github.com/mfdlabs/grid-service-websrv/assetdelivery_v1"
	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

func batchGetAssets(w http.ResponseWriter, r *http.Request) {
	var request []assetdeliveryv1.RobloxWebAssetsBatchAssetRequestItem
	if !httphelpers.ReadJSON(w, r, &request) {
		return
	}

	var (
		robloxPlaceId                                 int64
		acceptHeader, robloxBrowserAssetRequestHeader string
	)

	robloxPlaceId, err := strconv.ParseInt(r.Header.Get("Roblox-Place-Id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The place ID is invalid.")
		return
	}

	acceptHeader = r.Header.Get("Accept")
	robloxBrowserAssetRequestHeader = r.Header.Get("Roblox-Browser-Asset-Request")

	response, _, err := batchAssetsClient.BatchAPI.V1AssetsBatchPost(r.Context()).RobloxPlaceId(robloxPlaceId).Accept(acceptHeader).RobloxBrowserAssetRequest(robloxBrowserAssetRequestHeader).AssetRequestItems(request).Execute()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		httphelpers.WriteRobloxJSONError(w, err.Error())
		return
	}

	httphelpers.WriteJSON(w, response)
}
