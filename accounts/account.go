package accounts

import (
	"time"
	"math/big"
)

type account struct {
	blacklisted bool
	baseFee *big.Int `json:"-"`
	discountPercentage int64
	expiration int64
}

func (acct *account) Blacklisted() bool {
	return acct.blacklisted
}

func (acct *account) Discount() *big.Int {
	if acct.expiration < time.Now().Unix() {
		// Account is expired. No discount
		return new(big.Int)
	}
	discount := new(big.Int)
	discount.Mul(acct.baseFee, big.NewInt(acct.discountPercentage))
	return discount.Div(discount, big.NewInt(100))
}

func NewAccount(blacklisted bool, baseFee *big.Int, discountPercentage, expiration int64) Account {
	return &account{blacklisted, baseFee, discountPercentage, expiration}
}
