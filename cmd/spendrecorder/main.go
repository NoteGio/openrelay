package main

import (
	"github.com/notegio/openrelay/channels"
	dbModule "github.com/notegio/openrelay/db"
	// "github.com/notegio/openrelay/funds"
	"gopkg.in/redis.v3"
	"log"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	redisURL := os.Args[1]
	srcChannel := os.Args[2]
	db, err := dbModule.GetDB(os.Args[3], os.Args[4])
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err.Error())
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	destChannel := os.Args[5]
	consumerChannel, err := channels.ConsumerFromURI(srcChannel, redisClient)
	if err != nil {
		log.Fatalf("Error establishing consumer channel: %v", err.Error())
	}
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		concurrency = 5
	}
	publisher, err := channels.PublisherFromURI(destChannel, redisClient)
	if err != nil {
		log.Fatalf("Error establishing publisher channel: %v", err.Error())
	}
	consumerChannel.AddConsumer(dbModule.NewRecordSpendConsumer(db, concurrency, publisher))
	consumerChannel.StartConsuming()
	log.Printf("Starting spend recorder consumer on '%v'", srcChannel)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	consumerChannel.StopConsuming()
}
