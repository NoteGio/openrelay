package main

import (
	"github.com/notegio/openrelay/channels"
	"gopkg.in/redis.v3"
	"os"
	"os/signal"
	"log"
)

func main() {
	redisURL := os.Args[1]
	src := os.Args[2]
	dest := os.Args[3]
	if redisURL == "" {
		log.Fatalf("Please specify redis URL")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	consumerChannel, err := channels.ConsumerFromURI(src, redisClient)
	if err != nil { log.Fatalf(err.Error())}
	publisher, err := channels.PublisherFromURI(dest, redisClient)
	if err != nil { log.Fatalf(err.Error())}
	relay := channels.NewRelay(consumerChannel, publisher, &channels.IncludeAll{})
	log.Printf("Starting")
	relay.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	relay.Stop()
}
