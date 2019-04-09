package db

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/common"
	"log"
)

type RecordFillConsumer struct {
	idx *Indexer
	s   common.Semaphore
}

func (consumer *RecordFillConsumer) Consume(msg channels.Delivery) {
	consumer.s.Acquire()
	go func (){
		defer consumer.s.Release()
		fillRecord := &FillRecord{}
		json.Unmarshal([]byte(msg.Payload()), fillRecord)
		if err := consumer.idx.RecordFill(fillRecord); err == nil {
			msg.Ack()
			} else {
				log.Printf("Failed to record fill: '%v', '%v'", fillRecord.OrderHash, err.Error())
				msg.Reject()
			}
	}()
}

func NewRecordFillConsumer(db *gorm.DB, concurrency int, publisher channels.Publisher) *RecordFillConsumer {
	log.Printf("Publisher %v", publisher)
	return &RecordFillConsumer{NewIndexer(db, StatusOpen, publisher), make(common.Semaphore, concurrency)}
}
