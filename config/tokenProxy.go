package config

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/notegio/openrelay/types"
	orCommon "github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/exchangecontract"
	"gopkg.in/redis.v3"
	"time"
	"log"
)

type TokenProxy interface {
	Get(order *types.Order) (*types.Address, error)
	Set(*types.Address) error
}

type staticTokenProxy struct {
	value *types.Address
}

func (tokenProxy *staticTokenProxy) Get(order *types.Order) (*types.Address, error) {
	return tokenProxy.value, nil
}

func (tokenProxy *staticTokenProxy) Set(address *types.Address) error {
	tokenProxy.value = address
	return nil
}

type rpcTokenProxy struct {
	conn bind.ContractBackend
	exchangeProxyMap map[types.Address]*types.Address
}

func (tokenProxy *rpcTokenProxy) Get(order *types.Order) (*types.Address, error) {
	tokenProxyAddress := &types.Address{}
	if tokenProxyAddress, ok := tokenProxy.exchangeProxyMap[*order.ExchangeAddress]; ok {
		return tokenProxyAddress, nil
	}
	exchange, err := exchangecontract.NewExchange(orCommon.ToGethAddress(order.ExchangeAddress), tokenProxy.conn)
	if err != nil {
		log.Printf("Error intializing exchange contract '%v': '%v'", hex.EncodeToString(order.ExchangeAddress[:]), err.Error())
		return tokenProxyAddress, err
	}
	tokenProxyGethAddress, err := exchange.TOKEN_TRANSFER_PROXY_CONTRACT(nil)
	if err != nil {
		log.Printf("Error getting token proxy address for exhange %#x", order.ExchangeAddress)
		return nil, err
	}
	copy(tokenProxyAddress[:], tokenProxyGethAddress[:])
	tokenProxy.exchangeProxyMap[*order.ExchangeAddress] = tokenProxyAddress
	return tokenProxyAddress, nil
}

func (tokenProxy *rpcTokenProxy) Set(value *types.Address) error {
	// the rpcTokenProxy looks up from the RPC server, so we can't actually set
	// the value.
	return nil
}

type redisTokenProxy struct {
	redisClient     *redis.Client
	cachedValue     *types.Address
	cacheExpiration int64
}

func (tokenProxy *redisTokenProxy) Get(order *types.Order) (*types.Address, error) {
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

func NewRpcTokenProxy(rpcURL string) (FeeToken, error) {
	conn, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	if _, err = conn.SyncProgress(context.Background()); err != nil {
		// This is just here so that an NewRpcFeeToken can't be instantiated
		// successfully if the RPC server isn't responding properly. What RPC
		// function we call isn't important, but SyncProgress is pretty cheap.
		return nil, err
	}
	return &rpcTokenProxy{conn, make(map[types.Address]*types.Address)}, nil
}

func NewTokenProxy(client *redis.Client) TokenProxy {
	return &redisTokenProxy{client, &types.Address{}, 0}
}

func StaticTokenProxy(address *types.Address) TokenProxy {
	return &staticTokenProxy{address}
}
