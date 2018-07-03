package db

import (
	"github.com/notegio/openrelay/types"
	"github.com/jinzhu/gorm"
)

type Cancellation struct {
	Maker  *types.Address `gorm:"primary_key"`
	Sender *types.Address `gorm:"primary_key"`
	Epoch  *types.Uint256
}

func (cancellation *Cancellation) Save(db *gorm.DB) *gorm.DB {
	return db.Where(
		&Cancellation{Maker: cancellation.Maker, Sender: cancellation.Sender},
	).Assign(Cancellation{Epoch: cancellation.Epoch}).FirstOrCreate(cancellation)
}

func GetCancellationEpoch(maker, sender *types.Address, db *gorm.DB) (*types.Uint256) {
	cancellation := &Cancellation{}
	if err := db.Where("maker = ? AND sender = ?", maker, sender).First(cancellation).Error; err != nil {
		return &types.Uint256{}
	}
	return cancellation.Epoch
}
