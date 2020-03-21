package main

import (
	"github.com/notegio/openrelay/pool"
	"github.com/notegio/openrelay/ingest"
	"github.com/notegio/openrelay/channels"
	dbModule "github.com/notegio/openrelay/db"
	"net/http"
	"gopkg.in/redis.v3"
	"os"
	"log"
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

// pgurl, pgpass, redis, dstChannel

func main() {
	db, err := dbModule.GetDB(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err.Error())
	}
	redisURL := os.Args[3]
	dstChannel := os.Args[4]
	port := "8080"
	if redisURL == "" {
		log.Fatalf("Please specify redis URL")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	publisher, err := channels.PublisherFromURI(dstChannel, redisClient)
	if err != nil { log.Fatalf(err.Error()) }
	handler := pool.PoolDecoratorBaseFee(db, redisClient, pool.PoolAdminHandler(db, publisher))

	mux := &regexpHandler{[]*route{}}
	mux.HandleFunc(regexp.MustCompile("^(/[^/]+)?/v3/_admin$"), handler)
	mux.HandleFunc(regexp.MustCompile("^/_hc$"), ingest.HealthCheckHandler(redisClient))
	log.Printf("Order Ingest Serving on :%v", port)
	http.ListenAndServe(":"+port, mux)
}
