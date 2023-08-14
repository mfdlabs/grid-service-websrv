package middleware

// This package provides recovery middleware

import (
	"net/http"

	"github.com/golang/glog"
)

// RecoveryMiddleware recovers from panics.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if errInterface := recover(); errInterface != nil {
				err := errInterface.(error)
				glog.Errorf("Panic: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
