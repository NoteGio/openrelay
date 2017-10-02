package funds

import (
	"math/big"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/notegio/openrelay/types"
	orCommon "github.com/notegio/openrelay/common"
  "github.com/notegio/openrelay/exchangecontract"
)

type FilledLookup interface {
	GetAmountCancelled(order *types.Order) ([32]byte, error)
	GetAmountFilled(order *types.Order) ([32]byte, error)
}

type rpcFilledLookup struct {
	conn bind.ContractBackend
}

func (filled *rpcFilledLookup) GetAmountCancelled(order *types.Order) ([32]byte, error) {
	cancelledAmount := [32]byte{}
	exchange, err := exchangecontract.NewExchange(orCommon.BytesToAddress(order.ExchangeAddress), filled.conn)
	if err != nil {
		return cancelledAmount, err
	}
	hash := [32]byte{}
	copy(hash[:], order.Hash())
	amount, err := exchange.Cancelled(nil, hash)
	if err != nil {
		return cancelledAmount, err
	}
	cancelledSlice := common.LeftPadBytes(amount.Bytes(), 32)
	copy(cancelledAmount[:], cancelledSlice)
	return cancelledAmount, nil
}

func (filled *rpcFilledLookup) GetAmountFilled(order *types.Order) ([32]byte, error) {
	filledAmount := [32]byte{}
	exchange, err := exchangecontract.NewExchange(orCommon.BytesToAddress(order.ExchangeAddress), filled.conn)
	if err != nil {
		return filledAmount, err
	}
	hash := [32]byte{}
	copy(hash[:], order.Hash())
	amount, err := exchange.Filled(nil, hash)
	if err != nil {
		return filledAmount, err
	}
	cancelledSlice := common.LeftPadBytes(amount.Bytes(), 32)
	copy(filledAmount[:], cancelledSlice)
	return filledAmount, nil
}
// TODO: Test FilledChecker
// TODO: Make FundCheckRelay update order with TakerTokenAmountFilled

func NewRpcFilledLookup(rpcUrl string) (FilledLookup, error) {
	conn, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}
	return &rpcFilledLookup{conn}, nil
}

type MockFilledLookup struct {
	cancelled *big.Int
	filled *big.Int
	err error
}
func (filled *MockFilledLookup) GetAmountCancelled(order *types.Order) ([32]byte, error) {
	result := [32]byte{}
	if filled.err != nil {
		return [32]byte{}, filled.err
	}
	filledSlice := common.LeftPadBytes(filled.cancelled.Bytes(), 32)
	copy(result[:], filledSlice)
	return result, nil
}
func (filled *MockFilledLookup) GetAmountFilled(order *types.Order) ([32]byte, error) {
	result := [32]byte{}
	if filled.err != nil {
		return [32]byte{}, filled.err
	}
	filledSlice := common.LeftPadBytes(filled.filled.Bytes(), 32)
	copy(result[:], filledSlice)
	return result, nil
}

func NewMockFilledLookup(cancelled, filled string, err error) FilledLookup{
	cancelledInt := new(big.Int)
	cancelledInt.SetString(cancelled, 10)
	filledInt := new(big.Int)
	filledInt.SetString(filled, 10)
	return &MockFilledLookup{cancelledInt, filledInt, err}
}
