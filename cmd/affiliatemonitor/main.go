package main

import (
	"github.com/notegio/openrelay/monitor/affiliate"
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
	affiliateSignupAddress := os.Args[4]
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	consumerChannel, err := channels.ConsumerFromURI(src, redisClient)
	if err != nil {
		log.Fatalf("Error constructing consumer: %v", err.Error())
	}
	consumer, err := affiliate.NewRPCAffiliateBlockConsumer(rpcURL, affiliateSignupAddress, redisClient)
	if err != nil {
		log.Fatalf("Error constructing affiliate monitor: %v", err.Error())
	}
	consumerChannel.AddConsumer(consumer)
	consumerChannel.StartConsuming()
	log.Printf("Started consuming blocks from channel %v for signUp %v", src, affiliateSignupAddress)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	consumerChannel.StopConsuming()

}
