package main

import (
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/funds"
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
	db, err := dbModule.GetDB(os.Args[3], os.Args[4])
	if err != nil { log.Fatalf("Error opening database: %v", err.Error()) }
	if redisURL == "" {
		log.Fatalf("Please specify redis URL")
	}
	if rpcURL == "" {
		log.Fatalf("Please specify RPC URL")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	lookup:= funds.NewDBCancellationLookup(db)
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		concurrency = 5
	}
	consumerChannels := []channels.ConsumerChannel{}
	for _, channelString := range os.Args[5:] {
		consumerChannel, allPublisher, changePublisher, err := cmdutils.ParseChannels(channelString, redisClient)
		if err != nil { log.Fatalf(err.Error()) }
		fillConsumer := funds.NewCancellationConsumer(allPublisher, changePublisher, lookup, concurrency)
		consumerChannels = append(consumerChannels, consumerChannel)
		consumerChannel.AddConsumer(&fillConsumer)
		consumerChannel.StartConsuming()
		log.Printf("Starting cancellation updater consumer on '%v'", channelString)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	for _, consumerChannel := range consumerChannels {
		consumerChannel.StopConsuming()
	}
}
