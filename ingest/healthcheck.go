package ingest

import (
	"encoding/json"
	"gopkg.in/redis.v3"
	"net/http"
)

type HealthCheck struct {
	Time []string
}

func HealthCheckHandler(redisClient *redis.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := redisClient.Time().Result()
		if err != nil {
			returnError(w, IngestError{102, "Internal error", []ValidationError{}}, 500)
		}
		hc := &HealthCheck{t}
		response, err := json.Marshal(hc)
		if err != nil {
			returnError(w, IngestError{102, err.Error(), []ValidationError{}}, 500)
			return
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
