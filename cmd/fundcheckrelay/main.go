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
	"strings"
	"strconv"
)

type FundFilter struct {
	orderValidator funds.OrderValidator
}

func (filter *FundFilter) Filter(delivery channels.Delivery) bool {
	msg := []byte(delivery.Payload())
	orderBytes := [441]byte{}
	copy(orderBytes[:], msg[:])
	order := types.OrderFromBytes(orderBytes)
	if !order.Signature.Verify(order.Maker) {
		log.Printf("Invalid order signature");
		return false;
	}
	valid, _ := filter.orderValidator.ValidateOrder(order)
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
	if redisURL == "" {
		log.Fatalf("Please specify redis URL")
	}
	if rpcURL == "" {
		log.Fatalf("Please specify RPC URL")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	publishers := []channels.Publisher{}
	consumerChannel, err := channels.ConsumerFromURI(src, redisClient)
	if err != nil { log.Fatalf(err.Error()) }
	invert := false
	var invalidationChannel channels.ConsumerChannel
	for _, arg := range os.Args[4:] {
		if arg == "--invert" {
			invert = true
		} else if strings.HasPrefix(arg, "--invalidation=") {
			arg = strings.TrimPrefix(arg, "--invalidation=")
			invalidationChannel, err = channels.ConsumerFromURI(arg, redisClient)
			if err != nil { log.Fatalf(err.Error()) }
		} else {
			publisher, err := channels.PublisherFromURI(arg, redisClient)
			if err != nil { log.Fatalf(err.Error()) }
			publishers = append(publishers, publisher)
		}
	}
	feeToken, err := config.NewRpcFeeToken(rpcURL)
	if err != nil {
		log.Fatalf("Error creating RpcOrderValidator: '%v'", err.Error())
	}
	tokenProxy, err := config.NewRpcTokenProxy(rpcURL)
	if err != nil {
		log.Fatalf("Error creating RpcOrderValidator: '%v'", err.Error())
	}
	orderValidator, err := funds.NewRpcOrderValidator(rpcURL, feeToken, tokenProxy, invalidationChannel)
	if err != nil {
		log.Fatalf("Error creating RpcOrderValidator: '%v'", err.Error())
	}
	var fundFilter channels.RelayFilter
	fundFilter = &FundFilter{orderValidator}
	if invert {
		fundFilter = &channels.InvertFilter{fundFilter}
	}
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		concurrency = 5
	}

	relay := channels.NewRelay(consumerChannel, publishers, fundFilter, concurrency)
	log.Printf("Starting fundcheck")
	relay.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	relay.Stop()
}
