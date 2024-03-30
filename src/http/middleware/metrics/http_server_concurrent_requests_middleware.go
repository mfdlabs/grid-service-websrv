package metrics

import (
	"net/http"

	metrics_counters "github.com/mfdlabs/grid-service-websrv/metrics"
)

func HttpServerConcurrentRequestsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics_counters.ConcurrentRequestsGuage.Inc()

		h.ServeHTTP(w, r)

		metrics_counters.ConcurrentRequestsGuage.Dec()
	})
}
