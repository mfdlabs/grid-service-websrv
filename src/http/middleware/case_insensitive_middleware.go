package middleware

import (
	"net/http"
	"strings"
)

// CaseInsensitiveMiddleware is a middleware that makes all the request paths lowercase.
func CaseInsensitiveMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
