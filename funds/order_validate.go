package funds

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/notegio/openrelay/config"
	"github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/funds/balance"
	"log"
	"math/big"
)

type OrderValidator interface {
	ValidateOrder(order *types.Order) (bool, error)
}

type orderValidator struct {
	balanceChecker balance.BalanceChecker
	feeToken       config.FeeToken
	tokenProxy     config.TokenProxy
}

type boolOrErr struct {
	success bool
	err     error
}

func (funds *orderValidator) checkBalance(assetData types.AssetData, userAddress *types.Address, required []byte, respond chan boolOrErr) {
	requiredInt := new(big.Int)
	requiredInt.SetBytes(required[:])
	if requiredInt.Cmp(big.NewInt(0)) == 0 {
		// If the required amount is 0, there's no point in looking up the balance
		respond <- boolOrErr{true, nil}
		return
	}
	balance, err := funds.balanceChecker.GetBalance(assetData, userAddress)
	if err != nil {
		log.Printf("'%v': '%v', '%v'", err.Error(), hex.EncodeToString(assetData[:]), hex.EncodeToString(userAddress[:]))
		respond <- boolOrErr{false, err}
		return
	}
	respond <- boolOrErr{(requiredInt.Cmp(balance) <= 0), nil}
}

func (funds *orderValidator) checkAllowance(assetData types.AssetData, userAddress, proxyAddress *types.Address, required []byte, respond chan boolOrErr) {
	requiredInt := new(big.Int)
	requiredInt.SetBytes(required[:])
	if requiredInt.Cmp(big.NewInt(0)) == 0 {
		// If the required amount is 0, there's no point in looking up the allowance
		respond <- boolOrErr{true, nil}
		return
	}
	balance, err := funds.balanceChecker.GetAllowance(assetData, userAddress, proxyAddress)
	if err != nil {
		log.Printf("'%v': '%v', '%v'", err.Error(), hex.EncodeToString(assetData[:]), hex.EncodeToString(userAddress[:]))
		respond <- boolOrErr{false, err}
		return
	}
	respond <- boolOrErr{(requiredInt.Cmp(balance) <= 0), nil}
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
func (funds *orderValidator) ValidateOrder(order *types.Order) (bool, error) {
	feeToken := order.MakerFeeAssetData //funds.feeToken.Get(order)
	makerProxyAddress, err := funds.tokenProxy.Get(order)
	if err != nil {
		log.Printf("Error getting token proxy address '%v'", err.Error())
		return false, err
	}
	feeProxyAddress, err := funds.tokenProxy.GetById(order, order.MakerFeeAssetData.ProxyId())
	if err != nil {
		log.Printf("Error getting fee token proxy address '%v'", err.Error())
		return false, err
	}
	makerChan := make(chan boolOrErr)
	feeChan := make(chan boolOrErr)
	makerAllowanceChan := make(chan boolOrErr)
	feeAllowanceChan := make(chan boolOrErr)
	unavailableAmount := order.TakerAssetAmountFilled.Big()

	if bytes.Equal(order.MakerAssetData[:], feeToken[:]) {
		requiredMakerAmount := common.BigToUint256(big.NewInt(0).Add(order.MakerAssetAmount.Big(), order.MakerFee.Big()))
		go funds.checkBalance(
			order.MakerAssetData,
			order.Maker,
			getRemainingAmount(unavailableAmount.Bytes(), order.TakerAssetAmount[:], requiredMakerAmount[:]),
			makerChan,
		)
		go funds.checkAllowance(
			order.MakerAssetData,
			order.Maker,
			makerProxyAddress,
			getRemainingAmount(unavailableAmount.Bytes(), order.TakerAssetAmount[:], requiredMakerAmount[:]),
			makerAllowanceChan,
		)
		// If the MakerAssetData == MakerFeeAssetData, then we only need to check it
		// once with the combined amount, but the other channels still expect
		// results
		go func () {
			feeChan <- boolOrErr{success: true}
			feeAllowanceChan <- boolOrErr{success: true}
		}()
	} else {
		go funds.checkBalance(
			order.MakerAssetData,
			order.Maker,
			getRemainingAmount(unavailableAmount.Bytes(), order.TakerAssetAmount[:], order.MakerAssetAmount[:]),
			makerChan,
		)
		go funds.checkAllowance(
			order.MakerAssetData,
			order.Maker,
			makerProxyAddress,
			getRemainingAmount(unavailableAmount.Bytes(), order.TakerAssetAmount[:], order.MakerAssetAmount[:]),
			makerAllowanceChan,
		)
		if !bytes.Equal(feeToken[:], order.TakerAssetData[:]) && order.TakerAssetAmount.Big().Cmp(order.MakerFee.Big()) > 0 {
			go funds.checkBalance(
				feeToken,
				order.Maker,
				getRemainingAmount(unavailableAmount.Bytes(), order.TakerAssetAmount[:], order.MakerFee[:]),
				feeChan,
			)
		} else {
			// If the TakerToken is the FeeToken and the TakerAssetAmount is bigger
			// than the MakerFee, the maker will always have enough tokens to pay the
			// fee after the trade.
			go func() {
				feeChan <- boolOrErr{success: true}
			}()
		}
		go funds.checkAllowance(
			feeToken,
			order.Maker,
			feeProxyAddress,
			getRemainingAmount(unavailableAmount.Bytes(), order.TakerAssetAmount[:], order.MakerFee[:]),
			feeAllowanceChan,
		)
	}
	result := true
	if chanResult := <-makerChan; !chanResult.success {
		log.Printf("Insufficient maker token funds")
		if chanResult.err != nil {
			if chanResult.err.Error() == "no contract code at given address" {
				return false, chanResult.err
			}
			panic(fmt.Sprintf("RPC Communication Failed: '%v'", chanResult.err.Error()))
		}
		result = false
	}
	if chanResult := <-feeChan; !chanResult.success {
		log.Printf("Insufficient fee token allowance")
		if chanResult.err != nil {
			if chanResult.err.Error() == "no contract code at given address" {
				return false, chanResult.err
			}
			panic(fmt.Sprintf("RPC Communication Failed: '%v'", chanResult.err.Error()))
		}
		result = false
	}
	if chanResult := <-makerAllowanceChan; !chanResult.success {
		log.Printf("Insufficient makers token allowance")
		if chanResult.err != nil {
			if chanResult.err.Error() == "no contract code at given address" {
				return false, chanResult.err
			}
			panic(fmt.Sprintf("RPC Communication Failed: '%v'", chanResult.err.Error()))
		}
		result = false
	}
	if chanResult := <-feeAllowanceChan; !chanResult.success {
		log.Printf("Insufficient fee token allowance")
		if chanResult.err != nil {
			if chanResult.err.Error() == "no contract code at given address" {
				return false, chanResult.err
			}
			panic(fmt.Sprintf("RPC Communication Failed: '%v'", chanResult.err.Error()))
		}
		result = false
	}
	return result, nil
}

func NewRpcOrderValidator(rpcUrl string, feeToken config.FeeToken, tokenProxy config.TokenProxy, invalidationChannel channels.ConsumerChannel) (OrderValidator, error) {
	if checker, err := balance.NewRpcRoutingBalanceChecker(rpcUrl); err == nil {
		if invalidationChannel != nil {
			invalidationChannel.AddConsumer(checker)
			invalidationChannel.StartConsuming()
		}
		return &orderValidator{checker, feeToken, tokenProxy}, nil
	} else {
		return nil, err
	}
}
func NewOrderValidator(checker balance.BalanceChecker, feeToken config.FeeToken, tokenProxy config.TokenProxy) OrderValidator {
	return &orderValidator{checker, feeToken, tokenProxy}
}
