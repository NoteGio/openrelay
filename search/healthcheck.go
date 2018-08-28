package search

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/blockhash"
	"net/http"
	"strings"
)

type HealthCheck struct {
	BlockHash string
}

func HealthCheckHandler(db *gorm.DB, blockHash blockhash.BlockHash) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		hc := &HealthCheck{}
		if hash := strings.Trim(blockHash.Get(), "\""); hash != "" {
			hc.BlockHash = hash
		} else {
			returnError(w, errors.New("Got empty blockhash"), 500)
			return
		}
		if db != nil {
			if err := db.Raw("SELECT 1").Error; err != nil {
				returnError(w, err, 500)
				return
			}
		}
		response, err := json.Marshal(hc)
		if err != nil {
			returnError(w, err, 500)
			return
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
