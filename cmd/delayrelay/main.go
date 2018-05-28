package main

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/cmd/cmdutils"
	"gopkg.in/redis.v3"
	"os"
	"os/signal"
	"log"
	"strings"
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
	// src := os.Args[2]
	// dest := os.Args[3]
	// upstreamSignal := os.Args[4]
	// var downstreamSignal string
	// if len(os.Args) >= 6 {
	// 	downstreamSignal = os.Args[5]
	// } else {
	// 	downstreamSignal = ""
	// }
	if redisURL == "" {
		log.Fatalf("Please specify redis URL")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	// if strings.HasPrefix(src, "topic://") {
	// 	log.Fatal("Delay relay source must be queue, not topic")
	// }
	// sourceChannel, err := channels.ConsumerFromURI(src, redisClient)
	// if err != nil { log.Fatalf(err.Error())}
	// sourcePublisher, err := channels.PublisherFromURI(src, redisClient)
	// if err != nil { log.Fatalf(err.Error())}
	// destPublisher, err := channels.PublisherFromURI(dest, redisClient)
	// if err != nil { log.Fatalf(err.Error())}
	// signalConsumer, err := channels.ConsumerFromURI(upstreamSignal, pubsubClient)
	// if err != nil { log.Fatalf(err.Error())}
	// var signalPublisher channels.Publisher
	// if downstreamSignal != "" {
	// 	signalPublisher, err = channels.PublisherFromURI(downstreamSignal, redisClient)
	// 	if err != nil { log.Fatalf(err.Error())}
	// } else {
	// 	signalPublisher = nil
	// }
	var relays []channels.DelayRelay
	var signalConsumers []channels.ConsumerChannel
	for _, arg := range os.Args[2:] {
		channelStrings := strings.Split(arg, ",")
		if len(channelStrings) != 2 {
			log.Fatalf("Channel Strings must have two elements separated by ','")
		}
		if strings.HasPrefix(channelStrings[1], "topic://") {
			log.Fatal("Delay relay source must be queue, not topic")
		}
		signalConsumer, err := channels.ConsumerFromURI(channelStrings[0], redisClient)
		sourceChannel, publisher, signalPublisher, err := cmdutils.ParseChannels(channelStrings[1], redisClient)
		if err != nil { log.Fatalf(err.Error()) }
		relay := channels.NewDelayRelay(sourceChannel.Publisher(), sourceChannel, publisher, "pause")
		log.Printf("Starting delayrelay '%v'", arg)
		relay.Start()
		signalConsumer.AddConsumer(&DelayConsumer{&relay, signalPublisher})
		signalConsumer.StartConsuming()
		relays = append(relays, relay)
		signalConsumers = append(signalConsumers, signalConsumer)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	for _, relay := range relays {
		relay.Stop()
	}
	for _, signalConsumer := range signalConsumers {
		signalConsumer.StopConsuming()
	}
}
