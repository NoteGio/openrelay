package accounts

import (
	"github.com/notegio/openrelay/types"
	"math/big"
)

type Account interface {
	Blacklisted() bool
	Discount() *big.Int
}

type AccountService interface {
	Get(*types.Address) Account
	Set(*types.Address, Account) error
}
