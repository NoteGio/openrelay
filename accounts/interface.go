package accounts

import (
	"math/big"
	"github.com/notegio/openrelay/types"
)

type Account interface {
	Blacklisted() bool
	Discount() *big.Int
}

type AccountService interface {
	Get(*types.Address) Account
	Set(*types.Address, Account) error
}
