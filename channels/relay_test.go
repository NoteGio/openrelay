package channels_test

import (
	"github.com/notegio/openrelay/channels"
	"testing"
)

type TestFilter struct { }

func (filter *TestFilter) Filter(delivery channels.Delivery) bool {
	return delivery.Payload() == "test"
}

func TestRelay(t *testing.T) {
	sourcePublisher, sourceChannel := channels.MockChannel()
	destPublisher, destChannel := channels.MockChannel()
	testConsumer := testConsumer{make(chan string), make(chan bool), make(chan bool)}
	destChannel.AddConsumer(&testConsumer)
	destChannel.StartConsuming()
	relay := channels.NewRelay(sourceChannel, destPublisher, &channels.IncludeAll{})
	relay.Start()
	defer relay.Stop()
	sourcePublisher.Publish("test")
	message := <-testConsumer.channel
	if message != "test" {
		t.Errorf("Message did not get relayed")
	}
}

func TestInvertFilter(t *testing.T) {
	sourcePublisher, sourceChannel := channels.MockChannel()
	destPublisher, destChannel := channels.MockChannel()
	testConsumer := testConsumer{make(chan string), make(chan bool), make(chan bool)}
	destChannel.AddConsumer(&testConsumer)
	destChannel.StartConsuming()
	relay := channels.NewRelay(sourceChannel, destPublisher, &channels.InvertFilter{&TestFilter{}})
	relay.Start()
	defer relay.Stop()
	sourcePublisher.Publish("test")
	sourcePublisher.Publish("abc")
	message := <-testConsumer.channel
if message != "abc" {
		t.Errorf("Message did not get relayed")
	}
}
