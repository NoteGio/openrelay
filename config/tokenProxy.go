package config

import (
	"encoding/hex"
	"github.com/notegio/openrelay/types"
	"gopkg.in/redis.v3"
	"time"
)

type TokenProxy interface {
	Get() (*types.Address, error)
	Set(*types.Address) error
}

type staticTokenProxy struct {
	value *types.Address
}

func (tokenProxy *staticTokenProxy) Get() (*types.Address, error) {
	return tokenProxy.value, nil
}

func (tokenProxy *staticTokenProxy) Set(address *types.Address) error {
	tokenProxy.value = address
	return nil
}

type redisTokenProxy struct {
	redisClient     *redis.Client
	cachedValue     *types.Address
	cacheExpiration int64
}

func (tokenProxy *redisTokenProxy) Get() (*types.Address, error) {
	if tokenProxy.cacheExpiration > time.Now().Unix() {
		// The token proxy shouldn't change often, but it doesn't hurt to check
		// periodically.
		return tokenProxy.cachedValue, nil
	}
	result := &types.Address{}
	val, err := tokenProxy.redisClient.Get("tokenProxy::address").Result()
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

func (tokenProxy *redisTokenProxy) Set(value *types.Address) error {
	return tokenProxy.redisClient.Set("tokenProxy::address", hex.EncodeToString(value[:]), 0).Err()
}

func NewTokenProxy(client *redis.Client) TokenProxy {
	return &redisTokenProxy{client, &types.Address{}, 0}
}

func StaticTokenProxy(address *types.Address) TokenProxy {
	return &staticTokenProxy{address}
}
