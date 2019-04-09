package db

import (
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/common"
	"log"
)

type IndexConsumer struct {
	idx *Indexer
	s   common.Semaphore
}

func (consumer *IndexConsumer) Consume(msg channels.Delivery) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Failed to index order: %v", r)
			msg.Reject()
		}
	}()
	consumer.s.Acquire()
	go func(){
		defer consumer.s.Release()
		order, err := types.OrderFromBytes([]byte(msg.Payload()))
		if err != nil {
			log.Printf("Error parsing order: %v", err.Error())
			msg.Reject()
			return
		}
		if err := consumer.idx.Index(order); err == nil {
			msg.Ack()
		} else {
			log.Printf("Failed to index order: '%v', '%v'", order.Hash(), err.Error())
			msg.Reject()
		}
	}()
}

func NewIndexConsumer(db *gorm.DB, status int64, concurrency int, publisher channels.Publisher) *IndexConsumer {
	return &IndexConsumer{NewIndexer(db, status, publisher), make(common.Semaphore, concurrency)}
}
