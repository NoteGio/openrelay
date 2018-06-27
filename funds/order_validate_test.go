package funds_test

import (
	"encoding/json"
	"errors"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/config"
	"github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/types"
	"math/big"
	"reflect"
	"testing"
	// "log"
)

func hexToAddress(addressHex string) (*types.Address, error) {
	addressBytes, err := common.HexToBytes(addressHex)
	if err != nil {
		return nil, err
	}
	address := &types.Address{}
	copy(address[:], addressBytes[:])
	return address, nil
}

func createMockBalanceChecker(tokenAddressHex, userAddressHex string, tokenAmount string, feeTokenAmount string, t *testing.T) funds.BalanceChecker {
	tokenAddress, err := hexToAddress(tokenAddressHex)
	if err != nil {
		t.Errorf(err.Error())
	}
	userAddress, err := hexToAddress(userAddressHex)
	if err != nil {
		t.Errorf(err.Error())
	}
	feeTokenAddress, _ := hexToAddress("e41d2489571d322189246dafa5ebde1f4699f498")
	tokenInt := new(big.Int)
	tokenInt.SetString(tokenAmount, 10)
	feeTokenInt := new(big.Int)
	feeTokenInt.SetString(feeTokenAmount, 10)
	balanceMap := make(map[types.Address]map[types.Address]*big.Int)
	balanceMap[*tokenAddress] = make(map[types.Address]*big.Int)
	balanceMap[*feeTokenAddress] = make(map[types.Address]*big.Int)
	balanceMap[*tokenAddress][*userAddress] = tokenInt
	if reflect.DeepEqual(feeTokenAddress, tokenAddress) {
		balanceMap[*feeTokenAddress] = balanceMap[*tokenAddress]
	}
	balanceMap[*feeTokenAddress][*userAddress] = feeTokenInt
	return funds.NewMockBalanceChecker(balanceMap)
}

func TestOrderValidate(t *testing.T) {
	balanceChecker := createMockBalanceChecker("1dad4783cf3fe3085c1426157ab175a6119a04ba", "627306090abab3a6e1400e9345bc60c78a8bef57", "0", "0", t)
	feeTokenAddress, _ := hexToAddress("e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAddress), config.StaticTokenProxy(tokenProxyAddress))
	testOrderBytes := getTestOrderBytes()
	newOrder, err := types.OrderFromBytes(testOrderBytes)
	if err != nil {
		t.Errorf("Error parsing order: %v", err.Error())
	}
	if result, _ := validator.ValidateOrder(newOrder); result {
		t.Errorf("Expected insufficient funds")
	}
}

func TestOrderValidateSufficient(t *testing.T) {
	balanceChecker := createMockBalanceChecker("1dad4783cf3fe3085c1426157ab175a6119a04ba", "627306090abab3a6e1400e9345bc60c78a8bef57", "50000000000000000000", "0", t)
	feeTokenAddress, _ := hexToAddress("e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAddress), config.StaticTokenProxy(tokenProxyAddress))
	testOrderBytes := getTestOrderBytes()
	newOrder, err := types.OrderFromBytes(testOrderBytes)
	if err != nil {
		t.Errorf("Error parsing order: %v", err.Error())
	}

	if result, _ := validator.ValidateOrder(newOrder); !result {
		t.Errorf("Expected sufficient funds")
	}
}

func TestOrderValidateSpent(t *testing.T) {
	// 25000000000000000000 is half of the makerAssetAmount for the sample order
	balanceChecker := createMockBalanceChecker("1dad4783cf3fe3085c1426157ab175a6119a04ba", "627306090abab3a6e1400e9345bc60c78a8bef57", "25000000000000000000", "0", t)
	feeTokenAddress, _ := hexToAddress("e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAddress), config.StaticTokenProxy(tokenProxyAddress))
	testOrderBytes := getTestOrderBytes()
	newOrder, err := types.OrderFromBytes(testOrderBytes)
	if err != nil {
		t.Errorf("Error parsing order: %v", err.Error())
	}
	filledInt := new(big.Int)
	// Half the taker token amount for the order
	filledInt.SetString("500000000000000000", 10)
	copy(newOrder.TakerAssetAmountFilled[:], gethCommon.LeftPadBytes(filledInt.Bytes(), 32))

	if result, _ := validator.ValidateOrder(newOrder); !result {
		t.Errorf("Expected sufficient funds")
	}
}

func TestErrorPanic(t *testing.T) {
	var x interface{}
	err := json.Unmarshal([]byte("<oops>"), x)
	balanceChecker := funds.NewErrorMockBalanceChecker(err)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic from JSON error")
		}
	}()
	feeTokenAddress, _ := hexToAddress("e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAddress), config.StaticTokenProxy(tokenProxyAddress))

	testOrderBytes := getTestOrderBytes()
	newOrder, err := types.OrderFromBytes(testOrderBytes)
	if err != nil {
		t.Errorf("Error parsing order: %v", err.Error())
	}
	filledInt := new(big.Int)
	// Half the taker token amount for the order
	filledInt.SetString("500000000000000000", 10)
	copy(newOrder.TakerAssetAmountFilled[:], gethCommon.LeftPadBytes(filledInt.Bytes(), 32))

	validator.ValidateOrder(newOrder)
}

func TestErrorNoContract(t *testing.T) {
	err := errors.New("no contract code at given address")
	balanceChecker := funds.NewErrorMockBalanceChecker(err)
	feeTokenAddress, _ := hexToAddress("e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAddress), config.StaticTokenProxy(tokenProxyAddress))

	testOrderBytes := getTestOrderBytes()
	newOrder, err := types.OrderFromBytes(testOrderBytes)
	if err != nil {
		t.Errorf("Error parsing order: %v", err.Error())
	}
	filledInt := new(big.Int)
	// Half the taker token amount for the order
	filledInt.SetString("500000000000000000", 10)
	copy(newOrder.TakerAssetAmountFilled[:], gethCommon.LeftPadBytes(filledInt.Bytes(), 32))

	validator.ValidateOrder(newOrder)
}
