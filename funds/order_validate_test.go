package funds_test

import (
	"encoding/hex"
	"github.com/notegio/0xrelay/common"
	"github.com/notegio/0xrelay/funds"
	"github.com/notegio/0xrelay/types"
	"math/big"
	"reflect"
	"testing"
)

func createMockBalanceChecker(tokenAddress, userAddress string, tokenAmount int64, feeTokenAmount int64, t *testing.T) funds.BalanceChecker {
	tokenBytes, err := common.HexToBytes(tokenAddress)
	if err != nil {
		t.Errorf(err.Error())
	}
	userBytes, err := common.HexToBytes(userAddress)
	if err != nil {
		t.Errorf(err.Error())
	}
	feeTokenBytes, _ := common.HexToBytes("e41d2489571d322189246dafa5ebde1f4699f498")
	tokenInt := big.NewInt(tokenAmount)
	feeTokenInt := big.NewInt(feeTokenAmount)
	balanceMap := make(map[[20]byte]map[[20]byte]*big.Int)
	balanceMap[tokenBytes] = make(map[[20]byte]*big.Int)
	balanceMap[feeTokenBytes] = make(map[[20]byte]*big.Int)
	balanceMap[tokenBytes][userBytes] = tokenInt
	if reflect.DeepEqual(feeTokenBytes, tokenBytes) {
		balanceMap[feeTokenBytes] = balanceMap[tokenBytes]
	}
	balanceMap[feeTokenBytes][userBytes] = feeTokenInt
	return funds.NewMockBalanceChecker(balanceMap)
}

func TestOrderValidate(t *testing.T) {
	balanceChecker := createMockBalanceChecker("1dad4783cf3fe3085c1426157ab175a6119a04ba", "324454186bb728a3ea55750e0618ff1b18ce6cf8", 0, 0, t)
	validator := funds.NewOrderValidator(balanceChecker)
	testOrderBytes, _ := hex.DecodeString("324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c000000000000000000000000000000000000000090fe2af704b34e0224bf2299c838e04d4dcf1364000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000059938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b00021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	var testOrderByteArray [377]byte
	copy(testOrderByteArray[:], testOrderBytes[:])
	newOrder := types.OrderFromBytes(testOrderByteArray)
	if validator.ValidateOrder(newOrder) {
		t.Errorf("Expected insufficient funds")
	}
}
