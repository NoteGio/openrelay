package search

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/types"
	dbModule "github.com/notegio/openrelay/db"
	"net/http"
	// "log"
)

type OrderBook struct {
	Asks *PagedResult `json:"asks"`
	Bids *PagedResult `json:"bids"`
}

func OrderBookHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request, types.Pool) {
	return func(w http.ResponseWriter, r *http.Request, pool types.Pool) {
		queryObject := r.URL.Query()
		errs := []ValidationError{}
		pageInt, perPageInt, err := getPages(queryObject)
		if err != nil {
			errs = append(errs, ValidationError{err.Error(), 1001, "page"})
		}
		baseAssetDataHex := queryObject.Get("baseAssetData")
		quoteAssetDataHex := queryObject.Get("quoteAssetData")
		if baseAssetDataHex == "" || quoteAssetDataHex == "" {
			returnError(w, errors.New("Must provide baseAssetData and quoteAssetData "), 404)
			return
		}
		baseAssetData, err := common.HexToAssetData(baseAssetDataHex)
		if err != nil {
			errs = append(errs, ValidationError{err.Error(), 1001, "baseAssetData"})
		}
		quoteAssetData, err := common.HexToAssetData(quoteAssetDataHex)
		if err != nil {
			errs = append(errs, ValidationError{err.Error(), 1001, "quoteAssetData"})
		}
		currentTime := getExpTime(queryObject)
		baseQuery, err := pool.Filter(db.Model(&dbModule.Order{}).Where("status = ?", dbModule.StatusOpen).Where("expiration_timestamp_in_sec > ?", currentTime))
		if err != nil {
			errs = append(errs, ValidationError{err.Error(), 1001, "pool"})
		}
		if len(errs) > 0 {
			returnErrorList(w, errs)
			return
		}
		bids := []dbModule.Order{}
		asks := []dbModule.Order{}
		var bidCount int
		var askCount int

		// orderBook := &OrderBook{[]dbModule.Order{}, []dbModule.Order{}}
		baseQuery.Where("taker_asset_data = ? AND maker_asset_data = ?", []byte(baseAssetData[:]), []byte(quoteAssetData[:])).Order("price, fee_rate, expiration_timestamp_in_sec").Count(&bidCount)
		baseQuery.Where("maker_asset_data = ? AND taker_asset_data = ?", []byte(baseAssetData[:]), []byte(quoteAssetData[:])).Order("price, fee_rate, expiration_timestamp_in_sec").Count(&askCount)
		if bidCount > (pageInt - 1) * perPageInt {
			// We don't need to bother with this query if te total is less than the
			// offset
			baseQuery.Where("taker_asset_data = ? AND maker_asset_data = ?", []byte(baseAssetData[:]), []byte(quoteAssetData[:])).Order("price, fee_rate, expiration_timestamp_in_sec").Offset((pageInt - 1) * perPageInt).Limit(perPageInt).Find(&bids)
		}
		if askCount > (pageInt - 1) * perPageInt {
			baseQuery.Where("maker_asset_data = ? AND taker_asset_data = ?", []byte(baseAssetData[:]), []byte(quoteAssetData[:])).Order("price, fee_rate, expiration_timestamp_in_sec").Offset((pageInt - 1) * perPageInt).Limit(perPageInt).Find(&asks)
		}
		formattedAsks := []FormattedOrder{}
		formattedBids := []FormattedOrder{}
		for _, order := range asks {
			formattedAsks = append(formattedAsks, *GetFormattedOrder(&order))
		}
		for _, order := range bids {
			formattedBids = append(formattedBids, *GetFormattedOrder(&order))
		}
		orderBook := &OrderBook{
			GetPagedResult(askCount, pageInt, perPageInt, formattedAsks),
			GetPagedResult(bidCount, pageInt, perPageInt, formattedBids),
		}
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
