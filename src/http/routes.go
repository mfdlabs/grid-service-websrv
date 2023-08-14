package http

import (
	"github.com/gorilla/mux"
	"github.com/mfdlabs/grid-service-websrv/routes/avatar"
	"github.com/mfdlabs/grid-service-websrv/routes/clientsettings"
	"github.com/mfdlabs/grid-service-websrv/routes/ephemeralcounters"
	"github.com/mfdlabs/grid-service-websrv/routes/versioncompatibility"
)

// RegisterRoutes registers the HTTP routes.
func registerRoutes(r *mux.Router) {
	clientsettings.RegisterRoutes(r)
	avatar.RegisterRoutes(r)
	versioncompatibility.RegisterRoutes(r)
	ephemeralcounters.RegisterRoutes(r)
}
