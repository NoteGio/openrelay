package splitter_test

import (
	// "encoding/hex"
	// "encoding/json"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/splitter"
	"github.com/notegio/openrelay/channels"
	// "io/ioutil"
	// "reflect"
	"testing"
	"time"
	"encoding/json"
	// "bytes"
	"io/ioutil"
)

func getTestOrderBytes(t *testing.T) []byte {
	return sampleOrder(t).Bytes()
}

func sampleOrder(t *testing.T) *types.Order {
	order := &types.Order{}
	if orderData, err := ioutil.ReadFile("../formatted_transaction.json"); err == nil {
		if err := json.Unmarshal(orderData, order); err != nil {
			t.Fatalf(err.Error())
		}
	}
	return order
}

func TestExchangeSplitter(t *testing.T) {
	sourcePublisher, sourceConsumerChannel := channels.MockChannel()
	defaultPublisher, defaultConsumerChannel := channels.MockChannel()
	otherPublisher, otherConsumerChannel := channels.MockChannel()
	mapping := make(map[types.Address]channels.Publisher)
	mapping[*sampleOrder(t).ExchangeAddress] = otherPublisher
	exchangeSplitter := splitter.NewExchangeSplitterConsumer(mapping, defaultPublisher, 1)
	sourceConsumerChannel.AddConsumer(exchangeSplitter)
	sourceConsumerChannel.StartConsuming()
	defaultConsumerChannel.StartConsuming()
	otherConsumerChannel.StartConsuming()
	orderBytes := getTestOrderBytes(t)
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
	mapping[*sampleOrder(t).Taker] = otherPublisher
	exchangeSplitter := splitter.NewExchangeSplitterConsumer(mapping, defaultPublisher, 1)
	sourceConsumerChannel.AddConsumer(exchangeSplitter)
	sourceConsumerChannel.StartConsuming()
	defaultConsumerChannel.StartConsuming()
	otherConsumerChannel.StartConsuming()
	orderBytes := getTestOrderBytes(t)
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
