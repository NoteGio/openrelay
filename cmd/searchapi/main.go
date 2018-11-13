package main

import (
	"github.com/notegio/openrelay/search"
	"github.com/notegio/openrelay/pool"
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
	searchHandler := search.BlockHashDecorator(blockHash, pool.PoolDecorator(search.SearchHandler(db)))
	orderHandler := search.BlockHashDecorator(blockHash, search.OrderHandler(db))
	orderBookHandler := search.BlockHashDecorator(blockHash, search.OrderBookHandler(db))
	feeRecipientsHandler := search.BlockHashDecorator(blockHash, search.FeeRecipientHandler(affiliates.NewRedisAffiliateService(redisClient)))
	pairHandler := search.PairHandler(db)

	mux := &regexpHandler{[]*route{}}
	mux.HandleFunc(regexp.MustCompile("^(/[^/]+)?/v2/orders$"), searchHandler)
	mux.HandleFunc(regexp.MustCompile("^(/[^/]+)?/v2/order/$"), orderHandler)
	mux.HandleFunc(regexp.MustCompile("^(/[^/]+)?/v2/asset_pairs$"), pairHandler)
	mux.HandleFunc(regexp.MustCompile("^(/[^/]+)?/v2/orderbook$"), orderBookHandler)
	mux.HandleFunc(regexp.MustCompile("^(/[^/]+)?/v2/fee_recipients$"), feeRecipientsHandler)
	mux.HandleFunc(regexp.MustCompile("^/_hc$"), search.HealthCheckHandler(db, blockHash))
	corsHandler := cors.Default().Handler(mux)
	log.Printf("Order Search Serving on :%v", port)
	http.ListenAndServe(":"+port, corsHandler)
}
