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
	StatusOpen      = int64(0)
	StatusFilled    = int64(1)
	StatusUnfunded  = int64(2)
	StatusCancelled = int64(3)
)

type Order struct {
	types.Order
	CreatedAt time.Time
	UpdatedAt time.Time
	OrderHash []byte  `gorm:"primary_key"`
	Status    int64   `gorm:"index"`
	Price     float64 `gorm:"index:price"`
	FeeRate   float64 `gorm:"index:price"`
	MakerAssetRemaining *types.Uint256
	MakerFeeRemaining *types.Uint256
}

func (order *Order) TableName() string {
	return "orderv2"
}

// Save records the order in the database, defaulting to the specified status.
// Status should either be db.StatusOpen, or db.StatusUnfunded. If the order
// is filled based on order.TakerAssetAmountFilled + order.TakerAssetAmountCancelled
// the status will be recorded as db.StatusFilled regardless of the specified status.
func (order *Order) Save(db *gorm.DB, status int64) *gorm.DB {
	if !order.Signature.Verify(order.Maker, order.Hash()) {
		scope := db.New()
		scope.AddError(errors.New("Failed to verify signature"))
		return scope
	}
	order.OrderHash = order.Hash()
	log.Printf("Attempting to save order %#x", order.Hash())
	if order.Cancelled {
		order.Status = StatusCancelled
	}

	remainingAmount := order.TakerAssetAmount.Big()
	thresholdAmount := order.TakerAssetAmount.Big()
	// If the order is larger than 100 base units, we want to remove it from the
	// orderbook when it's 99% filled. For fewer than 100 base units, integer
	// arithmetic makes this not work the way we want it to.
	if thresholdAmount.Cmp(big.NewInt(100)) > 0 {
		thresholdAmount.Mul(thresholdAmount, big.NewInt(99))
		thresholdAmount.Div(thresholdAmount, big.NewInt(100))
	}
	remainingAmount.Sub(remainingAmount, order.TakerAssetAmountFilled.Big())
	thresholdAmount.Sub(thresholdAmount, order.TakerAssetAmountFilled.Big())

	// makerRemainingInt = (MakerAssetAmount * RemainingAmount) / TakerAssetAmount
	makerRemainingInt := new(big.Int).Div(new(big.Int).Mul(order.MakerAssetAmount.Big(), remainingAmount), order.TakerAssetAmount.Big())
	makerRemaining := &types.Uint256{}
	copy(makerRemaining[:], abi.U256(makerRemainingInt))
	makerFeeRemainingInt := new(big.Int).Div(new(big.Int).Mul(order.MakerFee.Big(), remainingAmount), order.TakerAssetAmount.Big())
	makerFeeRemaining := &types.Uint256{}
	copy(makerFeeRemaining[:], abi.U256(makerFeeRemainingInt))
	updates := map[string]interface{}{
		"taker_asset_amount_filled":    order.TakerAssetAmountFilled,
		"maker_asset_remaining":        makerRemaining,
		"maker_fee_remaining":          makerFeeRemaining,
		"status":                       status,
	}
	if thresholdAmount.Cmp(big.NewInt(0)) <= 0 {
		updates["status"] = StatusFilled
	}
	updateScope := db.Model(Order{}).Where("order_hash = ?", order.OrderHash).Updates(updates)
	if updateScope.Error != nil {
		log.Printf(updateScope.Error.Error())
	}
	if updateScope.RowsAffected > 0 {
		return updateScope
	}

	takerTokenAmount := new(big.Float).SetInt(order.TakerAssetAmount.Big())
	makerTokenAmount := new(big.Float).SetInt(order.MakerAssetAmount.Big())
	takerFeeAmount := new(big.Float).SetInt(order.TakerFee.Big())

	order.Price, _ = new(big.Float).Quo(takerTokenAmount, makerTokenAmount).Float64()
	order.FeeRate, _ = new(big.Float).Quo(takerFeeAmount, takerTokenAmount).Float64()
	order.Status = status
	order.MakerAssetRemaining = makerRemaining
	order.MakerFeeRemaining = makerFeeRemaining
	return db.Create(order)
}
