package main

import(
	"gopkg.in/redis.v3"
	"os"
	// "github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/splitter"
	"log"
	"os/signal"
	"strconv"
)

func main() {
	redisURL := os.Args[1]
	src := os.Args[2]
	suffix := os.Args[3]
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	sourceConsumerChannel, err := channels.ConsumerFromURI(src, redisClient)
	if err != nil { log.Fatalf(err.Error()) }
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		concurrency = 5
	}
	translator := channels.NewRedisURITraslator(redisClient)
	exchangeSplitter := splitter.NewExchangeSplitterConsumer(translator, suffix, concurrency)
	sourceConsumerChannel.AddConsumer(exchangeSplitter)
	sourceConsumerChannel.StartConsuming()
	log.Printf("Consuming on '%v'", src)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	sourceConsumerChannel.StopConsuming()
}
