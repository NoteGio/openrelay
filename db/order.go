package db

import (
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/types"
	"log"
	"math/big"
	"time"
	"errors"
)

const (
	StatusOpen     = int64(0)
	StatusFilled   = int64(1)
	StatusUnfunded = int64(2)
)

type Order struct {
	types.Order
	CreatedAt time.Time
	UpdatedAt time.Time
	OrderHash []byte  `gorm:"primary_key"`
	Status    int64   `gorm:"index"`
	Price     float64 `gorm:"index:price"`
	FeeRate   float64 `gorm:"index:price"`
}

// Save records the order in the database, defaulting to the specified status.
// Status should either be db.StatusOpen, or db.StatusUnfunded. If the order
// is filled based on order.TakerTokenAmountFilled + order.TakerTokenAmountCancelled
// the status will be recorded as db.StatusFilled regardless of the specified status.
func (order *Order) Save(db *gorm.DB, status int64) *gorm.DB {
	copy(order.Signature.Hash[:], order.Hash())
	if !order.Signature.Verify(order.Maker) {
		scope := db.New()
		scope.AddError(errors.New("Failed to verify signature"))
		return scope
	}
	order.OrderHash = order.Hash()
	log.Printf("Attempting to save order %#x", order.Hash())

	remainingAmount := new(big.Int)
	remainingAmount.SetBytes(order.TakerTokenAmount[:])
	// We want to consider it filled once it's 99% filled
	remainingAmount.Mul(remainingAmount, new(big.Int).SetInt64(99))
	remainingAmount.Div(remainingAmount, new(big.Int).SetInt64(100))
	remainingAmount.Sub(remainingAmount, new(big.Int).SetBytes(order.TakerTokenAmountFilled[:]))
	remainingAmount.Sub(remainingAmount, new(big.Int).SetBytes(order.TakerTokenAmountCancelled[:]))
	updates := map[string]interface{}{
		"taker_token_amount_filled":    order.TakerTokenAmountFilled,
		"taker_token_amount_cancelled": order.TakerTokenAmountCancelled,
		"status":                       status,
	}
	if remainingAmount.Cmp(new(big.Int).SetInt64(0)) <= 0 {
		updates["status"] = StatusFilled
	}
	updateScope := db.Model(Order{}).Where("order_hash = ?", order.OrderHash).Updates(updates)
	if updateScope.Error != nil {
		log.Printf(updateScope.Error.Error())
	}
	if updateScope.RowsAffected > 0 {
		return updateScope
	}

	takerTokenAmount := new(big.Float).SetInt(new(big.Int).SetBytes(order.TakerTokenAmount[:]))
	makerTokenAmount := new(big.Float).SetInt(new(big.Int).SetBytes(order.MakerTokenAmount[:]))
	takerFeeAmount := new(big.Float).SetInt(new(big.Int).SetBytes(order.TakerFee[:]))

	order.Price, _ = new(big.Float).Quo(takerTokenAmount, makerTokenAmount).Float64()
	order.FeeRate, _ = new(big.Float).Quo(takerFeeAmount, takerTokenAmount).Float64()
	order.Status = status
	return db.Create(order)
}
