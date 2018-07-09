package funds

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/common"

	"log"
)

type CancellationConsumer struct {
	allPublisher    channels.Publisher
	changePublisher channels.Publisher
	lookup          CancellationLookup
	s               common.Semaphore
}

func (consumer *CancellationConsumer) Consume(msg channels.Delivery) {
	go func() {
		consumer.s.Acquire()
		defer consumer.s.Release()
		order, err := types.OrderFromBytes([]byte(msg.Payload()))
		if err != nil {
			log.Printf("Error parsing order: %#x", msg.Payload())
			msg.Reject()
			return
		}
		oldCancelled := order.Cancelled
		order.Cancelled, err = consumer.lookup.GetCancelled(order)
		if err != nil {
			log.Printf("Error gettin cancelled status: %v", err.Error())
			msg.Reject()
			return
		}
		payload := string(order.Bytes())
		if oldCancelled != order.Cancelled && consumer.changePublisher != nil {
			consumer.changePublisher.Publish(payload)
		}
		consumer.allPublisher.Publish(payload)
		msg.Ack()
	}()
}

func NewCancellationConsumer(allPublisher channels.Publisher, changePublisher channels.Publisher, lookup CancellationLookup, concurrency int) CancellationConsumer {
	return CancellationConsumer{allPublisher, changePublisher, lookup, make(common.Semaphore, concurrency)}
}
