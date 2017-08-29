package affiliates

import (
	"math/big"
)

type Affiliate interface {
	Fee() *big.Int
}

type AffiliateService interface {
	Get([20]byte) (Affiliate, error)
	Set([20]byte, Affiliate) error
}
