package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	GroupReads = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_settings_application_reads_total",
		Help: "The total number of application reads",
	},
		[]string{"application_name"},
	)

	GroupWrites = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_settings_application_writes_total",
		Help: "The total number of application writes",
	},
		[]string{"application_name"},
	)

	SettingsRefreshes = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_settings_refreshes_total",
		Help: "The total number of settings refreshes",
	},
		[]string{"application_name"},
	)
)
