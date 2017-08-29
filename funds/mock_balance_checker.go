package funds

import (
	"math/big"
	"errors"
)

type mockBalanceChecker struct {
	balances map[[20]byte]map[[20]byte]*big.Int
}

func (funds *mockBalanceChecker)GetBalance(tokenAddrBytes, userAddrBytes [20]byte) (*big.Int, error) {

	if tokenMap, ok := funds.balances[tokenAddrBytes]; ok {
		if balance, ok := tokenMap[userAddrBytes]; ok {
			return balance, nil
		}
		return nil, errors.New("User address not found")
	}
	return nil, errors.New("Token not found")
}

func NewMockBalanceChecker(balanceMap map[[20]byte]map[[20]byte]*big.Int) (BalanceChecker){
	return &mockBalanceChecker{balanceMap}
}
