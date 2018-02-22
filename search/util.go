package search

import (
	"math/big"
	"github.com/notegio/openrelay/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	urlModule "net/url"
	"time"
)


func getExpTime(queryObject urlModule.Values) (*types.Uint256) {
	var currentTime *big.Int

	if expTimeString := queryObject.Get("_expTime"); expTimeString != "" {
		currentTime, _ = new(big.Int).SetString(expTimeString, 10)
	}
	if currentTime == nil {
		// _expTime was not set, or was not parseable. Use current time.
		currentTime = new(big.Int).SetInt64(time.Now().Unix())
	}
	currentTimeBytes := &types.Uint256{}
	copy(currentTimeBytes[:], abi.U256(currentTime))
	return currentTimeBytes
}
