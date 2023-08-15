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
		r.HandleFunc("/{route:[G|g]et[A|a]llowed(?:MD|Md|mD|md)5[H|h]ashes\\/?}", routes.StubEmpty).Methods("GET")
		r.HandleFunc("/{route:[G|g]et[A|a]llowed[V|v]ersions\\/?}", routes.StubEmpty).Methods("GET")
	})
}
