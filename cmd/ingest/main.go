package main

import (
	"github.com/notegio/openrelay/ingest"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/affiliates"
	"github.com/notegio/openrelay/accounts"
	dbModule "github.com/notegio/openrelay/db"
	"encoding/hex"
	"net/http"
	"gopkg.in/redis.v3"
	"os"
	"log"
	"github.com/rs/cors"
)

func main() {
	db, err := dbModule.GetDB(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err.Error())
	}
	redisURL := os.Args[3]
	defaultFeeRecipientString := os.Args[4]
	dstChannel := os.Args[5]
	var port string
	if len(os.Args) >= 7 {
		port = os.Args[6]
	} else {
		port = "8080"
	}
	if redisURL == "" {
		log.Fatalf("Please specify redis URL")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	defaultFeeRecipientSlice, err := hex.DecodeString(defaultFeeRecipientString)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defaultFeeRecipientBytes := [20]byte{}
	copy(defaultFeeRecipientBytes[:], defaultFeeRecipientSlice[:])
	affiliateService := affiliates.NewRedisAffiliateService(redisClient)
	accountService := accounts.NewRedisAccountService(redisClient)
	publisher, err := channels.PublisherFromURI(dstChannel, redisClient)
	enforceTerms := os.Getenv("OR_ENFORCE_TERMS") != "false"
	if err != nil { log.Fatalf(err.Error()) }
	handler := ingest.Handler(publisher, accountService, affiliateService, enforceTerms, dbModule.NewTermsManager(db), dbModule.NewExchangeLookup(db))
	feeHandler := ingest.FeeHandler(publisher, accountService, affiliateService, defaultFeeRecipientBytes)

	mux := http.NewServeMux()
	mux.HandleFunc("/v2/order", handler)
	mux.HandleFunc("/v2/order_config", feeHandler)
	mux.HandleFunc("/_hc", ingest.HealthCheckHandler(redisClient))
	corsHandler := cors.Default().Handler(mux)
	log.Printf("Order Ingest Serving on :%v", port)
	http.ListenAndServe(":"+port, corsHandler)
}
