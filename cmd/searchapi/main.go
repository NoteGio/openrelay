package main

import (
	"github.com/notegio/openrelay/search"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/blockhash"
	"github.com/notegio/openrelay/affiliates"
	dbModule "github.com/notegio/openrelay/db"
	"net/http"
	"gopkg.in/redis.v3"
	"os"
	"log"
	"github.com/rs/cors"
	"strconv"
)

func main() {
	redisURL := os.Args[1]
	blockChannel := os.Args[2]
	db, err := dbModule.GetDB(os.Args[3], os.Args[4])
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err.Error())
	}
	port := "8080"
	for _, arg := range os.Args[5:] {
		if _, err := strconv.Atoi(arg); err == nil {
			// If the argument is castable as an integer,
			port = arg
		}
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
	orderBookHandler := search.BlockHashDecorator(blockHash, search.OrderBookHandler(db))
	feeRecipientsHandler := search.BlockHashDecorator(blockHash, search.FeeRecipientHandler(affiliates.NewRedisAffiliateService(redisClient)))
	pairHandler := search.PairHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/v2/orders", searchHandler)
	mux.HandleFunc("/v2/order/", orderHandler)
	mux.HandleFunc("/v2/asset_pairs", pairHandler)
	mux.HandleFunc("/v2/orderbook", orderBookHandler)
	mux.HandleFunc("/v2/fee_recipients", feeRecipientsHandler)
	mux.HandleFunc("/_hc", search.HealthCheckHandler(db, blockHash))
	corsHandler := cors.Default().Handler(mux)
	log.Printf("Order Search Serving on :%v", port)
	http.ListenAndServe(":"+port, corsHandler)
}
