package metrics

import (
	"net/http"

	metrics_counters "github.com/mfdlabs/grid-service-websrv/metrics"
)

func HttpServerRequestCountMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics_counters.HttpRequestCounter.WithLabelValues(r.Method, r.URL.Path).Inc()

		h.ServeHTTP(w, r)
	})
}
