package config

import (
	"gopkg.in/redis.v3"
	"math/big"
	"time"
)

const (
	cacheDuration = 60 * 5
)

type BaseFee interface {
	Get() (*big.Int, error)
	Set(*big.Int) error
}

type redisBaseFee struct {
	redisClient     *redis.Client
	cachedValue     *big.Int
	cacheExpiration int64
}

func (baseFee *redisBaseFee) Get() (*big.Int, error) {
	if baseFee.cacheExpiration > time.Now().Unix() {
		// It should be okay to cache the base fee for a while. It might mean
		// that we require the old fee until the cache expires, but that
		// shouldn't be a big deal
		return baseFee.cachedValue, nil
	}
	val, err := baseFee.redisClient.Get("baseFee").Result()
	if err != nil {
		return nil, err
	}
	result := new(big.Int)
	result.SetBytes([]byte(val))
	return result, nil
}

func (baseFee *redisBaseFee) Set(value *big.Int) error {
	return baseFee.redisClient.Set("baseFee", string(value.Bytes()), 0).Err()
}

func NewBaseFee(client *redis.Client) BaseFee {
	return &redisBaseFee{client, nil, 0}
}
