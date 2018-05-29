package splitter

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/common"
)

// AddressSplitterConsumer relays orders to different channels based on the
// exchange specified in the order.
type AddressSplitterConsumer struct {
	exchanges        map[types.Address]channels.Publisher
	defaultPublisher channels.Publisher
	addrGetter       func(*types.Order) (types.Address)
	s                common.Semaphore
}

func (consumer *AddressSplitterConsumer) Consume(delivery channels.Delivery) {
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

func NewExchangeSplitterConsumer(exchanges map[types.Address]channels.Publisher, defaultPublisher channels.Publisher, concurrency int) (*AddressSplitterConsumer) {
	return &AddressSplitterConsumer{exchanges, defaultPublisher, func(order *types.Order) (types.Address) {return *order.ExchangeAddress}, make(common.Semaphore, concurrency)}
}

func NewMakerSplitterConsumer(exchanges map[types.Address]channels.Publisher, defaultPublisher channels.Publisher, concurrency int) (*AddressSplitterConsumer) {
	return &AddressSplitterConsumer{exchanges, defaultPublisher, func(order *types.Order) (types.Address) {return *order.Maker}, make(common.Semaphore, concurrency)}
}
