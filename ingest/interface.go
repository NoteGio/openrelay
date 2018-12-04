package ingest

import (
	"github.com/notegio/openrelay/types"
)

type TermsManager interface {
	CheckAddress(*types.Address) (<-chan bool)
}

type ExchangeLookup interface {
	ExchangeIsKnown(*types.Address) (<-chan bool)
}
