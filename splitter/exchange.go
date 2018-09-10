package splitter

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/common"
	"log"
	"fmt"
)

// AddressSplitterConsumer relays orders to different channels based on the
// exchange specified in the order.
type AddressSplitterConsumer struct {
	translator       channels.URITranslator
	suffix           string
	attrGetter       func(*types.Order) (interface{})
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
		order, err := types.OrderFromBytes([]byte(payload))
		if err != nil {
			log.Printf("Error parsing order: %v", err.Error())
			delivery.Reject()
			return
		}
		publisher, err := consumer.translator.PublisherFromURI(fmt.Sprintf("queue://%v-%v", consumer.attrGetter(order), consumer.suffix))
		if err != nil {
			log.Printf("Error producing publisher for %v", consumer.attrGetter(order))
			delivery.Reject()
			return
		}
		publisher.Publish(payload)
		delivery.Ack()
	}()
}

func NewExchangeSplitterConsumer(translator channels.URITranslator, suffix string, concurrency int) (*AddressSplitterConsumer) {
	return &AddressSplitterConsumer{translator, suffix, func(order *types.Order) (interface{}) {return order.ExchangeAddress}, make(common.Semaphore, concurrency)}
}

func NewMakerSplitterConsumer(translator channels.URITranslator, suffix string, concurrency int) (*AddressSplitterConsumer) {
	return &AddressSplitterConsumer{translator, suffix, func(order *types.Order) (interface{}) {return order.Maker}, make(common.Semaphore, concurrency)}
}
