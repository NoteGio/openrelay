package db_test

import (
	dbModule "github.com/notegio/openrelay/db"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"testing"
)


func TestIndexOrder(t *testing.T) {
	db, err := getDb()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	tx := db.Begin()
	defer func(){
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	indexer := dbModule.NewIndexer(tx)
	order := sampleOrder()
	if err := indexer.Index(order); err != nil {
		t.Errorf(err.Error())
	}
}
