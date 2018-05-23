package channels

import (
	"github.com/notegio/openrelay/common"
)

type DelayRelay struct {
	*Relay
	sourcePublisher Publisher
	sentinel        string
	delayChan       chan bool
}

func (relay *DelayRelay) Flush() {
	relay.sourcePublisher.Publish(relay.sentinel)
	relay.delayChan <- true
}

type DelayRelayFilter struct {
	sentinel  string
	delayChan chan bool
}

func (filter *DelayRelayFilter) Filter(delivery Delivery) bool {
	if delivery.Payload() == filter.sentinel {
		delivery.Ack()
		<-filter.delayChan
		return false
	}
	return true
}

// NewDelayRelay creates a DelayRelay that consumes from `channel` and relays
// messages to `publisher`. It also requires `sourcePublisher`, which must be
// able to publish messages to `channel`, and `sentinel` which should be a
// string that would never be a normal message from the consumer channel.
// Messages will pile up on the DelayRelay until DelayRelay.Flush() is called,
// at which point any queued messages will flush to `publisher`.
//
// Note: `channel` should be a queue based channel, not a topic based channel.
// Topic based channels process messages in parallel, and thus won't block
// until Flush() is called
func NewDelayRelay(sourcePublisher Publisher, channel ConsumerChannel, publisher Publisher, sentinel string) DelayRelay {
	delayChan := make(chan bool)
	relay := DelayRelay{
		&Relay{
			channel,
			[]Publisher{publisher},
			&DelayRelayFilter{sentinel, delayChan},
			make(common.Semaphore, 1), // DelayRelays can't handle concurrency > 1
		},
		sourcePublisher,
		sentinel,
		delayChan,
	}
	relay.consumerChannel.AddConsumer(&RelayConsumer{relay.Relay})
	sourcePublisher.Publish(sentinel)
	return relay
}
