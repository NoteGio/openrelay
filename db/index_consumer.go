package db

import (
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	"log"
)

type IndexConsumer struct {
	idx *Indexer
}

func (consumer *IndexConsumer) Consume(msg channels.Delivery) {
	orderBytes := [441]byte{}
	copy(orderBytes[:], []byte(msg.Payload()))
	order := types.OrderFromBytes(orderBytes)
	if err := consumer.idx.Index(order); err == nil {
		msg.Ack()
	} else {
		log.Printf("Failed to index order: '%v', '%v'", order.Hash(), err.Error())
		msg.Reject()
	}
}

func NewIndexConsumer(db *gorm.DB, status int64) *IndexConsumer {
	return &IndexConsumer{NewIndexer(db, status)}
}
