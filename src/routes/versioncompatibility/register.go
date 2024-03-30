package versioncompatibility

import (
	"sync"

	"github.com/gorilla/mux"
)

var registerOnce sync.Once

// RegisterRoutes registers the version compatibility routes.
func RegisterRoutes(r *mux.Router) {
	registerOnce.Do(func() {
		r.HandleFunc("/{route:(?:(?i)GetAllowedMD5Hashes)\\/?}", getAllowedMd5Hashes).Methods("GET")
	})
}
