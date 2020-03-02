package pool


import (
	"bytes"
	"context"
	"math/big"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/common"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"fmt"
	// "log"
)

type FilterContract struct {
	address       *types.Address
	conn          bind.ContractCaller
}


func (fc *FilterContract) Filter(poolID []byte, order *types.Order) (bool, error) {

	// The following code ABIv2 encodes a call to `OrderBookFilter(bytes32 poolid, LibOrder.Order order)`
	// Because the asset data fields are variable length, the order is encoded
	// using offsets and lengths, which mean there's not a single order encoding
	// that can be used across multiple function calls.
	encoding := []byte{}
	encoding = append(encoding, []byte{128, 0, 25, 52}...) // keccak256(signature)[:4]
	if len(poolID) != 32 {
		return false, fmt.Errorf("Expected poolID to be 32 bytes, got %v", len(poolID))
	}
	encoding = append(encoding, poolID...)
	twelveNullBytes := make([]byte, 12)
	orderOffset := make([]byte, 32)
	// The number of bytes up to this point
	orderOffset[31] = 64
	encoding = append(encoding, orderOffset...)
	encoding = append(encoding, twelveNullBytes...)
	encoding = append(encoding, order.Maker[:]...)
	encoding = append(encoding, twelveNullBytes...)
	encoding = append(encoding, order.Taker[:]...)
	encoding = append(encoding, twelveNullBytes...)
	encoding = append(encoding, order.FeeRecipient[:]...)
	encoding = append(encoding, twelveNullBytes...)
	encoding = append(encoding, order.SenderAddress[:]...)
	encoding = append(encoding, order.MakerAssetAmount[:]...)
	encoding = append(encoding, order.TakerAssetAmount[:]...)
	encoding = append(encoding, order.MakerFee[:]...)
	encoding = append(encoding, order.TakerFee[:]...)
	encoding = append(encoding, order.ExpirationTimestampInSec[:]...)
	encoding = append(encoding, order.Salt[:]...)

	// The number of bytes up to this point
	makerAssetDataOffset := 384
	encoding = append(encoding, common.BigToUint256(big.NewInt(int64(makerAssetDataOffset)))[:]...)
	roundUp := 0
	if len(order.MakerAssetData) % 32 > 0 {
		roundUp = 1
	}
	makerAssetDataLength := ((len(order.MakerAssetData) / 32) + roundUp) * 32  // Round up to the next 32 byte chunk

	takerAssetDataOffset := makerAssetDataOffset + 32 + makerAssetDataLength
	encoding = append(encoding, common.BigToUint256(big.NewInt(int64(takerAssetDataOffset)))[:]...)
	encoding = append(encoding, common.BigToUint256(big.NewInt(int64(makerAssetDataLength)))[:]...)
	makerAssetDataBytes := make([]byte, makerAssetDataLength)
	copy(makerAssetDataBytes[:], order.MakerAssetData[:])
	encoding = append(encoding, makerAssetDataBytes...)

	roundUp = 0
	if len(order.TakerAssetData) % 32 > 0 {
		roundUp = 1
	}
	takerAssetDataLength := ((len(order.TakerAssetData) / 32) + roundUp) * 32

	encoding = append(encoding, common.BigToUint256(big.NewInt(int64(takerAssetDataLength)))[:]...)
	takerAssetDataBytes := make([]byte, takerAssetDataLength)
	copy(takerAssetDataBytes[:], order.TakerAssetData[:])
	encoding = append(encoding, takerAssetDataBytes...)


	target := fc.address.ToGethAddress()
	callMsg := ethereum.CallMsg{
		From: order.Maker.ToGethAddress(),
		To: &target,
		Gas: 1000000,
		GasPrice: big.NewInt(1),
		Value: big.NewInt(0),
		Data: encoding,
	}

	result, err := fc.conn.CallContract(context.Background(), callMsg, nil)
	return !bytes.Equal(result[:], make([]byte, 32)), err
}


func NewRPCFilterContract(address *types.Address, rpcURL string) (*FilterContract, error) {
	conn, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	if _, err = conn.SyncProgress(context.Background()); err != nil {
		// This is just here so that an RpcBalanceChecker can't be instantiated
		// successfully if the RPC server isn't responding properly. What RPC
		// function we call isn't important, but SyncProgress is pretty cheap.
		return nil, err
	}
	return NewFilterContract(address, conn), nil

}

func NewFilterContract(address *types.Address, conn bind.ContractCaller) (*FilterContract) {
	return &FilterContract{
		address,
		conn,
	}
}
