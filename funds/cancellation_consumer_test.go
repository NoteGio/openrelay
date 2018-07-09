package funds_test

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/types"
	"testing"
)

func TestCancellationConsumer(t *testing.T) {
	sourcePublisher, consumerChannel := channels.MockChannel()
	changePublisher, changeChan := channels.MockPublisher()
	allPublisher, allChan := channels.MockPublisher()
	lookup := funds.NewMockCancellationLookup(false)
	consumer := funds.NewCancellationConsumer(allPublisher, changePublisher, lookup, 1)
	consumerChannel.AddConsumer(&consumer)
	orderBytes := getTestOrderBytes()
	consumerChannel.StartConsuming()
	sourcePublisher.Publish(string(orderBytes[:]))
	updatedPayload := <-allChan
	if updatedPayload.Payload() != string(orderBytes[:]) {
		t.Errorf("Unexpected change in processing")
	}
	select {
	case _, ok := <-changeChan:
		if ok {
			t.Errorf("Change chan should have been empty")
		} else {
			t.Errorf("Change chan was closed")
		}
	default:
	}

	consumerChannel.StopConsuming()
}

func TestCancellationChangeConsumer(t *testing.T) {
	sourcePublisher, consumerChannel := channels.MockChannel()
	changePublisher, changeChan := channels.MockPublisher()
	allPublisher, allChan := channels.MockPublisher()
	lookup := funds.NewMockCancellationLookup(true)
	consumer := funds.NewCancellationConsumer(allPublisher, changePublisher, lookup, 1)
	consumerChannel.AddConsumer(&consumer)
	orderBytes := getTestOrderBytes()
	consumerChannel.StartConsuming()
	sourcePublisher.Publish(string(orderBytes[:]))
	updatedPayload := <-allChan
	order, err := types.OrderFromBytes([]byte(updatedPayload.Payload()))
	if err != nil { t.Errorf("Error: %v", err.Error()) }
	if !order.Cancelled {
		t.Errorf("Expected order to be cancelld")
	}
	select {
	case _, ok := <-changeChan:
		if !ok {
			t.Errorf("Change chan should have had a message")
		}
	default:
	}

	consumerChannel.StopConsuming()
}
