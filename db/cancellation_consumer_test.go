package db_test

import (
	"github.com/notegio/openrelay/channels"
	dbModule "github.com/notegio/openrelay/db"
	"reflect"
	"testing"
	"time"
	"fmt"
)

func TestCancellationConsumer(t *testing.T) {
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
	if err := tx.AutoMigrate(&dbModule.Cancellation{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	indexer := dbModule.NewIndexer(tx, dbModule.StatusOpen, nil)
	order := sampleOrder(t)
	if !order.Signature.Verify(order.Maker, order.Hash()) {
		t.Errorf("Failed to verify signature")
	}
	if err := indexer.Index(order); err != nil {
		t.Errorf(err.Error())
	}
	fillString := fmt.Sprintf("{\"Maker\": \"%#x\", \"Sender\": \"0x0000000000000000000000000000000000000000\", \"Epoch\": \"11065671350908846865864045738088581419204014210814002044381812654087807532\"}", order.Maker[:])
	publisher, channel := channels.MockChannel()
	dsPublisher, ch := channels.MockPublisher()
	consumer := dbModule.NewRecordCancellationConsumer(tx, 1, dsPublisher)
	channel.AddConsumer(consumer)
	channel.StartConsuming()
	defer channel.StopConsuming()
	publisher.Publish(fillString)
	time.Sleep(100 * time.Millisecond)

	dbOrder := &dbModule.Order{}
	dbOrder.Initialize()
	if err := tx.Model(&dbModule.Order{}).Where("order_hash = ?", order.Hash()).First(dbOrder).Error; err != nil {
		t.Errorf(err.Error())
	}
	if dbOrder.Status != dbModule.StatusCancelled {
		t.Errorf("Order status should be cancelled, got %v", dbOrder.Status)
	}
	cancellation := &dbModule.Cancellation{}
	if err := tx.Model(&dbModule.Cancellation{}).Where("maker = ? AND sender = ?", dbOrder.Maker, dbOrder.SenderAddress).First(cancellation).Error; err != nil {
		t.Errorf("Error getting cancellation: %v", err.Error())
	}
	if !reflect.DeepEqual(cancellation.Maker, dbOrder.Maker) || !reflect.DeepEqual(cancellation.Sender, dbOrder.SenderAddress) {
		t.Errorf("Cancellation does not match order: %v", cancellation.Maker)
	}
	select {
	case <-ch:
	default:
		t.Errorf("Expected canceled item to be published")
	}
}
