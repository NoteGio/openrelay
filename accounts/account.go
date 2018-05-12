package accounts

import (
	"math/big"
	"time"
)

type account struct {
	Bl        bool
	BaseFee            *big.Int `json:"-"`
	Dp int64
	Expiration         int64
}

func (acct *account) Blacklisted() bool {
	return acct.Bl
}

func (acct *account) Discount() *big.Int {
	if acct.Expiration < time.Now().Unix() {
		// Account is expired. No discount
		return new(big.Int)
	}
	discount := new(big.Int)
	discount.Mul(acct.BaseFee, big.NewInt(acct.Dp))
	return discount.Div(discount, big.NewInt(100))
}

func NewAccount(Bl bool, BaseFee *big.Int, Dp, Expiration int64) Account {
	return &account{Bl, BaseFee, Dp, Expiration}
}
