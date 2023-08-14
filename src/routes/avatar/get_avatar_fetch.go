package avatar

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/mfdlabs/grid-service-websrv/flags"
	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

func getAvatarCacheKey(userId, placeId int64) string {
	return fmt.Sprintf("avatar_fetch:%d:%d", userId, placeId)
}

func getAvatarFetch(w http.ResponseWriter, r *http.Request) {
	userId, err := httphelpers.ParseInt64FromQuery(r, "userId")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The user ID is invalid.")
		return
	}

	placeId, err := httphelpers.ParseInt64FromQuery(r, "placeId")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		httphelpers.WriteRobloxJSONError(w, "The place ID is invalid.")
		return
	}

	cacheKey := getAvatarCacheKey(userId, placeId)
	cachedResponse, ok := localCache.Get(cacheKey)
	if ok {
		httphelpers.WriteJSON(w, cachedResponse)
		return
	}

	avatarFetchResponse, _, err := avatarApiClient.AvatarApi.V1AvatarFetchGet(context.Background()).UserId(userId).PlaceId(placeId).Execute()
	if err != nil {
		glog.Errorf("Failed to fetch avatar: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		httphelpers.WriteRobloxJSONError(w, "An unexpected error occurred.")
		return
	}

	response := fromNewAvatarFetchResponse(avatarFetchResponse)
	localCache.Set(cacheKey, response, *flags.AvatarFetchCacheItemExpiration)

	httphelpers.WriteJSON(w, response)
}
