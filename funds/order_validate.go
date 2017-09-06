package funds

import (
	"encoding/hex"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/config"
	"log"
	"math/big"
)

type OrderValidator interface {
	ValidateOrder(order *types.Order) bool
}

type orderValidator struct {
	balanceChecker BalanceChecker
	feeToken config.FeeToken
	tokenProxy config.TokenProxy
}

func (funds *orderValidator) checkBalance(tokenAddress, userAddress [20]byte, required [32]byte, respond chan bool) {
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

func (funds *orderValidator) checkAllowance(tokenAddress, userAddress [20]byte, required [32]byte, respond chan bool) {
	requiredInt := new(big.Int)
	requiredInt.SetBytes(required[:])
	proxyAddress, err := funds.tokenProxy.Get()
	if err != nil {
		respond <- false
		return
	}
	balance, err := funds.balanceChecker.GetAllowance(tokenAddress, userAddress, proxyAddress)
	if err != nil {
		log.Printf("'%v': '%v'", err.Error(), hex.EncodeToString(tokenAddress[:]))
		respond <- false
		return
	}
	respond <- (requiredInt.Cmp(balance) <= 0)
}

// ValidateOrder makes sure that the maker of an order has sufficient funds to
// fill the order and pay makerFees
func (funds *orderValidator) ValidateOrder(order *types.Order) bool {
	feeToken, err := funds.feeToken.Get() //common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")
	if err != nil { return false }
	makerChan := make(chan bool)
	feeChan := make(chan bool)
	makerAllowanceChan := make(chan bool)
	feeAllowanceChan := make(chan bool)
	go funds.checkBalance(order.MakerToken, order.Maker, order.MakerTokenAmount, makerChan)
	go funds.checkBalance(feeToken, order.Maker, order.MakerFee, feeChan)
	go funds.checkAllowance(order.MakerToken, order.Maker, order.MakerTokenAmount, makerAllowanceChan)
	go funds.checkAllowance(feeToken, order.Maker, order.MakerFee, feeAllowanceChan)
	result := true
	if !<-makerChan {
		log.Printf("Insufficient maker token funds")
		result = false
	}
	if !<-feeChan {
		log.Printf("Insufficient fee token funds")
		result = false
	}
	if !<-makerAllowanceChan {
		log.Printf("Insufficient maker token allowance")
		result = false
	}
	if !<-feeAllowanceChan {
		log.Printf("Insufficient fee token allowance")
		result = false
	}
	return result
}

func NewRpcOrderValidator(rpcUrl string, feeToken config.FeeToken, tokenProxy config.TokenProxy) (OrderValidator, error) {
	if checker, err := NewRpcBalanceChecker(rpcUrl); err == nil {
		return &orderValidator{checker, feeToken, tokenProxy}, nil
	} else {
		return nil, err
	}
}
func NewOrderValidator(checker BalanceChecker, feeToken config.FeeToken, tokenProxy config.TokenProxy) OrderValidator {
	return &orderValidator{checker, feeToken, tokenProxy}
}
