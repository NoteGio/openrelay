package search

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/types"
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

func Handler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
