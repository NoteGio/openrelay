package main

import (
	"github.com/notegio/openrelay/search"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/blockhash"
	"net/http"
	"io/ioutil"
	"gopkg.in/redis.v3"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"log"
	"github.com/rs/cors"
	"fmt"
	"strconv"
)

func main() {
	redisURL := os.Args[1]
	blockChannel := os.Args[2]
	pgHost := os.Args[3]
	pgUser := os.Args[4]
	pgPassword := ""
	port := "8080"
	for _, arg := range os.Args[5:] {
		if _, err := strconv.Atoi(arg); err == nil {
			// If the argument is castable as an integer,
			port = arg
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
		"host=%v sslmode=disable user=%v password=%v",
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
	blockChannelConsumer, err := channels.ConsumerFromURI(blockChannel, redisClient)
	if err != nil {
		log.Fatalf("Error establishing block channel: %v", err.Error())
	}
	blockHash := blockhash.NewChanneledBlockHash(blockChannelConsumer)
	searchHandler := search.BlockHashDecorator(blockHash, search.SearchHandler(db))
	orderHandler := search.BlockHashDecorator(blockHash, search.OrderHandler(db))
	pairHandler := search.PairHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/v0/orders", searchHandler)
	mux.HandleFunc("/v0/order/", orderHandler)
	mux.HandleFunc("/v0/token_pairs", pairHandler)
	corsHandler := cors.Default().Handler(mux)
	log.Printf("Order Search Serving on :%v", port)
	http.ListenAndServe(":"+port, corsHandler)
}
