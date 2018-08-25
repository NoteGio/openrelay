package config

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/exchangecontract"
	orCommon "github.com/notegio/openrelay/common"
	"log"
)

type FeeToken interface {
	Get(order *types.Order) (types.AssetData, error)
	Set(types.AssetData) error
}

type staticFeeToken struct {
	value types.AssetData
}

func (feeToken *staticFeeToken) Get(order *types.Order) (types.AssetData, error) {
		return feeToken.value, nil
}

func (feeToken *staticFeeToken) Set(assetData types.AssetData) error {
	feeToken.value = assetData
	return nil
}

type rpcFeeToken struct {
	conn bind.ContractBackend
	exchangeTokenMap map[types.Address]types.AssetData
}

func (feeToken *rpcFeeToken) Get(order *types.Order) (types.AssetData, error) {
	feeTokenAssetData := types.AssetData{}
	if feeTokenAssetData, ok := feeToken.exchangeTokenMap[*order.ExchangeAddress]; ok {
		return feeTokenAssetData, nil
	}
	exchange, err := exchangecontract.NewExchange(orCommon.ToGethAddress(order.ExchangeAddress), feeToken.conn)
	if err != nil {
		log.Printf("Error intializing exchange contract '%v': '%v'", hex.EncodeToString(order.ExchangeAddress[:]), err.Error())
		return feeTokenAssetData, err
	}
	feeTokenAssetDataBytes, err := exchange.ZRX_ASSET_DATA(nil)
	if err != nil {
		log.Printf("Error getting fee token address for exchange %#x", order.ExchangeAddress)
		return nil, err
	}
	feeTokenAssetData = make(types.AssetData, len(feeTokenAssetDataBytes))
	copy(feeTokenAssetData[:], feeTokenAssetDataBytes[:])
	feeToken.exchangeTokenMap[*order.ExchangeAddress] = feeTokenAssetData
	return feeTokenAssetData, nil
}

func (feeToken *rpcFeeToken) Set(value types.AssetData) error {
	// the rpcFeeToken looks up from the RPC server, so we can't actually set
	// the value.
	return nil
}

func NewRpcFeeToken(rpcURL string) (FeeToken, error) {
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
	return &rpcFeeToken{conn, make(map[types.Address]types.AssetData)}, nil
}

func StaticFeeToken(address types.AssetData) FeeToken {
	return &staticFeeToken{address}
}
