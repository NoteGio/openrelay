package db

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/channels"
	"log"
)

type RecordFillConsumer struct {
	idx *Indexer
}

func (consumer *RecordFillConsumer) Consume(msg channels.Delivery) {
	fillRecord := &FillRecord{}
	json.Unmarshal([]byte(msg.Payload()), fillRecord)
	if err := consumer.idx.RecordFill(fillRecord); err == nil {
		msg.Ack()
	} else {
		log.Printf("Failed to record fill: '%v', '%v'", fillRecord.OrderHash, err.Error())
		msg.Reject()
	}
}

func NewRecordFillConsumer(db *gorm.DB) *RecordFillConsumer {
	return &RecordFillConsumer{NewIndexer(db, StatusOpen)}
}
