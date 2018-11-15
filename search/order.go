package search

import (
	"encoding/hex"
	"errors"
	"github.com/jinzhu/gorm"
	dbModule "github.com/notegio/openrelay/db"
	"net/http"
	"regexp"
	"strings"
)

func OrderHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	orderRegex := regexp.MustCompile(".*/order/0x([0-9a-fA-F]+)")
	return func(w http.ResponseWriter, r *http.Request) {
		pathMatch := orderRegex.FindStringSubmatch(r.URL.Path)
		if len(pathMatch) == 0 {
			returnError(w, errors.New("Malformed order hash"), 404)
			return
		}
		hashHex := pathMatch[1]
		hashBytes, err := hex.DecodeString(hashHex)
		if err != nil {
			returnError(w, err, 400)
			return
		}
		order := &dbModule.Order{}
		query := db.Model(&dbModule.Order{}).Where("order_hash = ?", hashBytes).First(order)
		if query.Error != nil {
			if query.Error.Error() == "record not found" {
				returnError(w, query.Error, 404)
			} else {
				returnError(w, query.Error, 500)
			}
			return
		}
		var acceptHeader string
		if acceptVal, ok := r.Header["Accept"]; ok {
			acceptHeader = strings.Split(acceptVal[0], ";")[0]
		} else {
			acceptHeader = "unknown"
		}
		response, contentType, err := FormatSingleResponse(order, acceptHeader)
		if err == nil {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", contentType)
			w.Write(response)
		} else {
			returnError(w, err, 500)
		}
	}
}
