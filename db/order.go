package db

import (
	"github.com/notegio/openrelay/types"
	"github.com/jinzhu/gorm"
	"math/big"
	"time"
	"log"
)

type Order struct {
	types.Order
	CreatedAt time.Time
	UpdatedAt time.Time
	OrderHash     []byte `gorm:"primary_key"`
	Filled bool `gorm:"index"`
}

func (order *Order) Save(db *gorm.DB) (*gorm.DB) {
	order.OrderHash = order.Hash()
	remainingAmount := new(big.Int)
	remainingAmount.SetBytes(order.TakerTokenAmount[:])
	remainingAmount.Sub(remainingAmount, new(big.Int).SetBytes(order.TakerTokenAmountFilled[:]))
	remainingAmount.Sub(remainingAmount, new(big.Int).SetBytes(order.TakerTokenAmountCancelled[:]))
	updateScope := db.Model(Order{}).Where("order_hash = ?", order.OrderHash).Updates(map[string]interface{}{
		"taker_token_amount_filled": order.TakerTokenAmountFilled,
		"taker_token_amount_cancelled": order.TakerTokenAmountCancelled,
		"update_at": time.Now(),
		"filled": remainingAmount.Cmp(new(big.Int).SetInt64(0)) <= 0,
	})
	if updateScope.Error != nil {
		log.Printf(updateScope.Error.Error())
	}
	if updateScope.RowsAffected > 0 {
		return updateScope
	}
	order.CreatedAt = time.Now()
	order.UpdatedAt = order.CreatedAt
	return db.Create(order)
}
