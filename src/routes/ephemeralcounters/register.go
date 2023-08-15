package ephemeralcounters

import (
	"sync"

	"github.com/gorilla/mux"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/mfdlabs/grid-service-websrv/flags"
)

var (
	registerOnce sync.Once

	client   influxdb2.Client
	writeAPI api.WriteAPI
)

// RegisterRoutes registers the routes for this package
func RegisterRoutes(r *mux.Router) {
	registerOnce.Do(func() {
		if *flags.EphemeralCountersInfluxDbReportingEnabled {
			client = influxdb2.NewClientWithOptions(
				*flags.EphemeralCountersInfluxDbReportingUrl,
				*flags.EphemeralCountersInfluxDbReportingToken,
				influxdb2.DefaultOptions().SetBatchSize(uint(*flags.EphemeralCountersInfluxDbMaxBatchSize)).SetFlushInterval(uint(*flags.EphemeralCountersInfluxDbFlushInterval)),
			)

			writeAPI = client.WriteAPI(*flags.EphemeralCountersInfluxDbReportingOrganization, *flags.EphemeralCountersInfluxDbReportingDatabase)
		}

		r.HandleFunc("/v1.0/SequenceStatistics/{route:BatchAddToSequencesV2\\/?}", batchAddToSequencesV2).Methods("POST")
		r.HandleFunc("/v1.1/Counters/{route:BatchIncrement\\/?}", batchIncrement).Methods("POST")
	})
}

// CloseWriteAPI closes the write API
func CloseWriteAPI() {
	if writeAPI != nil {
		client.Close()
	}
}
