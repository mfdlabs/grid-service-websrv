package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ConcurrentRequestsGuage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_server_concurrent_requests_total",
		Help: "Number of concurrent requests",
	})

	HttpRequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_server_requests_total",
			Help: "Total number of http requests",
		},
		[]string{"method", "endpoint"},
	)
)
