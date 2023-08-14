package ephemeralcounters

import (
	"net/http"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
	httphelpers "github.com/mfdlabs/grid-service-websrv/http_helpers"
)

func parseBatchAddToSequencesV2Request(w http.ResponseWriter, r *http.Request) []*write.Point {
	points := []*write.Point{}

	// Request is a map of strings to values
	var request []*ephemeralStatistic
	if !httphelpers.ReadJSON(w, r, &request) {
		return points
	}

	for _, stat := range request {
		points = append(points, write.NewPoint(
			"ephemeral_statistics",
			map[string]string{
				"key": stat.Key,
			},
			map[string]interface{}{
				"value": stat.Value,
			},
			time.Now(),
		))
	}

	return points
}

func batchAddToSequencesV2(w http.ResponseWriter, r *http.Request) {
	if writeAPI != nil {
		points := parseBatchAddToSequencesV2Request(w, r)

		for _, point := range points {
			writeAPI.WritePoint(point)
		}
	}

	w.WriteHeader(http.StatusOK)
}
