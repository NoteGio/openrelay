package channels_test

import (
	"github.com/notegio/openrelay/channels"
	"testing"
)

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
