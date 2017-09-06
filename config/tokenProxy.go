package config

import (
	"encoding/hex"
	"gopkg.in/redis.v3"
	"time"
)

type TokenProxy interface {
	Get() ([20]byte, error)
	Set([20]byte) error
}

type staticTokenProxy struct {
	value [20]byte
}

func (tokenProxy *staticTokenProxy) Get() ([20]byte, error) {
	return tokenProxy.value, nil
}

func (tokenProxy *staticTokenProxy) Set(address [20]byte) error {
	tokenProxy.value = address
	return nil
}

type redisTokenProxy struct {
	redisClient     *redis.Client
	cachedValue     [20]byte
	cacheExpiration int64
}

func (tokenProxy *redisTokenProxy) Get() ([20]byte, error) {
	if tokenProxy.cacheExpiration > time.Now().Unix() {
		// The token proxy shouldn't change often, but it doesn't hurt to check
		// periodically.
		return tokenProxy.cachedValue, nil
	}
	result := [20]byte{}
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

func (tokenProxy *redisTokenProxy) Set(value [20]byte) error {
	return tokenProxy.redisClient.Set("tokenProxy::address", hex.EncodeToString(value[:]), 0).Err()
}

func NewTokenProxy(client *redis.Client) TokenProxy {
	return &redisTokenProxy{client, [20]byte{}, 0}
}

func StaticTokenProxy(address [20]byte) TokenProxy {
	return &staticTokenProxy{address}
}
