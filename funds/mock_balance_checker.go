package funds

import (
	"errors"
	"math/big"
	"encoding/hex"
)

type mockBalanceChecker struct {
	balances map[[20]byte]map[[20]byte]*big.Int
}

func (funds *mockBalanceChecker) GetBalance(tokenAddrBytes, userAddrBytes [20]byte) (*big.Int, error) {
	if tokenMap, ok := funds.balances[tokenAddrBytes]; ok {
		if balance, ok := tokenMap[userAddrBytes]; ok {
			return balance, nil
		}
		return nil, errors.New("(GetBalance) User address not found " + hex.EncodeToString(userAddrBytes[:]))
	}
	return nil, errors.New("(GetBalance) Token not found " + hex.EncodeToString(tokenAddrBytes[:]))
}

func (funds *mockBalanceChecker) GetAllowance(tokenAddrBytes, userAddrBytes, senderAddress [20]byte) (*big.Int, error) {
	// For now I'm just making GetAllowance match GetBalance for the mock version.
	// Eventually we'll need to test differences between allowance and balance, but
	// for now this will do.
	if tokenMap, ok := funds.balances[tokenAddrBytes]; ok {
		if balance, ok := tokenMap[userAddrBytes]; ok {
			return balance, nil
		}
		return nil, errors.New("(GetAllowance) User address not found " + hex.EncodeToString(userAddrBytes[:]))
	}
	return nil, errors.New("(GetAllowance) Token not found " + hex.EncodeToString(tokenAddrBytes[:]))
}

func NewMockBalanceChecker(balanceMap map[[20]byte]map[[20]byte]*big.Int) BalanceChecker {
	return &mockBalanceChecker{balanceMap}
}
