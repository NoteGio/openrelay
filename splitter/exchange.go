package splitter

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/common"
)

// ExchangeSplitterConsumer relays orders to different channels based on the
// exchange specified in the order.
type ExchangeSplitterConsumer struct {
	exchanges        map[types.Address]channels.Publisher
	defaultPublisher channels.Publisher
	s                common.Semaphore
}

func (consumer *ExchangeSplitterConsumer) Consume(delivery channels.Delivery) {
	consumer.s.Acquire()
	go func(){
		defer consumer.s.Release()
		payload := delivery.Payload()
		if len(payload) == 0 {
			// Sometimes we get the odd empty message
			delivery.Ack()
			return
		}
		orderBytes := [441]byte{}
		copy(orderBytes[:], []byte(payload))
		order := types.OrderFromBytes(orderBytes)
		publisher, ok := consumer.exchanges[*order.ExchangeAddress]
		if !ok {
			publisher = consumer.defaultPublisher
		}
		publisher.Publish(payload)
		delivery.Ack()
	}()
}

func NewExchangeSplitterConsumer(exchanges map[types.Address]channels.Publisher, defaultPublisher channels.Publisher, concurrency int) (*ExchangeSplitterConsumer) {
	return &ExchangeSplitterConsumer{exchanges, defaultPublisher, make(common.Semaphore, concurrency)}
}
