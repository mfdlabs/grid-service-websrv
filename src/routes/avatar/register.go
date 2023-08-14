package avatar

import (
	"sync"

	"github.com/gorilla/mux"

	avatarv1 "github.com/mfdlabs/grid-service-websrv/avatar_v1"
	"github.com/mfdlabs/grid-service-websrv/cache"
	"github.com/mfdlabs/grid-service-websrv/flags"
)

var (
	registerOnce sync.Once

	avatarApiClient *avatarv1.APIClient
	localCache      *cache.LocalCache
)

// RegisterRoutes registers the avatar routes.
func RegisterRoutes(r *mux.Router) {
	registerOnce.Do(func() {
		avatarApiClient = avatarv1.NewAPIClient(avatarv1.NewConfiguration(*flags.AvatarApiUrl))
		localCache = cache.NewLocalCache(*flags.AvatarFetchCacheInvalidationInterval)

		r.HandleFunc("/v1/avatar-fetch", getAvatarFetch).Methods("GET")
	})
}
