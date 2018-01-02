package search

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/blockhash"
	"github.com/notegio/openrelay/common"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/types"
	"net/http"
	urlModule "net/url"
	"strconv"
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

func FormatSingleResponse(order *dbModule.Order, format string) ([]byte, string, error) {
	if format == "application/octet-stream" {
		result := order.Bytes()
		return result[:], "application/octet-stream", nil
	}
	result, err := json.Marshal(order)
	return result, "application/json", err
}

func applyFilter(query *gorm.DB, queryField, dbField string, queryObject urlModule.Values) (*gorm.DB, error) {
	if address := queryObject.Get(queryField); address != "" {
		addressBytes, err := common.HexToBytes(address)
		if err != nil {
			return query, err
		}
		whereClause := fmt.Sprintf("%v = ?", dbField)
		filteredQuery := query.Where(whereClause, common.BytesToOrAddress(addressBytes))
		return filteredQuery, filteredQuery.Error
	}
	return query, nil
}

func applyOrFilter(query *gorm.DB, queryField, dbField1, dbField2 string, queryObject urlModule.Values) (*gorm.DB, error) {
	if address := queryObject.Get(queryField); address != "" {
		addressBytes, err := common.HexToBytes(address)
		if err != nil {
			return query, err
		}
		whereClause := fmt.Sprintf("%v = ? or %v = ?", dbField1, dbField2)
		filteredQuery := query.Where(whereClause, common.BytesToOrAddress(addressBytes), common.BytesToOrAddress(addressBytes))
		return filteredQuery, filteredQuery.Error
	}
	return query, nil
}

func returnError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err.Error())))
}

func getPages(queryObject urlModule.Values) (int, int, error) {
	pageStr := queryObject.Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	perPageStr := queryObject.Get("per_page")
	if perPageStr == "" {
		perPageStr = "20"
	}
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, err
	}
	perPageInt, err := strconv.Atoi(perPageStr)
	if err != nil {
		return 0, 0, err
	}
	return pageInt, perPageInt, nil
}

func BlockHashDecorator(blockHash blockhash.BlockHash, fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	blockHash.Get() // Start the go routines, if necessary
	return func(w http.ResponseWriter, r *http.Request) {
		queryObject := r.URL.Query()
		hash := queryObject.Get("blockhash")
		if hash == "" {
			queryObject.Set("blockhash", strings.Trim(blockHash.Get(), "\""))
			url := *r.URL
			url.RawQuery = queryObject.Encode()
			w.Header().Set("Cache-Control", "max-age=5, public")
			http.Redirect(w, r, (&url).RequestURI(), 307)
			return
		}
		fn(w, r)
	}
}

func SearchHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queryObject := r.URL.Query()
		query := db.Model(&dbModule.Order{}).Where("status = ?", dbModule.StatusOpen)

		query, err := applyFilter(query, "exchangeContractAddress", "exchange_address", queryObject)
		if err != nil {
			returnError(w, err, 400)
			return
		}
		query, err = applyFilter(query, "makerTokenAddress", "maker_token", queryObject)
		if err != nil {
			returnError(w, err, 400)
			return
		}
		query, err = applyFilter(query, "takerTokenAddress", "taker_token", queryObject)
		if err != nil {
			returnError(w, err, 400)
			return
		}
		query, err = applyFilter(query, "maker", "maker", queryObject)
		if err != nil {
			returnError(w, err, 400)
			return
		}
		query, err = applyFilter(query, "taker", "taker", queryObject)
		if err != nil {
			returnError(w, err, 400)
			return
		}
		query, err = applyFilter(query, "feeRecipient", "fee_recipient", queryObject)
		if err != nil {
			returnError(w, err, 400)
			return
		}
		query, err = applyOrFilter(query, "tokenAddress", "maker_token", "taker_token", queryObject)
		if err != nil {
			returnError(w, err, 400)
			return
		}
		query, err = applyOrFilter(query, "trader", "maker", "taker", queryObject)
		if err != nil {
			returnError(w, err, 400)
			return
		}

		pageInt, perPageInt, err := getPages(queryObject)
		if err != nil {
			returnError(w, err, 400)
			return
		}
		query = query.Offset((pageInt - 1) * perPageInt).Limit(perPageInt)
		if query.Error != nil {
			returnError(w, query.Error, 400)
			return
		}
		if queryObject.Get("makerTokenAddress") != "" && queryObject.Get("takerTokenAddress") != "" {
			query := query.Order("price asc, fee_rate asc")
			if query.Error != nil {
				returnError(w, query.Error, 500)
				return
			}
		} else {
			query := query.Order("updated_at")
			if query.Error != nil {
				returnError(w, query.Error, 400)
				return
			}
		}

		orders := []dbModule.Order{}
		if err := query.Find(&orders).Error; err != nil {
			returnError(w, err, 500)
			return
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
			returnError(w, err, 500)
		}
	}
}
