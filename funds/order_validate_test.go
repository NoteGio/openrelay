package funds_test

import (
	"encoding/json"
	"encoding/hex"
	"errors"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/config"
	"github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/funds/balance"
	"github.com/notegio/openrelay/types"
	"math/big"
	"reflect"
	"testing"
	"strings"
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

func hexToAssetData(assetDataHex string) (types.AssetData, error) {
	assetDataBytes, err := hex.DecodeString(strings.TrimPrefix(assetDataHex, "0x"))
	if err != nil {
		return nil, err
	}
	assetData := make(types.AssetData, len(assetDataBytes))
	copy(assetData[:], assetDataBytes[:])
	return assetData, nil
}

func createMockBalanceChecker(tokenAddressHex, userAddressHex string, tokenAmount string, feeTokenAmount string, t *testing.T) balance.BalanceChecker {
	tokenAddress, err := hexToAssetData(tokenAddressHex)
	if err != nil {
		t.Errorf(err.Error())
	}
	userAddress, err := hexToAddress(userAddressHex)
	if err != nil {
		t.Errorf(err.Error())
	}
	feeTokenAsset, _ := hexToAssetData("f47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498")
	tokenInt := new(big.Int)
	tokenInt.SetString(tokenAmount, 10)
	feeTokenInt := new(big.Int)
	feeTokenInt.SetString(feeTokenAmount, 10)
	balanceMap := make(map[string]map[types.Address]*big.Int)
	balanceMap[string(tokenAddress[:])] = make(map[types.Address]*big.Int)
	balanceMap[string(feeTokenAsset[:])] = make(map[types.Address]*big.Int)
	balanceMap[string(tokenAddress[:])][*userAddress] = tokenInt
	if reflect.DeepEqual(feeTokenAsset, tokenAddress) {
		balanceMap[string(feeTokenAsset[:])] = balanceMap[string(tokenAddress[:])]
	}
	balanceMap[string(feeTokenAsset[:])][*userAddress] = feeTokenInt
	return balance.NewMockBalanceChecker(balanceMap)
}

func TestOrderValidate(t *testing.T) {
	balanceChecker := createMockBalanceChecker("f47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba", "627306090abab3a6e1400e9345bc60c78a8bef57", "0", "0", t)
	feeTokenAsset, _ := hexToAssetData("f47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAsset), config.StaticTokenProxy(tokenProxyAddress))
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
	balanceChecker := createMockBalanceChecker("f47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba", "627306090abab3a6e1400e9345bc60c78a8bef57", "50000000000000000000", "0", t)
	feeTokenAsset, _ := hexToAssetData("f47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAsset), config.StaticTokenProxy(tokenProxyAddress))
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
	balanceChecker := createMockBalanceChecker("f47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba", "627306090abab3a6e1400e9345bc60c78a8bef57", "25000000000000000000", "0", t)
	feeTokenAsset, _ := hexToAssetData("f47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAsset), config.StaticTokenProxy(tokenProxyAddress))
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
	balanceChecker := balance.NewErrorMockBalanceChecker(err)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic from JSON error")
		}
	}()
	feeTokenAsset, _ := hexToAssetData("f47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAsset), config.StaticTokenProxy(tokenProxyAddress))

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
	balanceChecker := balance.NewErrorMockBalanceChecker(err)
	feeTokenAsset, _ := hexToAssetData("f47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAsset), config.StaticTokenProxy(tokenProxyAddress))

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
