package main

import (
	"github.com/notegio/0xrelay/channels"
	"gopkg.in/redis.v3"
	"os"
	"os/signal"
	"log"
	"strings"
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
	var consumerChannel channels.ConsumerChannel
	if strings.HasPrefix(src, "topic://") {
		srcTopic := src[len("topic://"):]
		consumerChannel = channels.NewTopicConsumerChannel(srcTopic, redisClient)
	} else if strings.HasPrefix(src, "queue://") {
		srcQueue := src[len("queue://"):]
		consumerChannel = channels.NewQueueConsumerChannel(srcQueue, redisClient)
	} else {
		log.Fatalf("Must specify src starting with queue:// or topic://")
	}
	var publisher channels.Publisher
	if strings.HasPrefix(dest, "topic://") {
		destTopic := dest[len("topic://"):]
		publisher = channels.NewRedisTopicPublisher(destTopic, redisClient)
	} else if strings.HasPrefix(dest, "queue://") {
		destQueue := dest[len("queue://"):]
		publisher = channels.NewRedisQueuePublisher(destQueue, redisClient)
	} else {
		log.Fatalf("Must specify dest starting with queue:// or topic://")
	}

	relay := channels.NewRelay(consumerChannel, publisher, &channels.IncludeAll{})
	log.Printf("Starting")
	relay.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
}
