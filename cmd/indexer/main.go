package main

import (
	"github.com/notegio/openrelay/channels"
	dbModule "github.com/notegio/openrelay/db"
	// "github.com/notegio/openrelay/funds"
	"gopkg.in/redis.v3"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"os/signal"
	"io/ioutil"
	"fmt"
)

func main() {
	redisURL := os.Args[1]
	srcChannel := os.Args[2]
	pgHost := os.Args[3]
	pgUser := os.Args[4]
	pgPassword := ""
	status := dbModule.StatusOpen
	for _, arg := range os.Args[5:] {
		if arg == "--unfunded" {
			status = dbModule.StatusUnfunded
		} else {
			pgPasswordFile := arg
			pgPasswordBytes, err := ioutil.ReadFile(pgPasswordFile)
			if err != nil {
				log.Fatalf("Could not read password file: %v", err.Error());
			}
			pgPassword = string(pgPasswordBytes)
		}
	}
	if pgPassword == "" {
		pgPassword = os.Getenv("POSTGRES_PASSWORD")
	}
	connectionString := fmt.Sprintf(
		"host=%v dbname=postgres sslmode=disable user=%v password=%v",
		pgHost,
		pgUser,
		pgPassword,
	)
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Could not open postgres connection: %v", err.Error())
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
