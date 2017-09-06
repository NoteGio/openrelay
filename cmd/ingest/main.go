package main

import (
	"github.com/notegio/openrelay/ingest"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/affiliates"
	"github.com/notegio/openrelay/accounts"
	"net/http"
	"gopkg.in/redis.v3"
	"os"
	"log"
)

func main() {
	redisURL := os.Args[1]
	var port string
	if len(os.Args) >= 3 {
		port = os.Args[2]
	} else {
		port = "8080"
	}
	if redisURL == "" {
		log.Fatalf("Please specify redis URL")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	affiliateService := affiliates.NewRedisAffiliateService(redisClient)
	accountService := accounts.NewRedisAccountService(redisClient)
	publisher := channels.NewRedisQueuePublisher("ingest", redisClient)
	handler := ingest.Handler(publisher, accountService, affiliateService)

    http.HandleFunc("/", handler)
	log.Printf("Serving on :%v", port)
    http.ListenAndServe(":"+port, nil)
}
