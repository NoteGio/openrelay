package funds_test

import (
	"encoding/hex"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/funds"
	"testing"
)

func getTestOrderBytes() [441]byte {
	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000059938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	var testOrderByteArray [441]byte
	copy(testOrderByteArray[:], testOrderBytes[:])
	return testOrderByteArray
}

func TestFilledConsumer(t *testing.T) {
	sourcePublisher, consumerChannel := channels.MockChannel()
	changePublisher, changeChan := channels.MockPublisher()
	allPublisher, allChan := channels.MockPublisher()
	lookup := funds.NewMockFilledLookup("0", "0", nil)
	consumer := funds.NewFillConsumer(allPublisher, changePublisher, lookup)
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
	lookup := funds.NewMockFilledLookup("1", "1", nil)
	consumer := funds.NewFillConsumer(allPublisher, changePublisher, lookup)
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
