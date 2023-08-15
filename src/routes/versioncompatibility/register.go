package versioncompatibility

import (
	"sync"

	"github.com/gorilla/mux"
	"github.com/mfdlabs/grid-service-websrv/routes"
)

var registerOnce sync.Once

// RegisterRoutes registers the version compatibility routes.
func RegisterRoutes(r *mux.Router) {
	registerOnce.Do(func() {
		r.HandleFunc("/{route:GetAllowedMD5Hashes\\/?}", routes.StubEmpty).Methods("GET")
		r.HandleFunc("/{route:GetAllowedVersions\\/?}", routes.StubEmpty).Methods("GET")
	})
}
