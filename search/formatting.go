package search

import (
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/types"
	"math/big"
	"fmt"
)

type FormattedOrder struct {
	Order *types.Order       `json:"order"`
	Metadata *OrderMetadata  `json:"metaData"`
}

type OrderMetadata struct {
	Hash string                   `json:"hash"`
	FeeRate float64               `json:"feeRate"`
	Status int64                  `json:"status"`
	TakerAssetAmountRemaining string `json:"takerAssetAmountRemaining"`
}

func GetFormattedOrder(order *dbModule.Order) (*FormattedOrder) {
	return &FormattedOrder{
		&order.Order,
		&OrderMetadata{
			fmt.Sprintf("%#x", order.OrderHash[:]),
			order.FeeRate,
			order.Status,
			new(big.Int).Sub(order.TakerAssetAmount.Big(), order.TakerAssetAmountFilled.Big()).String(),
		},
	}
}


type PagedResult struct {
	Total int           `json:"total"`
	Page int            `json:"page"`
	PerPage int         `json:"perPage"`
	Records interface{} `json:"records"`
}


func GetPagedResult(total, page, per_page int, records interface{}) (*PagedResult) {
	return &PagedResult{total, page, per_page, records}
}


type PagedOrders struct {
	Total int           `json:"total"`
	Page int            `json:"page"`
	PerPage int         `json:"perPage"`
	Records []FormattedOrder `json:"records"`
}
