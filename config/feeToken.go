package config

import (
	"encoding/hex"
	"github.com/notegio/openrelay/types"
	"gopkg.in/redis.v3"
	"time"
)

type FeeToken interface {
	Get() (*types.Address, error)
	Set(*types.Address) error
}

type staticFeeToken struct {
	value *types.Address
}

func (feeToken *staticFeeToken) Get() (*types.Address, error) {
	return feeToken.value, nil
}

func (feeToken *staticFeeToken) Set(address *types.Address) error {
	feeToken.value = address
	return nil
}

type redisFeeToken struct {
	redisClient     *redis.Client
	cachedValue     *types.Address
	cacheExpiration int64
}

func (feeToken *redisFeeToken) Get() (*types.Address, error) {
	if feeToken.cacheExpiration > time.Now().Unix() {
		// The fee token is unlikely to change, so caching it should be fine.
		// Doesn't hurt to check periodically just in case though.
		return feeToken.cachedValue, nil
	}
	result := &types.Address{}
	val, err := feeToken.redisClient.Get("feeToken::address").Result()
	if err != nil {
		return result, err
	}
	addressSlice, err := hex.DecodeString(val)
	if err != nil {
		return result, err
	}
	copy(result[:], addressSlice[:])
	return result, nil
}

func (feeToken *redisFeeToken) Set(value *types.Address) error {
	return feeToken.redisClient.Set("feeToken::address", hex.EncodeToString(value[:]), 0).Err()
}

func NewFeeToken(client *redis.Client) FeeToken {
	return &redisFeeToken{client, &types.Address{}, 0}
}

func StaticFeeToken(address *types.Address) FeeToken {
	return &staticFeeToken{address}
}
