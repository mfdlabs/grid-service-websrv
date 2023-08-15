package http

// Sets up the HTTP server.

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mfdlabs/grid-service-websrv/flags"
	"github.com/mfdlabs/grid-service-websrv/http/middleware"
	"github.com/mfdlabs/grid-service-websrv/http/middleware/metrics"
)

// Start starts the HTTP server.
func Start() {
	r := mux.NewRouter()

	r.Use(middleware.CaseInsensitiveMiddleware)

	// Add middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoveryMiddleware)

	r.Handle("/metrics", promhttp.Handler())

	r.Use(metrics.HttpServerConcurrentRequestsMiddleware)
	r.Use(metrics.HttpServerRequestCountMiddleware)

	registerRoutes(r)

	// Start the server
	srv := &http.Server{
		Addr:         *flags.BindAddressIpv4,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Starting HTTP server on %s\n", *flags.BindAddressIpv4)

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
