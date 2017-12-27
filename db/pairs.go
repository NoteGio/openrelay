package db

import (
	"github.com/notegio/openrelay/types"
	"github.com/jinzhu/gorm"
	"fmt"
	// "log"
)

// Pair tracks pairs of tokens TokenA and TokenB
type Pair struct {
	TokenA  *types.Address
	TokenB  *types.Address
}

func (pair *Pair) MarshalJSON() ([]byte, error) {
		return []byte(fmt.Sprintf("{\"tokenA\":{\"address\":\"%#x\",\"minAmount\":\"1\",\"maxAmount\":\"115792089237316195423570985008687907853269984665640564039457584007913129639935\",\"precision\":5},\"tokenB\":{\"address\":\"%#x\",\"minAmount\":\"1\",\"maxAmount\":\"115792089237316195423570985008687907853269984665640564039457584007913129639935\",\"precision\":5}}", *pair.TokenA, *pair.TokenB)), nil
}

// GetAllTokenPairs returns an unfilitered list of Pairs based on the trading
// pairs currently present in the database, limited by a count and offset.
func GetAllTokenPairs(db *gorm.DB, offset, count int) ([]Pair, error) {
	tokenPairs := []Pair{}
	if err := db.Model(&Order{}).Select("DISTINCT LEAST(maker_token, taker_token) as token_a, GREATEST(maker_token, taker_token) as token_b").Offset(offset).Limit(count).Scan(&tokenPairs).Error; err != nil {
		return tokenPairs, err
	}
	return tokenPairs, nil
}

// GetTokenAPairs returns a list of Pairs based on the trading pairs currrently
// present in the database, filtered to include only pairs that include tokenA
// and limited by a count and offset.
func GetTokenAPairs(db *gorm.DB, tokenA *types.Address, offset, count int) ([]Pair, error) {
	tokenPairs := []Pair{}
	if err := db.Model(&Order{}).Select("DISTINCT LEAST(maker_token, taker_token) as token_a, GREATEST(maker_token, taker_token) as token_b").Where("taker_token = ? or maker_token = ?", tokenA, tokenA).Offset(offset).Limit(count).Scan(&tokenPairs).Error; err != nil {
		return tokenPairs, err
	}
	return tokenPairs, nil
}

// GetTokenABPairs returns a list of Pairs based on the trading pairs
// currrently present in the database, filtered to include only pairs that
// include both tokenA and tokenB. There should only be one distinct
// combination of both token pairs, so there is no offset or limit, but it
// still returns a list to provide the same return value as the other retrieval
// methods.
func GetTokenABPairs(db *gorm.DB, tokenA, tokenB *types.Address) ([]Pair, error) {
	tokenPairs := []Pair{}
	if err := db.Model(&Order{}).Select("DISTINCT LEAST(maker_token, taker_token) as token_a, GREATEST(maker_token, taker_token) as token_b").Where("(taker_token = ? AND maker_token = ?) or (maker_token = ? and taker_token = ?)", tokenA, tokenB, tokenA, tokenB).Limit(1).Scan(&tokenPairs).Error; err != nil {
		return tokenPairs, err
	}
	return tokenPairs, nil
}
