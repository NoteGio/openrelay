package pool_test

import (
	"context"
	"math/big"
	"testing"
	orCommon "github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/pool"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	// "os"
)

type testContractCaller struct {
	code   []byte
	result bool
	err    error
}

func (conn *testContractCaller) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return conn.code, conn.err
}

func (conn *testContractCaller) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error){
	result := make([]byte, 32)
	if conn.result {
		result[31] = 1
	}
	return result, conn.err
}



func TestFilterOk(t *testing.T) {
	address, _ := orCommon.HexToAddress("0xa31e64ea55b9b6bbb9d6a676738e9a5b23149f84")
	// address, _ := common.HexToAddress("0xb23672f74749bf7916ba6827c64111a4d6de7f11")
	fc := pool.NewFilterContract(address, &testContractCaller{[]byte{}, true, nil})
	order := &types.Order{}
	order.Initialize()
	ok, err := fc.Filter(make([]byte, 32), order)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !ok {
		t.Error("Expected filter to pass")
	}
}
func TestFilterFail(t *testing.T) {
	address, _ := orCommon.HexToAddress("0xa31e64ea55b9b6bbb9d6a676738e9a5b23149f84")
	// address, _ := common.HexToAddress("0xb23672f74749bf7916ba6827c64111a4d6de7f11")
	fc := pool.NewFilterContract(address, &testContractCaller{[]byte{}, false, nil})
	order := &types.Order{}
	order.Initialize()
	ok, err := fc.Filter(make([]byte, 32), order)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ok {
		t.Error("Expected filter to fail")
	}
}
