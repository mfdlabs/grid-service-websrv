package routes

import (
	"net/http"

	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

// StubEmpty is a stub HTTP handler.
func StubEmpty(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	httphelpers.WriteJSON(w, []string{})
}
