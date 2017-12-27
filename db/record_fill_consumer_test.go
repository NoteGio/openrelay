package db_test

import (
	"github.com/notegio/openrelay/channels"
	dbModule "github.com/notegio/openrelay/db"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
	"testing"
	"reflect"
	"math/big"
	"fmt"
)

func TestFillConsumer(t *testing.T) {
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
	publisher, channel := channels.MockChannel()
	consumer := dbModule.NewRecordFillConsumer(tx)
	channel.AddConsumer(consumer)
	channel.StartConsuming()
	defer channel.StopConsuming()
	publisher.Publish(fillString)
	time.Sleep(100 * time.Millisecond)

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
