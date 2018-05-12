package funds_test

import (
	"encoding/hex"
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
	balanceChecker := createMockBalanceChecker("1dad4783cf3fe3085c1426157ab175a6119a04ba", "324454186bb728a3ea55750e0618ff1b18ce6cf8", "0", "0", t)
	feeTokenAddress, _ := hexToAddress("e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAddress), config.StaticTokenProxy(tokenProxyAddress))
	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	var testOrderByteArray [441]byte
	copy(testOrderByteArray[:], testOrderBytes[:])
	newOrder := types.OrderFromBytes(testOrderByteArray)
	if result, _ := validator.ValidateOrder(newOrder); result {
		t.Errorf("Expected insufficient funds")
	}
}

func TestOrderValidateSufficient(t *testing.T) {
	balanceChecker := createMockBalanceChecker("1dad4783cf3fe3085c1426157ab175a6119a04ba", "324454186bb728a3ea55750e0618ff1b18ce6cf8", "50000000000000000000", "0", t)
	feeTokenAddress, _ := hexToAddress("e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAddress), config.StaticTokenProxy(tokenProxyAddress))
	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	var testOrderByteArray [441]byte
	copy(testOrderByteArray[:], testOrderBytes[:])
	newOrder := types.OrderFromBytes(testOrderByteArray)

	if result, _ := validator.ValidateOrder(newOrder); !result {
		t.Errorf("Expected sufficient funds")
	}
}

func TestOrderValidateSpent(t *testing.T) {
	// 25000000000000000000 is half of the makerTokenAmount for the sample order
	balanceChecker := createMockBalanceChecker("1dad4783cf3fe3085c1426157ab175a6119a04ba", "324454186bb728a3ea55750e0618ff1b18ce6cf8", "25000000000000000000", "0", t)
	feeTokenAddress, _ := hexToAddress("e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAddress), config.StaticTokenProxy(tokenProxyAddress))
	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	var testOrderByteArray [441]byte
	copy(testOrderByteArray[:], testOrderBytes[:])
	newOrder := types.OrderFromBytes(testOrderByteArray)
	filledInt := new(big.Int)
	// Half the taker token amount for the order
	filledInt.SetString("500000000000000000", 10)
	copy(newOrder.TakerTokenAmountFilled[:], gethCommon.LeftPadBytes(filledInt.Bytes(), 32))

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

	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	var testOrderByteArray [441]byte
	copy(testOrderByteArray[:], testOrderBytes[:])
	newOrder := types.OrderFromBytes(testOrderByteArray)
	filledInt := new(big.Int)
	// Half the taker token amount for the order
	filledInt.SetString("500000000000000000", 10)
	copy(newOrder.TakerTokenAmountFilled[:], gethCommon.LeftPadBytes(filledInt.Bytes(), 32))

	validator.ValidateOrder(newOrder)
}

func TestErrorNoContract(t *testing.T) {
	err := errors.New("no contract code at given address")
	balanceChecker := funds.NewErrorMockBalanceChecker(err)
	feeTokenAddress, _ := hexToAddress("e41d2489571d322189246dafa5ebde1f4699f498")
	tokenProxyAddress, _ := hexToAddress("d4fd252d7d2c9479a8d616f510eac6243b5dddf9")
	validator := funds.NewOrderValidator(balanceChecker, config.StaticFeeToken(feeTokenAddress), config.StaticTokenProxy(tokenProxyAddress))

	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	var testOrderByteArray [441]byte
	copy(testOrderByteArray[:], testOrderBytes[:])
	newOrder := types.OrderFromBytes(testOrderByteArray)
	filledInt := new(big.Int)
	// Half the taker token amount for the order
	filledInt.SetString("500000000000000000", 10)
	copy(newOrder.TakerTokenAmountFilled[:], gethCommon.LeftPadBytes(filledInt.Bytes(), 32))

	validator.ValidateOrder(newOrder)
}
