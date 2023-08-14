package http

import (
	"github.com/gorilla/mux"
	"github.com/mfdlabs/grid-service-websrv/routes/clientsettings"
)

// RegisterRoutes registers the HTTP routes.
func registerRoutes(r *mux.Router) {
	clientsettings.RegisterRoutes(r)
}
