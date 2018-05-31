package splitter_test

import (
	"encoding/hex"
	// "encoding/json"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/splitter"
	"github.com/notegio/openrelay/channels"
	// "io/ioutil"
	// "reflect"
	"testing"
	"time"
)

func getTestOrderBytes() [441]byte {
	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	var testOrderByteArray [441]byte
	copy(testOrderByteArray[:], testOrderBytes[:])
	return testOrderByteArray
}

func sampleOrder() *types.Order {
	order := &types.Order{}
	order.FromBytes(getTestOrderBytes())
	return order
}

func TestExchangeSplitter(t *testing.T) {
	sourcePublisher, sourceConsumerChannel := channels.MockChannel()
	defaultPublisher, defaultConsumerChannel := channels.MockChannel()
	otherPublisher, otherConsumerChannel := channels.MockChannel()
	mapping := make(map[types.Address]channels.Publisher)
	mapping[*sampleOrder().ExchangeAddress] = otherPublisher
	exchangeSplitter := splitter.NewExchangeSplitterConsumer(mapping, defaultPublisher, 1)
	sourceConsumerChannel.AddConsumer(exchangeSplitter)
	sourceConsumerChannel.StartConsuming()
	defaultConsumerChannel.StartConsuming()
	otherConsumerChannel.StartConsuming()
	orderBytes := getTestOrderBytes()
	sourcePublisher.Publish(string(orderBytes[:]))
	time.Sleep(200 * time.Millisecond)
	if unackedCount := sourceConsumerChannel.ReturnAllUnacked(); unackedCount != 0 {
		t.Errorf("Expected 0 unacked on source, got %v", unackedCount)
	}
	if unackedCount := defaultConsumerChannel.ReturnAllUnacked(); unackedCount != 0 {
		t.Errorf("Expected 0 unacked on default, got %v", unackedCount)
	}
	if unackedCount := otherConsumerChannel.ReturnAllUnacked(); unackedCount != 1 {
		t.Errorf("Expected 1 unacked on other, got %v", unackedCount)
	}
}
func TestExchangeSplitterDefault(t *testing.T) {
	sourcePublisher, sourceConsumerChannel := channels.MockChannel()
	defaultPublisher, defaultConsumerChannel := channels.MockChannel()
	otherPublisher, otherConsumerChannel := channels.MockChannel()
	mapping := make(map[types.Address]channels.Publisher)
	mapping[*sampleOrder().Taker] = otherPublisher
	exchangeSplitter := splitter.NewExchangeSplitterConsumer(mapping, defaultPublisher, 1)
	sourceConsumerChannel.AddConsumer(exchangeSplitter)
	sourceConsumerChannel.StartConsuming()
	defaultConsumerChannel.StartConsuming()
	otherConsumerChannel.StartConsuming()
	orderBytes := getTestOrderBytes()
	sourcePublisher.Publish(string(orderBytes[:]))
	time.Sleep(200 * time.Millisecond)
	if unackedCount := sourceConsumerChannel.ReturnAllUnacked(); unackedCount != 0 {
		t.Errorf("Expected 0 unacked on source, got %v", unackedCount)
	}
	if unackedCount := defaultConsumerChannel.ReturnAllUnacked(); unackedCount != 1 {
		t.Errorf("Expected 1 unacked on default, got %v", unackedCount)
	}
	if unackedCount := otherConsumerChannel.ReturnAllUnacked(); unackedCount != 0 {
		t.Errorf("Expected 0 unacked on other, got %v", unackedCount)
	}
}
