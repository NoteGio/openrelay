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
	order := sampleOrder(t)
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *order
	if err := dbOrder.Save(tx, dbModule.StatusOpen, nil).Error; err != nil {
		t.Errorf(err.Error())
	}
	fillString := fmt.Sprintf(
		"{\"tokenAddress\": \"%v\",\"spenderAddress\": \"%v\",\"zrxToken\": \"%v\",\"balance\": \"%v\"}",
		order.MakerAssetData.Address(),
		order.Maker,
		order.TakerAssetData.Address(),
		order.MakerAssetAmount,
	)
	publisher, channel := channels.MockChannel()
	consumer := dbModule.NewRecordSpendConsumer(tx, 1, nil)
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

func TestSpendConsumerMakerZRX(t *testing.T) {
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
	order := sampleOrder(t)
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *order
	if err := dbOrder.Save(tx, dbModule.StatusOpen, nil).Error; err != nil {
		t.Errorf(err.Error())
	}
	fillString := fmt.Sprintf(
		"{\"tokenAddress\": \"%v\",\"spenderAddress\": \"%v\",\"zrxToken\": \"%v\",\"balance\": \"%v\"}",
		order.TakerAssetData.Address(),
		order.Maker,
		order.TakerAssetData.Address(),
		order.MakerAssetAmount,
	)
	tx.LogMode(true)
	defer tx.LogMode(false)
	publisher, channel := channels.MockChannel()
	consumer := dbModule.NewRecordSpendConsumer(tx, 1, nil)
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
	order := sampleOrder(t)
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *order
	if err := dbOrder.Save(tx, dbModule.StatusOpen, nil).Error; err != nil {
		t.Errorf(err.Error())
	}
	fillString := fmt.Sprintf(
		"{\"tokenAddress\": \"%v\",\"spenderAddress\": \"%v\",\"zrxToken\": \"%v\",\"balance\": \"0\"}",
		order.MakerAssetData.Address(),
		order.Maker,
		order.TakerAssetData.Address(),
	)
	publisher, channel := channels.MockChannel()
	dsPublisher, ch := channels.MockPublisher()
	consumer := dbModule.NewRecordSpendConsumer(tx, 1, dsPublisher)
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
	select {
	case <-ch:
	default:
		t.Errorf("Expected spent item to be published")
	}
}
