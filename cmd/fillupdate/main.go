package main

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/fillbloom"
	"github.com/notegio/openrelay/cmd/cmdutils"
	"gopkg.in/redis.v3"

	"log"
	"os"
	"os/signal"
	"strconv"
)


func main() {
	redisURL := os.Args[1]
	rpcURL := os.Args[2]
	fillSrc := os.Args[3]
	bloomURI := os.Args[4]
	// src := os.Args[3]
	// allDest := os.Args[6]
	// var changeDest string
	// if len(os.Args) >= 8 {
	// 	changeDest = os.Args[7]
	// }
	if redisURL == "" {
		log.Fatalf("Please specify redis URL")
	}
	if rpcURL == "" {
		log.Fatalf("Please specify RPC URL")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	// consumerChannel, err := channels.ConsumerFromURI(src, redisClient)
	// if err != nil { log.Fatalf(err.Error()) }
	fillConsumerChannel, err := channels.ConsumerFromURI(fillSrc, redisClient)
	if err != nil { log.Fatalf(err.Error()) }
	fillBloom, err := fillbloom.NewFillBloom(bloomURI)
	if err != nil { log.Fatalf(err.Error()) }
	lookup, err := funds.NewRPCFilledLookup(rpcURL, fillBloom)
	if err != nil { log.Fatalf(err.Error()) }
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		concurrency = 5
	}
	consumerChannels := []channels.ConsumerChannel{}
	for _, channelString := range os.Args[5:] {
		consumerChannel, allPublisher, changePublisher, err := cmdutils.ParseChannels(channelString, redisClient)
		if err != nil { log.Fatalf(err.Error()) }
		fillConsumer := funds.NewFillConsumer(allPublisher, changePublisher, lookup, concurrency)
		consumerChannels = append(consumerChannels, consumerChannel)
		consumerChannel.AddConsumer(&fillConsumer)
		consumerChannel.StartConsuming()
		log.Printf("Starting fillupdate consumer on '%v'", channelString)
	}
	fillConsumerChannel.AddConsumer(fillBloom)
	fillConsumerChannel.StartConsuming()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	for _, consumerChannel := range consumerChannels {
		consumerChannel.StopConsuming()
	}
	fillConsumerChannel.StopConsuming()
}
