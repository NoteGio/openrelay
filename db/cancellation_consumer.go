package db

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/common"
	"log"
)

type CancellationConsumer struct {
	idx *Indexer
	s   common.Semaphore
}

func (consumer *CancellationConsumer) Consume(msg channels.Delivery) {
	consumer.s.Acquire()
	go func(){
		defer consumer.s.Release()
		cancellation := &Cancellation{}
		err := json.Unmarshal([]byte(msg.Payload()), cancellation)
		if err != nil {
			log.Printf("Failed to parse JSON: %v", err.Error())
			msg.Reject()
			return
		}
		if err := consumer.idx.RecordCancellation(cancellation); err == nil {
			msg.Ack()
		} else {
			log.Printf("Failed to record cancellation: '%v', '%v'", msg.Payload(), err.Error())
			msg.Reject()
			return
		}
	}()
}

func NewRecordCancellationConsumer(db *gorm.DB, concurrency int, publisher channels.Publisher) *CancellationConsumer {
	return &CancellationConsumer{NewIndexer(db, StatusCancelled, publisher), make(common.Semaphore, concurrency)}
}
