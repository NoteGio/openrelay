package db_test

import (
	"github.com/notegio/openrelay/channels"
	dbModule "github.com/notegio/openrelay/db"
	"reflect"
	"testing"
	"time"
)

func IndexConsumerDefaultStatus(status int64, t *testing.T) {
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
	publisher, channel := channels.MockChannel()
	consumer := dbModule.NewIndexConsumer(tx, status, 1)
	channel.AddConsumer(consumer)
	channel.StartConsuming()
	defer channel.StopConsuming()
	orderBytes := order.Bytes()
	publisher.Publish(string(orderBytes[:]))
	time.Sleep(100 * time.Millisecond)
	dbOrder := &dbModule.Order{}
	tx.Model(&dbModule.Order{}).Where("order_hash = ?", order.Hash()).First(dbOrder)
	if !reflect.DeepEqual(dbOrder.Signature[:], order.Signature[:]) {
		t.Errorf("Queried signature doesn't match saved signature: '%#x' != '%#x'", dbOrder.Signature, order.Signature)
	}
	if !reflect.DeepEqual(dbOrder.Bytes(), order.Bytes()) {
		t.Errorf(
			"Queried order not equal to saved order; '%#x' != '%#x'", dbOrder.Bytes(), order.Bytes(),
		)
	}
	if dbOrder.Status != status {
		t.Errorf("Unexpected order status. Wanted %v, got %v", status, dbOrder.Status)
	}
	if channel.PurgeRejected() > 0 {
		t.Errorf("Failed to record order")
	}

}

func TestIndexConsumerOpenStatus(t *testing.T) {
	IndexConsumerDefaultStatus(dbModule.StatusOpen, t)
}
func TestIndexConsumerUnfundedStatus(t *testing.T) {
	IndexConsumerDefaultStatus(dbModule.StatusUnfunded, t)
}
