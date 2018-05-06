package main

import (
	"github.com/notegio/openrelay/monitor/multisig"
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
	multisigAddress := os.Args[4]
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	consumerChannel, err := channels.ConsumerFromURI(src, redisClient)
	if err != nil {
		log.Fatalf("Error constructing consumer: %v", err.Error())
	}
	consumer, err := multisig.NewRPCMultisigBlockConsumer(rpcURL, multisigAddress)
	if err != nil {
		log.Fatalf("Error constructing multisig monitor: %v", err.Error())
	}
	consumerChannel.AddConsumer(consumer)
	consumerChannel.StartConsuming()
	log.Printf("Started consuming blocks from channel %v for exchange %v", src, multisigAddress)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	consumerChannel.StopConsuming()

}
