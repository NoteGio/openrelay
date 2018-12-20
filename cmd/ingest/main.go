package main

import (
	"github.com/notegio/openrelay/ingest"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/affiliates"
	"github.com/notegio/openrelay/accounts"
	"github.com/notegio/openrelay/pool"
	dbModule "github.com/notegio/openrelay/db"
	"encoding/hex"
	"net/http"
	"gopkg.in/redis.v3"
	"os"
	"log"
	"github.com/rs/cors"
	"regexp"
)

type route struct {
    pattern *regexp.Regexp
    handler http.Handler
}

type regexpHandler struct {
    routes []*route
}

func (h *regexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
    h.routes = append(h.routes, &route{pattern, handler})
}

func (h *regexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
    h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *regexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    for _, route := range h.routes {
        if route.pattern.MatchString(r.URL.Path) {
            route.handler.ServeHTTP(w, r)
            return
        }
    }
    // no pattern matched; send 404 response
    http.NotFound(w, r)
}

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
	exchangeLookup := dbModule.NewExchangeLookup(db)
	handler := pool.PoolDecoratorBaseFee(db, redisClient, ingest.Handler(publisher, accountService, affiliateService, enforceTerms, dbModule.NewTermsManager(db), exchangeLookup))
	feeHandler := pool.PoolDecoratorBaseFee(db, redisClient, ingest.FeeHandler(publisher, accountService, affiliateService, defaultFeeRecipientBytes, exchangeLookup))

	mux := &regexpHandler{[]*route{}}
	mux.HandleFunc(regexp.MustCompile("^(/[^/]+)?/v2/order$"), handler)
	mux.HandleFunc(regexp.MustCompile("^(/[^/]+)?/v2/order_config$"), feeHandler)
	mux.HandleFunc(regexp.MustCompile("^/_hc$"), ingest.HealthCheckHandler(redisClient))
	corsHandler := cors.Default().Handler(mux)
	log.Printf("Order Ingest Serving on :%v", port)
	http.ListenAndServe(":"+port, corsHandler)
}
