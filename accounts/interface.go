package accounts

import "math/big"

type Account interface {
	Blacklisted() bool
	Discount() *big.Int
}

type AccountService interface {
	Get([20]byte) Account
	Set([20]byte, Account) error
}
