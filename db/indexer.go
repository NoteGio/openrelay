package db

import (
	"github.com/notegio/openrelay/types"
	"github.com/jinzhu/gorm"
)

type Indexer struct {
	db *gorm.DB
}

func (indexer *Indexer) Index(order *types.Order) (error) {
	dbOrder := Order{}
	dbOrder.Order = *order
	return dbOrder.Save(indexer.db).Error
}

func NewIndexer(db *gorm.DB) (*Indexer) {
	return &Indexer{db}
}
