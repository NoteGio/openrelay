package main

import (
	"github.com/notegio/openrelay/channels"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/common"
	// "github.com/notegio/openrelay/funds"
	"gopkg.in/redis.v3"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"os/signal"
	"fmt"
)

func main() {
	redisURL := os.Args[1]
	srcChannel := os.Args[2]
	db, err := dbModule.GetDB(os.Args[3], os.args[4])
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err.Error())
	}
	status := dbModule.StatusOpen
	for _, arg := range os.Args[5:] {
		if arg == "--unfunded" {
			status = dbModule.StatusUnfunded
		}
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	consumerChannel, err := channels.ConsumerFromURI(srcChannel, redisClient)
	if err != nil {
		log.Fatalf("Error establishing consumer channel: %v", err.Error())
	}
	consumerChannel.AddConsumer(dbModule.NewIndexConsumer(db, status))
	consumerChannel.StartConsuming()
	log.Printf("Starting db indexer consumer on '%v'", srcChannel)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	consumerChannel.StopConsuming()
}
