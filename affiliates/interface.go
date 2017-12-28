package affiliates

import (
	"github.com/notegio/openrelay/types"
	"math/big"
)

type Affiliate interface {
	Fee() *big.Int
}

type AffiliateService interface {
	Get(*types.Address) (Affiliate, error)
	Set(*types.Address, Affiliate) error
}
