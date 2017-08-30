package main

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/types"
	"gopkg.in/redis.v3"
	"os"
	"os/signal"
	"log"
)

type FundFilter struct {
	orderValidator funds.OrderValidator
}

func (filter *FundFilter) Filter(delivery channels.Delivery) bool {
	msg := []byte(delivery.Payload())
	orderBytes := [377]byte{}
	copy(orderBytes[:], msg[:])
	order := types.OrderFromBytes(orderBytes)
	return filter.orderValidator.ValidateOrder(order)
}

func main() {
	redisURL := os.Args[1]
	rpcURL := os.Args[2]
	if redisURL == "" {
		log.Fatalf("Please specify redis URL")
	}
	if rpcURL == "" {
		log.Fatalf("Please specify RPC URL")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	consumerChannel := channels.NewQueueConsumerChannel("ingest", redisClient)
	publisher := channels.NewRedisTopicPublisher("instant", redisClient)
	orderValidator, err := funds.NewRpcOrderValidator(rpcURL)
	if err != nil { log.Fatalf(err.Error()) }
	fundFilter := &FundFilter{orderValidator}
	relay := channels.NewRelay(consumerChannel, publisher, fundFilter)
	log.Printf("Starting")
	relay.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
}
