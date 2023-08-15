package middleware

import (
	"net/http"

	"github.com/golang/glog"
)

// LoggingMiddleware logs the request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/metrics" {
			glog.Infof("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		}

		next.ServeHTTP(w, r)

	})
}
