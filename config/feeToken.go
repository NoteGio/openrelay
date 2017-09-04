package config

import (
	"gopkg.in/redis.v3"
	"time"
	"encoding/hex"
)

type FeeToken interface {
	Get() ([20]byte, error)
	Set([20]byte) error
}

type staticFeeToken struct {
	value [20]byte
}

func (feeToken *staticFeeToken) Get() ([20]byte, error) {
	return feeToken.value, nil
}

func (feeToken *staticFeeToken) Set(address [20]byte) error {
	feeToken.value = address
	return nil
}

type redisFeeToken struct {
	redisClient     *redis.Client
	cachedValue     [20]byte
	cacheExpiration int64
}

func (feeToken *redisFeeToken) Get() ([20]byte, error) {
	if feeToken.cacheExpiration > time.Now().Unix() {
		// The fee token is unlikely to change, so caching it should be fine.
		// Doesn't hurt to check periodically just in case though.
		return feeToken.cachedValue, nil
	}
	result := [20]byte{}
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

func (feeToken *redisFeeToken) Set(value [20]byte) error {
	return feeToken.redisClient.Set("feeToken::address", hex.EncodeToString(value[:]), 0).Err()
}

func NewFeeToken(client *redis.Client) FeeToken {
	return &redisFeeToken{client, [20]byte{}, 0}
}

func StaticFeeToken(address [20]byte) FeeToken {
	return &staticFeeToken{address}
}
