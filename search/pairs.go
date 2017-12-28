package search

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/common"
	dbModule "github.com/notegio/openrelay/db"
	"net/http"
)

func PairHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queryObject := r.URL.Query()
		tokenAString := queryObject.Get("tokenA")
		tokenBString := queryObject.Get("tokenB")
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
		if tokenAString == "" {
			pairs, err = dbModule.GetAllTokenPairs(db, offset, perPageInt)
			if err != nil {
				returnError(w, err, 400)
				return
			}
		} else {
			tokenABytes, err := common.HexToBytes(tokenAString)
			if err != nil {
				returnError(w, err, 400)
				return
			}
			tokenAAddress := common.BytesToOrAddress(tokenABytes)
			if tokenBString == "" {
				pairs, err = dbModule.GetTokenAPairs(db, tokenAAddress, offset, perPageInt)
			} else {
				tokenBBytes, err := common.HexToBytes(tokenBString)
				if err != nil {
					returnError(w, err, 400)
					return
				}
				tokenBAddress := common.BytesToOrAddress(tokenBBytes)
				pairs, err = dbModule.GetTokenABPairs(db, tokenAAddress, tokenBAddress)
			}
			if err != nil {
				returnError(w, err, 400)
				return
			}
		}
		response, err := json.Marshal(pairs)
		if err != nil {
			returnError(w, err, 500)
			return
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
