package main

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/config"
	"gopkg.in/redis.v3"
	"os"
	"os/signal"
	"encoding/hex"
	"log"
)

type FundFilter struct {
	orderValidator funds.OrderValidator
}

func (filter *FundFilter) Filter(delivery channels.Delivery) bool {
	msg := []byte(delivery.Payload())
	orderBytes := [409]byte{}
	copy(orderBytes[:], msg[:])
	order := types.OrderFromBytes(orderBytes)
	if !order.Signature.Verify(order.Maker) {
		log.Printf("Invalid order signature");
		return false;
	}
	valid := filter.orderValidator.ValidateOrder(order)
	if valid {
		log.Printf("Order '%v' has funds", hex.EncodeToString(order.Hash()))
	} else {
		log.Printf("Order '%v' lacks funds", hex.EncodeToString(order.Hash()))
	}
	return valid
}

func main() {
	redisURL := os.Args[1]
	rpcURL := os.Args[2]
	src := os.Args[3]
	dest := os.Args[4]
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
	publisher, err := channels.PublisherFromURI(dest, redisClient)
	if err != nil { log.Fatalf(err.Error()) }
	orderValidator, err := funds.NewRpcOrderValidator(rpcURL, config.NewFeeToken(redisClient), config.NewTokenProxy(redisClient))
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
	relay.Stop()
}
