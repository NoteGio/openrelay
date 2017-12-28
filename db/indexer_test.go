package db_test

import (
	"encoding/json"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	dbModule "github.com/notegio/openrelay/db"
	"math/big"
	"reflect"
	"testing"
)

func TestIndexOrder(t *testing.T) {
	db, err := getDb()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	indexer := dbModule.NewIndexer(tx, dbModule.StatusOpen)
	order := sampleOrder()
	if err := indexer.Index(order); err != nil {
		t.Errorf(err.Error())
	}
}

func TestFillIndex(t *testing.T) {
	db, err := getDb()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	indexer := dbModule.NewIndexer(tx, dbModule.StatusOpen)
	order := sampleOrder()
	if err := indexer.Index(order); err != nil {
		t.Errorf(err.Error())
	}
	takerTokenAmount := new(big.Int).SetBytes(order.TakerTokenAmount[:])
	fillString := fmt.Sprintf(
		"{\"orderHash\": \"%#x\", \"filledTakerTokenAmount\": \"%v\"}",
		order.Hash(),
		takerTokenAmount.String(),
	)
	fillRecord := &dbModule.FillRecord{}
	if err := json.Unmarshal([]byte(fillString), fillRecord); err != nil {
		t.Errorf(err.Error())
	}
	if err := indexer.RecordFill(fillRecord); err != nil {
		t.Errorf(err.Error())
	}
	dbOrder := &dbModule.Order{}
	dbOrder.Initialize()
	tx.Model(&dbModule.Order{}).Where("order_hash = ?", order.Hash()).First(dbOrder)
	if !reflect.DeepEqual(dbOrder.TakerTokenAmount, dbOrder.TakerTokenAmountFilled) {
		t.Errorf("TakerTokenAmount should match TakerTokenAmountFilled, got %#x != %#x", dbOrder.TakerTokenAmount[:], dbOrder.TakerTokenAmountFilled[:])
	}
	if dbOrder.Status != dbModule.StatusFilled {
		t.Errorf("Order status should be filled, got %v", dbOrder.Status)
	}
}
