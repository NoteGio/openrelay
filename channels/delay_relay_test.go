package channels_test

import (
	"github.com/notegio/openrelay/channels"
	"testing"
)

func TestDelayRelay(t *testing.T) {
	sourcePublisher, sourceChannel := channels.MockChannel()
	destPublisher, destChannel := channels.MockChannel()
	testConsumer := testConsumer{make(chan string, 5), make(chan bool), make(chan bool)}
	destChannel.AddConsumer(&testConsumer)
	destChannel.StartConsuming()
	relay := channels.NewDelayRelay(sourcePublisher, sourceChannel, destPublisher, "pause")
	relay.Start()
	defer relay.Stop()
	sourcePublisher.Publish("test1")
	sourcePublisher.Publish("test2")
	select {
	case message := <- testConsumer.channel:
		t.Errorf("Published messages shouldn't be available yet. Got '%v'", message)
	default:
	}
	relay.Flush()
	sourcePublisher.Publish("test3")
	message := <-testConsumer.channel
	if message != "test1" {
		t.Errorf("Message did not get relayed")
	}
	testConsumer.ack <- true
	<-testConsumer.done
	message = <-testConsumer.channel
	if message != "test2" {
		t.Errorf("Message did not get relayed")
	}
	testConsumer.ack <- true
	<-testConsumer.done
	select {
	case message := <- testConsumer.channel:
		t.Errorf("Should have exhausted published messages. Got '%v'", message)
	default:
	}

}
