package affiliates

import (
	"math/big"
	"github.com/notegio/openrelay/types"
)

type Affiliate interface {
	Fee() *big.Int
}

type AffiliateService interface {
	Get(*types.Address) (Affiliate, error)
	Set(*types.Address, Affiliate) error
}
