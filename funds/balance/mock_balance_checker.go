package balance

import (
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/channels"
	"math/big"
	"fmt"
)

type mockBalanceChecker struct {
	balances map[string]map[types.Address]*big.Int
}

func (funds *mockBalanceChecker) GetBalance(assetData types.AssetData, userAddrBytes *types.Address) (*big.Int, error) {
	if tokenMap, ok := funds.balances[string(assetData[:])]; ok {
		if balance, ok := tokenMap[*userAddrBytes]; ok {
			return balance, nil
		}
		return nil, fmt.Errorf("(GetBalance) User address not found: '%#x'", userAddrBytes[:])
	}
	return nil, fmt.Errorf("(GetBalance) Token not found: '%#x'", assetData[:])
}

func (funds *mockBalanceChecker) GetAllowance(assetData types.AssetData, userAddrBytes, senderAddress *types.Address) (*big.Int, error) {
	// For now I'm just making GetAllowance match GetBalance for the mock version.
	// Eventually we'll need to test differences between allowance and balance, but
	// for now this will do.
	if tokenMap, ok := funds.balances[string(assetData[:])]; ok {
		if balance, ok := tokenMap[*userAddrBytes]; ok {
			return balance, nil
		}
		return nil, fmt.Errorf("(GetBalance) User address not found: '%#x'", userAddrBytes[:])
	}
	return nil, fmt.Errorf("(GetBalance) Token not found: '%#x'", assetData[:])
}

func (funds *mockBalanceChecker) Consume(msg channels.Delivery) {
	msg.Ack()
}

func NewMockBalanceChecker(balanceMap map[string]map[types.Address]*big.Int) BalanceChecker {
	return &mockBalanceChecker{balanceMap}
}

type errorMockBalanceChecker struct {
	err error
}

func (funds *errorMockBalanceChecker) GetBalance(assetData types.AssetData, userAddrBytes *types.Address) (*big.Int, error) {
	return nil, funds.err
}

func (funds *errorMockBalanceChecker) GetAllowance(assetData types.AssetData, userAddrBytes, senderAddress *types.Address) (*big.Int, error) {
	return nil, funds.err
}

func (funds *errorMockBalanceChecker) Consume(msg channels.Delivery) {
	msg.Ack()
}

func NewErrorMockBalanceChecker(err error) BalanceChecker {
	return &errorMockBalanceChecker{err}
}
