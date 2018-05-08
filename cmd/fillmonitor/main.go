package main

import (
	"github.com/notegio/openrelay/monitor/fill"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/fillbloom"
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
	storageURI := os.Args[5]
	exchangeAddress := os.Args[6]
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
	fillBloom, err := fillbloom.NewFillBloom(storageURI)
	if err != nil {
		log.Fatalf("Error constructing fillbloom: %v", err.Error())
	}
	consumer, err := fill.NewRPCFillBlockConsumer(rpcURL, exchangeAddress, publisher, fillBloom)
	if err != nil {
		log.Fatalf("Error constructing fill monitor: %v", err.Error())
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
