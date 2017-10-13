package funds

import (
	"encoding/hex"
	"github.com/notegio/openrelay/config"
	"github.com/notegio/openrelay/types"
	"log"
	"math/big"
)

type OrderValidator interface {
	ValidateOrder(order *types.Order) bool
}

type orderValidator struct {
	balanceChecker BalanceChecker
	feeToken       config.FeeToken
	tokenProxy     config.TokenProxy
}

func (funds *orderValidator) checkBalance(tokenAddress, userAddress [20]byte, required []byte, respond chan bool) {
	requiredInt := new(big.Int)
	requiredInt.SetBytes(required[:])
	balance, err := funds.balanceChecker.GetBalance(tokenAddress, userAddress)
	if err != nil {
		log.Printf("'%v': '%v', '%v'", err.Error(), hex.EncodeToString(tokenAddress[:]), hex.EncodeToString(userAddress[:]))
		respond <- false
		return
	}
	respond <- (requiredInt.Cmp(balance) <= 0)
}

func (funds *orderValidator) checkAllowance(tokenAddress, userAddress [20]byte, required []byte, respond chan bool) {
	requiredInt := new(big.Int)
	requiredInt.SetBytes(required[:])
	proxyAddress, err := funds.tokenProxy.Get()
	if err != nil {
		respond <- false
		return
	}
	balance, err := funds.balanceChecker.GetAllowance(tokenAddress, userAddress, proxyAddress)
	if err != nil {
		log.Printf("'%v': '%v', '%v'", err.Error(), hex.EncodeToString(tokenAddress[:]), hex.EncodeToString(userAddress[:]))
		respond <- false
		return
	}
	respond <- (requiredInt.Cmp(balance) <= 0)
}

func getRemainingAmount(numerator, denominator, target []byte) []byte {
	numInt := new(big.Int)
	denomInt := new(big.Int)
	targetInt := new(big.Int)
	numInt.SetBytes(numerator)
	denomInt.SetBytes(denominator)
	targetInt.SetBytes(target)
	mulInt := new(big.Int).Mul(numInt, targetInt)
	return new(big.Int).Sub(targetInt, new(big.Int).Div(mulInt, denomInt)).Bytes()
}

// ValidateOrder makes sure that the maker of an order has sufficient funds to
// fill the order and pay makerFees. This assumes that TakerAmountFilled and
// TakerAmountCancelled reflect
func (funds *orderValidator) ValidateOrder(order *types.Order) bool {
	feeToken, err := funds.feeToken.Get() //common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")
	if err != nil {
		log.Printf("Error getting fee token '%v'", err.Error())
		return false
	}
	makerChan := make(chan bool)
	feeChan := make(chan bool)
	makerAllowanceChan := make(chan bool)
	feeAllowanceChan := make(chan bool)
	unavailableAmount := new(big.Int)
	cancelledAmount := new(big.Int)
	cancelledAmount.SetBytes(order.TakerTokenAmountCancelled[:])
	unavailableAmount.SetBytes(order.TakerTokenAmountFilled[:])
	unavailableAmount.Add(unavailableAmount, cancelledAmount)
	go funds.checkBalance(
		order.MakerToken,
		order.Maker,
		getRemainingAmount(unavailableAmount.Bytes(), order.TakerTokenAmount[:], order.MakerTokenAmount[:]),
		makerChan,
	)
	go funds.checkBalance(
		feeToken,
		order.Maker,
		getRemainingAmount(unavailableAmount.Bytes(), order.TakerTokenAmount[:], order.MakerFee[:]),
		feeChan,
	)
	go funds.checkAllowance(
		order.MakerToken,
		order.Maker,
		getRemainingAmount(unavailableAmount.Bytes(), order.TakerTokenAmount[:], order.MakerTokenAmount[:]),
		makerAllowanceChan,
	)
	go funds.checkAllowance(
		feeToken,
		order.Maker,
		getRemainingAmount(unavailableAmount.Bytes(), order.TakerTokenAmount[:], order.MakerFee[:]),
		feeAllowanceChan,
	)
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
