package ephemeralcounters

import (
	"net/http"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

func parseBatchIncrementRequest(w http.ResponseWriter, r *http.Request) []*write.Point {
	points := []*write.Point{}

	// Request is a map of strings to values
	var request map[string]int64
	if !httphelpers.ReadJSON(w, r, &request) {
		return points
	}

	for key, value := range request {
		points = append(points, write.NewPoint(
			"ephemeral_counters",
			map[string]string{
				"key": key,
			},
			map[string]interface{}{
				"value": value,
			},
			time.Now(),
		))
	}

	return points
}

func batchIncrement(w http.ResponseWriter, r *http.Request) {
	if writeAPI != nil {
		points := parseBatchIncrementRequest(w, r)

		for _, point := range points {
			writeAPI.WritePoint(point)
		}
	}

	w.WriteHeader(http.StatusOK)
}
