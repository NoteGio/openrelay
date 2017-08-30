package affiliates

import (
	"math/big"
)

type affiliate struct {
	baseFee    *big.Int
	feePercent int64
}

func (acct *affiliate) Fee() *big.Int {
	fee := new(big.Int)
	fee.Mul(acct.baseFee, big.NewInt(acct.feePercent))
	return fee.Div(fee, big.NewInt(100))
}

func NewAffiliate(baseFee *big.Int, feePercent int64) Affiliate {
	return &affiliate{baseFee, feePercent}
}
