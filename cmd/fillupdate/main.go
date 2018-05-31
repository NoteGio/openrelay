package main

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/fillbloom"
	"gopkg.in/redis.v3"

	"log"
	"os"
	"os/signal"
	"strconv"
)


func main() {
	redisURL := os.Args[1]
	rpcURL := os.Args[2]
	src := os.Args[3]
	fillSrc := os.Args[4]
	bloomURI := os.Args[5]
	allDest := os.Args[6]
	var changeDest string
	if len(os.Args) >= 8 {
		changeDest = os.Args[7]
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
	fillConsumerChannel, err := channels.ConsumerFromURI(fillSrc, redisClient)
	if err != nil { log.Fatalf(err.Error()) }
	fillBloom, err := fillbloom.NewFillBloom(bloomURI)
	if err != nil { log.Fatalf(err.Error()) }
	allPublisher, err := channels.PublisherFromURI(allDest, redisClient)
	if err != nil { log.Fatalf(err.Error()) }
	var changePublisher channels.Publisher
	if changeDest != "" {
		changePublisher, err = channels.PublisherFromURI(changeDest, redisClient)
		if err != nil { log.Fatalf(err.Error()) }
	}
	lookup, err := funds.NewRPCFilledLookup(rpcURL, fillBloom)
	if err != nil { log.Fatalf(err.Error()) }
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		concurrency = 5
	}
	fillConsumer := funds.NewFillConsumer(allPublisher, changePublisher, lookup, concurrency)
	consumerChannel.AddConsumer(&fillConsumer)
	consumerChannel.StartConsuming()
	fillConsumerChannel.AddConsumer(fillBloom)
	fillConsumerChannel.StartConsuming()
	log.Printf("Starting fillupdate consumer on '%v'", src)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	consumerChannel.StopConsuming()
	fillConsumerChannel.StopConsuming()
}
