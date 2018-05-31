package channels

import (
	"time"
	"errors"
)

type mockConsumerChannel struct {
	channel   chan Delivery
	consumers []Consumer
	unacked   *deliveries
	rejected  *deliveries
	processed uint
	publisher *mockPublisher
}

func (mock *mockConsumerChannel) AddConsumer(consumer Consumer) bool {
	mock.consumers = append(mock.consumers, consumer)
	return true
}

func (mock *mockConsumerChannel) StartConsuming() bool {
	go func() {
		for message := range mock.channel {
			mock.unacked.deliveries = append(mock.unacked.deliveries, message)
			for _, consumer := range mock.consumers {
				consumer.Consume(message)
			}
			mock.processed++
		}
	}()
	return true
}

func (mock *mockConsumerChannel) StopConsuming() bool {
	close(mock.channel)
	return true
}

func (mock *mockConsumerChannel) ReturnAllUnacked() int {
	returned := len(mock.unacked.deliveries)
	go func() {
		var message Delivery
		for {
			index := len(mock.unacked.deliveries) - 1
			if index >= 0 {
				message, mock.unacked.deliveries = mock.unacked.deliveries[index], mock.unacked.deliveries[:index]
				mock.channel <- message
			}
		}
	}()
	return returned
}

func (mock *mockConsumerChannel) PurgeRejected() int {
	rejected := len(mock.rejected.deliveries)
	mock.rejected.deliveries = []Delivery{}
	return rejected
}

func (mock *mockConsumerChannel) Publisher() Publisher {
	return mock.publisher
}

type mockPublisher struct {
	channel  chan Delivery
	unacked  *deliveries
	rejected *deliveries
}

func (mock *mockPublisher) Publish(payload string) bool {
	mock.channel <- &mockDelivery{payload, mock.unacked, mock.rejected}
	return true
}

type deliveries struct {
	deliveries []Delivery
}

type mockDelivery struct {
	payload  string
	unacked  *deliveries
	rejected *deliveries
}

func (mock *mockDelivery) Payload() string {
	return mock.payload
}

func (mock *mockDelivery) Ack() bool {
	for i, value := range mock.unacked.deliveries {
		if value == mock {
			mock.unacked.deliveries = append(
				mock.unacked.deliveries[:i],
				mock.unacked.deliveries[i+1:]...,
			)
			return true
		}
	}
	return false
}
func (mock *mockDelivery) Reject() bool {
	for i, value := range mock.unacked.deliveries {
		if value == mock {
			mock.unacked.deliveries = append(
				mock.unacked.deliveries[:i],
				mock.unacked.deliveries[i+1:]...,
			)
			mock.rejected.deliveries = append(mock.rejected.deliveries, mock)
			return true
		}
	}
	return false
}

func (mock *mockDelivery) Return() bool {
	return true
}

func MockChannel() (Publisher, ConsumerChannel) {
	channel := make(chan Delivery, 5)
	unacked := &deliveries{}
	rejected := &deliveries{}
	pub := &mockPublisher{
		channel,
		unacked,
		rejected,
	}
	consume := &mockConsumerChannel{
		channel,
		[]Consumer{},
		unacked,
		rejected,
		0,
		pub,
	}
	return pub, consume
}

func MockPublisher() (Publisher, chan Delivery) {
	channel := make(chan Delivery, 5)
	unacked := &deliveries{}
	rejected := &deliveries{}
	return &mockPublisher{
		channel,
		unacked,
		rejected,
	}, channel
}

func MockFinish(channel ConsumerChannel, count uint) error {
	switch mockChannel := channel.(type) {
	case *mockConsumerChannel:
		for mockChannel.processed < count {
			time.Sleep(10 * time.Millisecond)
		}
		return nil
	default:
		return errors.New("Channel not a mockConsumerChannel")
	}
}
