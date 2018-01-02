package splitter

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
)

// ExchangeSplitterConsumer relays orders to different channels based on the
// exchange specified in the order.
type ExchangeSplitterConsumer struct {
	exchanges map[types.Address]channels.Publisher
	defaultPublisher channels.Publisher
}

func (consumer *ExchangeSplitterConsumer) Consume(delivery channels.Delivery) {
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
}

func NewExchangeSplitterConsumer(exchanges map[types.Address]channels.Publisher, defaultPublisher channels.Publisher) (*ExchangeSplitterConsumer) {
	return &ExchangeSplitterConsumer{exchanges, defaultPublisher}
}
