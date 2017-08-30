package main

import (
	"github.com/notegio/openrelay/channels"
	"gopkg.in/redis.v3"
	"os"
	"os/signal"
	"log"
)

// DelayConsumer Flushes the DelayRelay every time it receives a message. If
// the publisher is specified, it will also pass the message it received to the
// publisher.
type DelayConsumer struct {
	relay *channels.DelayRelay
	publisher channels.Publisher
}

func (consumer *DelayConsumer) Consume(delivery channels.Delivery) {
	consumer.relay.Flush()
	if consumer.publisher != nil {
		consumer.publisher.Publish(delivery.Payload())
	}
	delivery.Ack()
}

func main() {
	redisURL := os.Args[1]
	src := os.Args[2]
	dest := os.Args[3]
	upstreamSignal := os.Args[4]
	var downstreamSignal string
	if len(os.Args) >= 6 {
		downstreamSignal = os.Args[5]
	} else {
		downstreamSignal = ""
	}
	if redisURL == "" {
		log.Fatalf("Please specify redis URL")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	sourceChannel, err := channels.ConsumerFromURI(src, redisClient)
	if err != nil { log.Fatalf(err.Error())}
	sourcePublisher, err := channels.PublisherFromURI(src, redisClient)
	if err != nil { log.Fatalf(err.Error())}
	destPublisher, err := channels.PublisherFromURI(dest, redisClient)
	if err != nil { log.Fatalf(err.Error())}
	signalConsumer, err := channels.ConsumerFromURI(upstreamSignal, redisClient)
	if err != nil { log.Fatalf(err.Error())}
	var signalPublisher channels.Publisher
	if downstreamSignal != "" {
		signalPublisher, err = channels.PublisherFromURI(downstreamSignal, redisClient)
		if err != nil { log.Fatalf(err.Error())}
	} else {
		signalPublisher = nil
	}

	relay := channels.NewDelayRelay(sourcePublisher, sourceChannel, destPublisher, "pause")
	log.Printf("started")
	relay.Start()
	signalConsumer.AddConsumer(&DelayConsumer{&relay, signalPublisher})
	signalConsumer.StartConsuming()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	relay.Stop()
	signalConsumer.StopConsuming()
}
