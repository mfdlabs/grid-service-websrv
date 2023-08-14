package middleware

import (
	"net/http"
	"time"

	"github.com/golang/glog"
)

// LoggingMiddleware logs the request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		start := time.Now()

		next.ServeHTTP(w, r)

		glog.Infof("%s %s %s %s", r.RemoteAddr, r.Method, r.URL, time.Since(start))
	})
}
