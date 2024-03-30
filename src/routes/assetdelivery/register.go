package assetdelivery

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	assetdeliveryv1 "github.com/mfdlabs/grid-service-websrv/assetdelivery_v1"
	"github.com/mfdlabs/grid-service-websrv/flags"
)

var (
	registerOnce sync.Once

	batchAssetsClient *assetdeliveryv1.APIClient
	rawAssetsClient   *http.Client
)

// RegisterRoutes registers the asset delivery routes.
func RegisterRoutes(r *mux.Router) {
	registerOnce.Do(func() {
		batchAssetsClient = assetdeliveryv1.NewAPIClient(assetdeliveryv1.NewConfiguration(*flags.AssetDeliveryApiUrl))
		rawAssetsClient = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		r.HandleFunc("/{route:(?:(?i)v1\\/asset)\\/?}", getAsset).Methods("GET")
		r.HandleFunc("/{route:(?:(?i)v1\\/assets/batch)\\/?}", batchGetAssets).Methods("POST")
	})
}
