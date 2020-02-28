package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/types"
	// "log"
)

// Pair tracks pairs of tokens TokenA and TokenB
type Pair struct {
	TokenA types.AssetData
	TokenB types.AssetData
}

func (pair *Pair) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{\"assetDataA\":{\"assetData\":\"%#x\",\"minAmount\":\"1\",\"maxAmount\":\"115792089237316195423570985008687907853269984665640564039457584007913129639935\",\"precision\":5},\"assetDataB\":{\"assetData\":\"%#x\",\"minAmount\":\"1\",\"maxAmount\":\"115792089237316195423570985008687907853269984665640564039457584007913129639935\",\"precision\":5}}", pair.TokenA, pair.TokenB)), nil
}

// GetAllTokenPairs returns an unfilitered list of Pairs based on the trading
// pairs currently present in the database, limited by a count and offset.
func GetAllTokenPairs(db *gorm.DB, offset, count, networkID int) ([]Pair, int, error) {
	tokenPairs := []Pair{}
	var total int
	// This uses a subquery, as `DISTINCT maker_asset_data, taker_asset_data` can be
	// determined easily based on indexes, but includes duplicate token pairs
	// showing both (A, B) and (B, A). Once we've done that, we reduce duplicates
	// by getting min(A, B), max(A, B).
	//
	// The results would be the same if we queried the orders table directly
	// instead of doing a subquery, but indexes would not be used, and the query
	// would be very inefficient.
	if err := db.Raw("SELECT COUNT(*) FROM (SELECT DISTINCT LEAST(x.maker_asset_data, x.taker_asset_data), GREATEST(x.maker_asset_data, x.taker_asset_data) from (SELECT DISTINCT maker_asset_data, taker_asset_data from orderv3 WHERE exchange_address IN (SELECT address FROM exchanges WHERE network = ?)) as x) as y", networkID).Count(&total).Error; err != nil {
		return tokenPairs, total, err
	}
	if err := db.Raw("SELECT DISTINCT LEAST(x.maker_asset_data, x.taker_asset_data) as token_a, GREATEST(x.maker_asset_data, x.taker_asset_data) as token_b from (SELECT DISTINCT maker_asset_data, taker_asset_data from orderv3 WHERE exchange_address IN (SELECT address FROM exchanges WHERE network = ?)) as x", networkID).Offset(offset).Limit(count).Scan(&tokenPairs).Error; err != nil {
		return tokenPairs, total, err
	}
	return tokenPairs, total, nil
}

// GetTokenAPairs returns a list of Pairs based on the trading pairs currrently
// present in the database, filtered to include only pairs that include tokenA
// and limited by a count and offset.
func GetTokenAPairs(db *gorm.DB, tokenA types.AssetData, offset, count, networkID int) ([]Pair, int, error) {
	tokenPairs := []Pair{}
	var total int
	if err := db.Raw("SELECT COUNT(*) FROM (SELECT DISTINCT LEAST(x.maker_asset_data, x.taker_asset_data), GREATEST(x.maker_asset_data, x.taker_asset_data) from (SELECT DISTINCT maker_asset_data, taker_asset_data from orderv3 WHERE exchange_address IN (SELECT address FROM exchanges WHERE network = ?)) as x WHERE x.taker_asset_data = ? or x.maker_asset_data = ?) as y", networkID, []byte(tokenA[:]), []byte(tokenA[:])).Count(&total).Error; err != nil {
		return tokenPairs, total, err
	}
	if err := db.Raw("SELECT DISTINCT LEAST(x.maker_asset_data, x.taker_asset_data) as token_a, GREATEST(x.maker_asset_data, x.taker_asset_data) as token_b from (SELECT DISTINCT maker_asset_data, taker_asset_data from orderv3 WHERE exchange_address IN (SELECT address FROM exchanges WHERE network = ?)) as x WHERE x.taker_asset_data = ? or x.maker_asset_data = ?", networkID, []byte(tokenA[:]), []byte(tokenA[:])).Offset(offset).Limit(count).Scan(&tokenPairs).Error; err != nil {
		return tokenPairs, total, err
	}
	return tokenPairs, total, nil
}

// GetTokenABPairs returns a list of Pairs based on the trading pairs
// currrently present in the database, filtered to include only pairs that
// include both tokenA and tokenB. There should only be one distinct
// combination of both token pairs, so there is no offset or limit, but it
// still returns a list to provide the same return value as the other retrieval
// methods.
func GetTokenABPairs(db *gorm.DB, tokenA, tokenB types.AssetData, networkID int) ([]Pair, int, error) {
	tokenPairs := []Pair{}
	if err := db.Raw("SELECT DISTINCT LEAST(x.maker_asset_data, x.taker_asset_data) as token_a, GREATEST(x.maker_asset_data, x.taker_asset_data) as token_b from (SELECT DISTINCT maker_asset_data, taker_asset_data from orderv3 WHERE exchange_address IN (SELECT address FROM exchanges WHERE network = ?)) as x WHERE (x.taker_asset_data = ? AND x.maker_asset_data = ?) or (x.maker_asset_data = ? and x.taker_asset_data = ?)", networkID, []byte(tokenA[:]), []byte(tokenB[:]), []byte(tokenA[:]), []byte(tokenB[:])).Scan(&tokenPairs).Error; err != nil {
		return tokenPairs, 0, err
	}
	return tokenPairs, len(tokenPairs), nil
}
