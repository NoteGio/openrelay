package affiliates

import (
	"math/big"
)

type affiliate struct {
	BaseFee    *big.Int
	FeePercent int64
}

func (acct *affiliate) Fee() *big.Int {
	fee := new(big.Int)
	fee.Mul(acct.BaseFee, big.NewInt(acct.FeePercent))
	return fee.Div(fee, big.NewInt(100))
}

func NewAffiliate(baseFee *big.Int, feePercent int64) Affiliate {
	return &affiliate{baseFee, feePercent}
}
