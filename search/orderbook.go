package search

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/common"
	dbModule "github.com/notegio/openrelay/db"
	"net/http"
)

type OrderBook struct {
	Asks []dbModule.Order `json:"asks"`
	Bids []dbModule.Order `json:"bids"`
}

func OrderBookHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queryObject := r.URL.Query()
		baseTokenAddressHex := queryObject.Get("baseTokenAddress")
		quoteTokenAddressHex := queryObject.Get("quoteTokenAddress")
		if baseTokenAddressHex == "" || quoteTokenAddressHex == "" {
			returnError(w, errors.New("Must provide baseTokenAddress and quoteTokenAddress "), 404)
			return
		}
		baseTokenAddressBytes, err := common.HexToBytes(baseTokenAddressHex)
		if err != nil {
			returnError(w, err, 400)
			return
		}
		quoteTokenAddressBytes, err := common.HexToBytes(quoteTokenAddressHex)
		if err != nil {
			returnError(w, err, 400)
			return
		}
		currentTime := getExpTime(queryObject)
		baseTokenAddress := common.BytesToOrAddress(baseTokenAddressBytes)
		quoteTokenAddress := common.BytesToOrAddress(quoteTokenAddressBytes)
		orderBook := &OrderBook{[]dbModule.Order{}, []dbModule.Order{}}
		db.Model(&dbModule.Order{}).Where("status = ?", dbModule.StatusOpen).Where("expiration_timestamp_in_sec > ?", currentTime).Where("taker_token = ? AND maker_token = ?", baseTokenAddress, quoteTokenAddress).Order("price, fee_rate, expiration_timestamp_in_sec").Find(&orderBook.Bids)
		db.Model(&dbModule.Order{}).Where("status = ?", dbModule.StatusOpen).Where("expiration_timestamp_in_sec > ?", currentTime).Where("maker_token = ? AND taker_token = ?", baseTokenAddress, quoteTokenAddress).Order("price, fee_rate, expiration_timestamp_in_sec").Find(&orderBook.Asks)
		response, err := json.Marshal(orderBook)
		if err != nil {
			returnError(w, err, 500)
		} else {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		}
	}
}
