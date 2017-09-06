package channels

import (
	"log"
)

// RelayFilter objects provide a predicate function to determine whether a
// message should be passed to the next stage
type RelayFilter interface {
	Filter(Delivery) bool
}

type IncludeAll struct{
	counter int64
}

func (filter *IncludeAll) Filter(delivery Delivery) bool {
	filter.counter++;
	log.Printf("Relayed message : '%v'", filter.counter)
	return true
}

type Relay struct {
	consumerChannel ConsumerChannel
	publisher       Publisher
	filter          RelayFilter
}

func (relay *Relay) Start() bool {
	return relay.consumerChannel.StartConsuming()
}

func (relay *Relay) Stop() bool {
	return relay.consumerChannel.StopConsuming()
}

type RelayConsumer struct {
	relay *Relay
}

func (consumer *RelayConsumer) Consume(delivery Delivery) {
	if consumer.relay.filter.Filter(delivery) {
		consumer.relay.publisher.Publish(delivery.Payload())
	}
	delivery.Ack()
}

func NewRelay(channel ConsumerChannel, publisher Publisher, filter RelayFilter) Relay {
	relay := Relay{
		channel,
		publisher,
		filter,
	}
	relay.consumerChannel.AddConsumer(&RelayConsumer{&relay})
	return relay
}
