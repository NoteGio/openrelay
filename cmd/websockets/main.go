package main

import (
	"github.com/notegio/openrelay/subscriptions"
	"github.com/notegio/openrelay/channels"
	dbModule "github.com/notegio/openrelay/db"
	"net/http"
	"gopkg.in/redis.v3"
	"os"
	"os/signal"
	"log"
	"strconv"
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
	redisURL := os.Args[1]
	orderChannel := os.Args[2]
	db, err := dbModule.GetDB(os.Args[3], os.Args[4])
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err.Error())
	}
	port := uint(8080)
	for _, arg := range os.Args[5:] {
		if portConv, err := strconv.Atoi(arg); err == nil {
			// If the argument is castable as an integer,
			port = uint(portConv)
		}
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	orderChannelConsumer, err := channels.ConsumerFromURI(orderChannel, redisClient)
	if err != nil {
		log.Fatalf("Error establishing block channel: %v", err.Error())
	}
	manager := subscriptions.NewWebsocketSubscriptionManager()
	quit, err := manager.ListenForSubscriptions(port, db)
	if err != nil {
		log.Fatalf("Error listening for subscriptions: %v", err.Error())
	}
	orderChannelConsumer.AddConsumer(manager)
	orderChannelConsumer.StartConsuming()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	orderChannelConsumer.StopConsuming()
	quit()
}
