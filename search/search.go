package search

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/blockhash"
	"net/http"
	"fmt"
	"strings"
)

func FormatResponse(orders []dbModule.Order, format string) ([]byte, string, error) {
	if format == "application/octet-stream" {
		result := []byte{}
		for _, order := range orders {
			orderBytes := order.Bytes()
			result = append(result, orderBytes[:]...)
		}
		return result, "application/octet-stream", nil
	} else {
		orderList := []types.Order{}
		for _, order := range orders {
			orderList = append(orderList, order.Order)
		}
		result, err := json.Marshal(orderList)
		return result, "application/json", err
	}
}

func Handler(db *gorm.DB, blockHash blockhash.BlockHash) func(http.ResponseWriter, *http.Request) {
	blockHash.Get() // Start the go routines, if necessary
	// TODO: Filters
	return func(w http.ResponseWriter, r *http.Request) {
		queryObject := r.URL.Query()
		hash := queryObject.Get("blockhash")
		if hash == "" {
			queryObject.Set("blockhash", blockHash.Get())
			url := *r.URL
			url.RawQuery = queryObject.Encode()
			http.Redirect(w, r, (&url).RequestURI(), 307)
			return
		}


		orders := []dbModule.Order{}
		query := db.Model(&dbModule.Order{})
		// Filter Stuff
		if err := query.Find(orders).Error; err != nil {
			w.WriteHeader(500)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err.Error())))
		}
		var acceptHeader string
		if acceptVal, ok := r.Header["Accept"]; ok {
			acceptHeader = strings.Split(acceptVal[0], ";")[0]
		} else {
			acceptHeader = "unknown"
		}
		response, contentType, err := FormatResponse(orders, acceptHeader)
		if err == nil {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", contentType)
			w.Write(response)
		} else {
			w.WriteHeader(500)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err.Error())))
		}
	}
}
