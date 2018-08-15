package funds_test

import (
	"encoding/hex"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/funds"
	"testing"
)

func getTestOrderBytes() []byte {
	testOrderBytes, _ := hex.DecodeString("f9020994627306090abab3a6e1400e9345bc60c78a8bef57940000000000000000000000000000000000000000941dad4783cf3fe3085c1426157ab175a6119a04ba9405d090b51c40b020eab3bfcb6a2dff130df22e9ca4f47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04baa4f47261b000000000000000000000000005d090b51c40b020eab3bfcb6a2dff130df22e9c9400000000000000000000000000000000000000009490fe2af704b34e0224bf2299c838e04d4dcf1364940000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000002b5e3af16b1880000a00000000000000000000000000000000000000000000000000de0b6b3a7640000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000159938ac4a0000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222bb8421c34f27415dc0177bc4016d48c3ec7eb19ee31124bcf4ca2eb3aba767c24e4712043bf8e49d1e28c6efa5a5e8a6824886700f356a403e0e66c75621e56b184b47b03a0000000000000000000000000000000000000000000000000000000000000000080")
	return testOrderBytes
}

func TestFilledConsumer(t *testing.T) {
	sourcePublisher, consumerChannel := channels.MockChannel()
	changePublisher, changeChan := channels.MockPublisher()
	allPublisher, allChan := channels.MockPublisher()
	lookup := funds.NewMockFilledLookup(false, "0", nil)
	consumer := funds.NewFillConsumer(allPublisher, changePublisher, lookup, 1)
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

func TestFilledConsumerChange(t *testing.T) {
	sourcePublisher, consumerChannel := channels.MockChannel()
	changePublisher, changeChan := channels.MockPublisher()
	allPublisher, allChan := channels.MockPublisher()
	lookup := funds.NewMockFilledLookup(false, "2", nil)
	consumer := funds.NewFillConsumer(allPublisher, changePublisher, lookup, 1)
	consumerChannel.AddConsumer(&consumer)
	orderBytes := getTestOrderBytes()
	consumerChannel.StartConsuming()
	sourcePublisher.Publish(string(orderBytes[:]))
	updatedPayload := <-allChan
	if updatedPayload.Payload() == string(orderBytes[:]) {
		t.Errorf("Expected change in processing")
	}
	select {
	case changedPayload, ok := <-changeChan:
		if ok {
			if changedPayload.Payload() == string(orderBytes[:]) {
				t.Errorf("Expected change in processing")
			}
		} else {
			t.Errorf("Change chan was closed")
		}
	default:
		t.Errorf("Change chan should have had value")
	}

	consumerChannel.StopConsuming()
}
func TestFilledConsumerCancelChange(t *testing.T) {
	sourcePublisher, consumerChannel := channels.MockChannel()
	changePublisher, changeChan := channels.MockPublisher()
	allPublisher, allChan := channels.MockPublisher()
	lookup := funds.NewMockFilledLookup(true, "0", nil)
	consumer := funds.NewFillConsumer(allPublisher, changePublisher, lookup, 1)
	consumerChannel.AddConsumer(&consumer)
	orderBytes := getTestOrderBytes()
	consumerChannel.StartConsuming()
	sourcePublisher.Publish(string(orderBytes[:]))
	updatedPayload := <-allChan
	if updatedPayload.Payload() == string(orderBytes[:]) {
		t.Errorf("Expected change in processing")
	}
	select {
	case changedPayload, ok := <-changeChan:
		if ok {
			if changedPayload.Payload() == string(orderBytes[:]) {
				t.Errorf("Expected change in processing")
			}
		} else {
			t.Errorf("Change chan was closed")
		}
	default:
		t.Errorf("Change chan should have had value")
	}

	consumerChannel.StopConsuming()
}
