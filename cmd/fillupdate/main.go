package main

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/funds"
	"gopkg.in/redis.v3"

	"log"
	"os"
	"os/signal"
)


func main() {
	redisURL := os.Args[1]
	rpcURL := os.Args[2]
	src := os.Args[3]
	allDest := os.Args[4]
	var changeDest string
	if len(os.Args) >= 6 {
		changeDest = os.Args[5]
	}
	if redisURL == "" {
		log.Fatalf("Please specify redis URL")
	}
	if rpcURL == "" {
		log.Fatalf("Please specify RPC URL")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	consumerChannel, err := channels.ConsumerFromURI(src, redisClient)
	if err != nil { log.Fatalf(err.Error()) }
	allPublisher, err := channels.PublisherFromURI(allDest, redisClient)
	if err != nil { log.Fatalf(err.Error()) }
	var changePublisher channels.Publisher
	if changeDest != "" {
		changePublisher, err = channels.PublisherFromURI(changeDest, redisClient)
		if err != nil { log.Fatalf(err.Error()) }
	}
	lookup, err := funds.NewRpcFilledLookup(rpcURL)
	if err != nil { log.Fatalf(err.Error()) }
	fillConsumer := funds.NewFillConsumer(allPublisher, changePublisher, lookup)
	consumerChannel.AddConsumer(&fillConsumer)
	consumerChannel.StartConsuming()
	log.Printf("Starting fillupdate consumer on '%v'", src)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	consumerChannel.StopConsuming()
}
