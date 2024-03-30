package apiproxy

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var (
	registerOnce sync.Once
)

// RegisterRoutes registers the asset delivery routes.
func RegisterRoutes(r *mux.Router) {
	registerOnce.Do(func() {
		r.HandleFunc("/{route:(?:(?i)universes\\/validate-place-join)\\/?}", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("true"))
		}).Methods("POST")
	})
}
