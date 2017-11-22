package funds

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"

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

type FillConsumer struct {
	allPublisher    channels.Publisher
	changePublisher channels.Publisher
	lookup          FilledLookup
}

func (consumer *FillConsumer) Consume(msg channels.Delivery) {
	orderBytes := [441]byte{}
	copy(orderBytes[:], []byte(msg.Payload()))
	order := types.OrderFromBytes(orderBytes)
	cancelledChan := make(chan *types.Uint256)
	filledChan := make(chan *types.Uint256)
	changes := make(chan bool, 2)
	go doLookup(order, order.TakerTokenAmountCancelled, consumer.lookup.GetAmountCancelled, cancelledChan, changes)
	go doLookup(order, order.TakerTokenAmountFilled, consumer.lookup.GetAmountFilled, filledChan, changes)
	order.TakerTokenAmountCancelled = <-cancelledChan
	order.TakerTokenAmountFilled = <-filledChan
	orderBytes = order.Bytes()
	payload := string(orderBytes[:])
	consumer.allPublisher.Publish(payload)
	if (<-changes || <-changes) && consumer.changePublisher != nil {
		consumer.changePublisher.Publish(payload)
	}
	msg.Ack()
}

func NewFillConsumer(allPublisher channels.Publisher, changePublisher channels.Publisher, lookup FilledLookup) FillConsumer {
	return FillConsumer{allPublisher, changePublisher, lookup}
}
