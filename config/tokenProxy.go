package config

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/notegio/openrelay/types"
	orCommon "github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/exchangecontract"
	"log"
)

type TokenProxy interface {
	Get(order *types.Order) (*types.Address, error)
	GetById(order *types.Order, proxyID [4]byte) (*types.Address, error)
	Set(*types.Address) error
}

type staticTokenProxy struct {
	value *types.Address
}

func (tokenProxy *staticTokenProxy) Get(order *types.Order) (*types.Address, error) {
	return tokenProxy.value, nil
}

func (tokenProxy *staticTokenProxy) GetById(order *types.Order, proxyID [4]byte) (*types.Address, error) {
	return tokenProxy.value, nil
}

func (tokenProxy *staticTokenProxy) Set(address *types.Address) error {
	tokenProxy.value = address
	return nil
}

type rpcTokenProxy struct {
	conn bind.ContractBackend
	exchangeProxyMap map[types.Address]map[[4]byte]*types.Address
}

func (tokenProxy *rpcTokenProxy) Get(order *types.Order) (*types.Address, error) {
	return tokenProxy.GetById(order, order.MakerAssetData.ProxyId())
}
func (tokenProxy *rpcTokenProxy) GetById(order *types.Order, proxyID [4]byte) (*types.Address, error) {
	tokenProxyAddress := &types.Address{}
	if tokenProxyAddress, ok := tokenProxy.exchangeProxyMap[*order.ExchangeAddress][proxyID]; ok {
		return tokenProxyAddress, nil
	}
	exchange, err := exchangecontract.NewExchangecontract(orCommon.ToGethAddress(order.ExchangeAddress), tokenProxy.conn)
	if err != nil {
		log.Printf("Error intializing exchange contract '%v': '%v'", hex.EncodeToString(order.ExchangeAddress[:]), err.Error())
		return tokenProxyAddress, err
	}
	tokenProxyGethAddress, err := exchange.GetAssetProxy(nil, proxyID)
	if err != nil {
		log.Printf("Error getting token proxy address for exhange %v", order.ExchangeAddress[:])
		return nil, err
	}
	copy(tokenProxyAddress[:], tokenProxyGethAddress[:])
	_, ok := tokenProxy.exchangeProxyMap[*order.ExchangeAddress]
	if !ok {
		tokenProxy.exchangeProxyMap[*order.ExchangeAddress] = make(map[[4]byte]*types.Address)
	}
	tokenProxy.exchangeProxyMap[*order.ExchangeAddress][proxyID] = tokenProxyAddress
	return tokenProxyAddress, nil
}

func (tokenProxy *rpcTokenProxy) Set(value *types.Address) error {
	// the rpcTokenProxy looks up from the RPC server, so we can't actually set
	// the value.
	return nil
}

func NewRpcTokenProxy(rpcURL string) (TokenProxy, error) {
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
	return &rpcTokenProxy{conn, make(map[types.Address]map[[4]byte]*types.Address)}, nil
}

func StaticTokenProxy(address *types.Address) TokenProxy {
	return &staticTokenProxy{address}
}
