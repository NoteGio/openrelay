package db_test

import (
	"fmt"
	"github.com/notegio/openrelay/channels"
	dbModule "github.com/notegio/openrelay/db"
	"testing"
)

func TestSpendConsumer(t *testing.T) {
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
	order := sampleOrder()
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *order
	if err := dbOrder.Save(tx, dbModule.StatusOpen).Error; err != nil {
		t.Errorf(err.Error())
	}
	fillString := fmt.Sprintf(
		"{\"tokenAddress\": \"%v\",\"spenderAddress\": \"%v\",\"zrxToken\": \"%v\",\"balance\": \"%v\"}",
		order.MakerToken,
		order.Maker,
		order.TakerToken,
		order.MakerTokenAmount,
	)
	publisher, channel := channels.MockChannel()
	consumer := dbModule.NewRecordSpendConsumer(tx)
	channel.AddConsumer(consumer)
	channel.StartConsuming()
	defer channel.StopConsuming()
	publisher.Publish(fillString)
	if err := channels.MockFinish(channel, 1); err != nil {
		t.Errorf(err.Error())
	}

	dbOrder = &dbModule.Order{}
	dbOrder.Initialize()
	if err := tx.Model(&dbModule.Order{}).Where("order_hash = ?", order.Hash()).First(dbOrder).Error; err != nil {
		t.Errorf(err.Error())
	}
	if dbOrder.Status != dbModule.StatusOpen {
		t.Errorf("Order status should be open, got %v", dbOrder.Status)
	}
}
func TestSpendConsumerInsufficient(t *testing.T) {
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
	order := sampleOrder()
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *order
	if err := dbOrder.Save(tx, dbModule.StatusOpen).Error; err != nil {
		t.Errorf(err.Error())
	}
	fillString := fmt.Sprintf(
		"{\"tokenAddress\": \"%v\",\"spenderAddress\": \"%v\",\"zrxToken\": \"%v\",\"balance\": \"0\"}",
		order.MakerToken,
		order.Maker,
		order.TakerToken,
	)
	publisher, channel := channels.MockChannel()
	consumer := dbModule.NewRecordSpendConsumer(tx)
	channel.AddConsumer(consumer)
	channel.StartConsuming()
	defer channel.StopConsuming()
	publisher.Publish(fillString)
	if err := channels.MockFinish(channel, 1); err != nil {
		t.Errorf(err.Error())
	}

	dbOrder = &dbModule.Order{}
	dbOrder.Initialize()
	if err := tx.Model(&dbModule.Order{}).Where("order_hash = ?", order.Hash()).First(dbOrder).Error; err != nil {
		t.Errorf(err.Error())
	}
	if dbOrder.Status != dbModule.StatusUnfunded {
		t.Errorf("Order status should be unfunded, got %v", dbOrder.Status)
	}
}
