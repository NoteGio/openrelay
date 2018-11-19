package search

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/common"
	dbModule "github.com/notegio/openrelay/db"
	"net/http"
	"strconv"
)

func PairHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queryObject := r.URL.Query()
		tokenAString := queryObject.Get("assetDataA")
		tokenBString := queryObject.Get("assetDataB")
		networkID, err := strconv.Atoi(queryObject.Get("networkId"))
		if err != nil {
			networkID = 1
		}
		if tokenAString == "" && tokenBString != "" {
			tokenAString, tokenBString = tokenBString, ""
		}
		pageInt, perPageInt, err := getPages(queryObject)
		offset := (pageInt - 1) * perPageInt
		if err != nil {
			returnError(w, err, 400)
			return
		}
		var pairs []dbModule.Pair
		var count int
		if tokenAString == "" {
			pairs, count, err = dbModule.GetAllTokenPairs(db, offset, perPageInt, networkID)
			if err != nil {
				returnError(w, err, 400)
				return
			}
		} else {
			assetDataA, err := common.HexToAssetData(tokenAString)
			if err != nil {
				returnError(w, err, 400)
				return
			}
			if tokenBString == "" {
				pairs, count, err = dbModule.GetTokenAPairs(db, assetDataA, offset, perPageInt, networkID)
			} else {
				assetDataB, err := common.HexToAssetData(tokenBString)
				if err != nil {
					returnError(w, err, 400)
					return
				}
				pairs, count, err = dbModule.GetTokenABPairs(db, assetDataA, assetDataB, networkID)
			}
			if err != nil {
				returnError(w, err, 400)
				return
			}
		}
		response, err := json.Marshal(GetPagedResult(count, pageInt, perPageInt, pairs))
		if err != nil {
			returnError(w, err, 500)
			return
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
