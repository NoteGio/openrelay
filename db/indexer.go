package db

import (
	"github.com/notegio/openrelay/types"
	"github.com/jinzhu/gorm"
)

type Indexer struct {
	db *gorm.DB
	status int64
}

func (indexer *Indexer) Index(order *types.Order) (error) {
	dbOrder := Order{}
	dbOrder.Order = *order
	return dbOrder.Save(indexer.db, indexer.status).Error
}

func NewIndexer(db *gorm.DB, status int64) (*Indexer) {
	return &Indexer{db, status}
}
