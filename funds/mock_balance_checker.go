package funds

import (
	"encoding/hex"
	"errors"
	"math/big"
	"github.com/notegio/openrelay/types"
)

type mockBalanceChecker struct {
	balances map[types.Address]map[types.Address]*big.Int
}

func (funds *mockBalanceChecker) GetBalance(tokenAddrBytes, userAddrBytes *types.Address) (*big.Int, error) {
	if tokenMap, ok := funds.balances[*tokenAddrBytes]; ok {
		if balance, ok := tokenMap[*userAddrBytes]; ok {
			return balance, nil
		}
		return nil, errors.New("(GetBalance) User address not found " + hex.EncodeToString(userAddrBytes[:]))
	}
	return nil, errors.New("(GetBalance) Token not found " + hex.EncodeToString(tokenAddrBytes[:]))
}

func (funds *mockBalanceChecker) GetAllowance(tokenAddrBytes, userAddrBytes, senderAddress *types.Address) (*big.Int, error) {
	// For now I'm just making GetAllowance match GetBalance for the mock version.
	// Eventually we'll need to test differences between allowance and balance, but
	// for now this will do.
	if tokenMap, ok := funds.balances[*tokenAddrBytes]; ok {
		if balance, ok := tokenMap[*userAddrBytes]; ok {
			return balance, nil
		}
		return nil, errors.New("(GetAllowance) User address not found " + hex.EncodeToString(userAddrBytes[:]))
	}
	return nil, errors.New("(GetAllowance) Token not found " + hex.EncodeToString(tokenAddrBytes[:]))
}

func NewMockBalanceChecker(balanceMap map[types.Address]map[types.Address]*big.Int) BalanceChecker {
	return &mockBalanceChecker{balanceMap}
}

type errorMockBalanceChecker struct {
	err error
}

func (funds *errorMockBalanceChecker) GetBalance(tokenAddrBytes, userAddrBytes *types.Address) (*big.Int, error) {
	return nil, funds.err
}

func (funds *errorMockBalanceChecker) GetAllowance(tokenAddrBytes, userAddrBytes, senderAddress *types.Address) (*big.Int, error) {
	return nil, funds.err
}

func NewErrorMockBalanceChecker(err error) BalanceChecker {
	return &errorMockBalanceChecker{err}
}
