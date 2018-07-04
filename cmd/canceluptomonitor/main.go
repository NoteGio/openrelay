package main

import (
	"github.com/notegio/openrelay/monitor/cancelupto"
	"github.com/notegio/openrelay/channels"
	"gopkg.in/redis.v3"
	"os/signal"
	"os"
	"log"
)

func main() {
	redisURL := os.Args[1]
	rpcURL := os.Args[2]
	src := os.Args[3]
	dst := os.Args[4]
	exchangeAddress := os.Args[5]
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	consumerChannel, err := channels.ConsumerFromURI(src, redisClient)
	if err != nil {
		log.Fatalf("Error constructing consumer: %v", err.Error())
	}
	publisher, err := channels.PublisherFromURI(dst, redisClient)
	if err != nil {
		log.Fatalf("Error constructing publisher: %v", err.Error())
	}
	consumer, err := cancelupto.NewRPCCancelUpToBlockConsumer(rpcURL, exchangeAddress, publisher)
	if err != nil {
		log.Fatalf("Error constructing cancelupto monitor: %v", err.Error())
	}
	consumerChannel.AddConsumer(consumer)
	consumerChannel.StartConsuming()
	log.Printf("Started consuming blocks from channel %v for exchange %v, publishing to %v", src, exchangeAddress, dst)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	consumerChannel.StopConsuming()

}
