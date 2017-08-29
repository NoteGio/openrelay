package funds

import (
	"github.com/notegio/0xrelay/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"encoding/hex"
	"log"
)

type OrderValidator interface {
	ValidateOrder(order *types.Order) (bool)
}

type orderValidator struct {
	balanceChecker BalanceChecker
}

func (funds *orderValidator) checkFunds(tokenAddress, userAddress [20]byte, required [32]byte, respond chan bool) {
	requiredInt := new(big.Int)
	requiredInt.SetBytes(required[:])
	balance, err := funds.balanceChecker.GetBalance(tokenAddress, userAddress)
	if err != nil {
		log.Printf("'%v': '%v'", err.Error(), hex.EncodeToString(tokenAddress[:]))
		respond <- false
		return
	}
	respond <- (requiredInt.Cmp(balance) < 0)
}

// ValidateOrder makes sure that the maker of an order has sufficient funds to
// fill the order and pay makerFees
func (funds *orderValidator)ValidateOrder(order *types.Order) bool {
	// TODO: Look this up from somewhere so it can work on different chains
	feeToken := common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")
	makerChan := make(chan bool)
	feeChan := make(chan bool)
	go funds.checkFunds(order.MakerToken, order.Maker, order.MakerTokenAmount, makerChan)
	go funds.checkFunds(feeToken, order.Maker, order.MakerFee, feeChan)
	return (<-makerChan && <-feeChan)
}

func NewRpcOrderValidator(rpcUrl string) (OrderValidator, error){
	if checker, err := NewRpcBalanceChecker(rpcUrl); err != nil {
		return &orderValidator{checker}, nil
	} else {
		return nil, err
	}
}
func NewOrderValidator(checker BalanceChecker) (OrderValidator){
	return &orderValidator{checker}
}
