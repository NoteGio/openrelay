package db

import (
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
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
	MakerTokenRemaining *types.Uint256
	MakerFeeRemaining *types.Uint256
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
	thresholdAmount := new(big.Int)
	remainingAmount.SetBytes(order.TakerTokenAmount[:])
	thresholdAmount.SetBytes(order.TakerTokenAmount[:])
	// We want to consider it filled once it's 99% filled
	thresholdAmount.Mul(thresholdAmount, new(big.Int).SetInt64(99))
	thresholdAmount.Div(thresholdAmount, new(big.Int).SetInt64(100))
	remainingAmount.Sub(remainingAmount, new(big.Int).SetBytes(order.TakerTokenAmountFilled[:]))
	remainingAmount.Sub(remainingAmount, new(big.Int).SetBytes(order.TakerTokenAmountCancelled[:]))
	thresholdAmount.Sub(thresholdAmount, new(big.Int).SetBytes(order.TakerTokenAmountFilled[:]))
	thresholdAmount.Sub(thresholdAmount, new(big.Int).SetBytes(order.TakerTokenAmountCancelled[:]))

	// makerRemainingInt = (MakerTokenAmount * RemainingAmount) / TakerTokenAmount
	makerRemainingInt := new(big.Int).Div(new(big.Int).Mul(new(big.Int).SetBytes(order.MakerTokenAmount[:]), remainingAmount), new(big.Int).SetBytes(order.TakerTokenAmount[:]))
	makerRemaining := &types.Uint256{}
	copy(makerRemaining[:], abi.U256(makerRemainingInt))
	makerFeeRemainingInt := new(big.Int).Div(new(big.Int).Mul(new(big.Int).SetBytes(order.MakerFee[:]), remainingAmount), new(big.Int).SetBytes(order.TakerTokenAmount[:]))
	makerFeeRemaining := &types.Uint256{}
	copy(makerFeeRemaining[:], abi.U256(makerFeeRemainingInt))
	updates := map[string]interface{}{
		"taker_token_amount_filled":    order.TakerTokenAmountFilled,
		"taker_token_amount_cancelled": order.TakerTokenAmountCancelled,
		"maker_token_remaining":        makerRemaining,
		"maker_fee_remaining":          makerFeeRemaining,
		"status":                       status,
	}
	if thresholdAmount.Cmp(new(big.Int).SetInt64(0)) <= 0 {
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
	order.MakerTokenRemaining = makerRemaining
	order.MakerFeeRemaining = makerFeeRemaining
	return db.Create(order)
}
