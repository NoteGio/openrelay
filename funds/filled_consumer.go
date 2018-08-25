package funds

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/common"

	"log"
	"reflect"
)

func doLookup(order *types.Order, oldValue *types.Uint256, lookupFn func(*types.Order) (*types.Uint256, error), valChan chan *types.Uint256, changeChan chan bool) {
	value, err := lookupFn(order)
	if err != nil {
		log.Printf(err.Error())
		value = oldValue
	}
	valChan <- value
	changeChan <- !reflect.DeepEqual(value, oldValue)
}

func doLookupBoolean(order *types.Order, oldValue bool, lookupFn func(*types.Order) (bool, error), valChan chan bool, changeChan chan bool) {
	value, err := lookupFn(order)
	if err != nil {
		log.Printf(err.Error())
		value = oldValue
	}
	valChan <- value
	changeChan <- value != oldValue
}


type FillConsumer struct {
	allPublisher    channels.Publisher
	changePublisher channels.Publisher
	lookup          FilledLookup
	s               common.Semaphore
}

func (consumer *FillConsumer) Consume(msg channels.Delivery) {
	go func() {
		consumer.s.Acquire()
		defer consumer.s.Release()
		order, err := types.OrderFromBytes([]byte(msg.Payload()))
		if err != nil {
			log.Printf("Error parsing order: %#x", msg.Payload())
			msg.Reject()
		}
		cancelledChan := make(chan bool)
		filledChan := make(chan *types.Uint256)
		changes := make(chan bool, 2)
		go doLookupBoolean(order, order.Cancelled, consumer.lookup.GetCancelled, cancelledChan, changes)
		go doLookup(order, order.TakerAssetAmountFilled, consumer.lookup.GetAmountFilled, filledChan, changes)
		order.Cancelled = <-cancelledChan || order.Cancelled
		order.TakerAssetAmountFilled = <-filledChan
		payload := string(order.Bytes())
		if (<-changes || <-changes) && consumer.changePublisher != nil {
			consumer.changePublisher.Publish(payload)
		}
		consumer.allPublisher.Publish(payload)
		msg.Ack()
	}()
}

func NewFillConsumer(allPublisher channels.Publisher, changePublisher channels.Publisher, lookup FilledLookup, concurrency int) FillConsumer {
	return FillConsumer{allPublisher, changePublisher, lookup, make(common.Semaphore, concurrency)}
}
