package main

import(
	"gopkg.in/redis.v3"
	"os"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/splitter"
	"strings"
	"log"
	"os/signal"
)

func main() {
	redisURL := os.Args[1]
	src := os.Args[2]
	defaultDest := os.Args[3]
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	sourceConsumerChannel, err := channels.ConsumerFromURI(src, redisClient)
	if err != nil { log.Fatalf(err.Error()) }
	mapping := make(map[types.Address]channels.Publisher)
	defaultPublisher, err := channels.PublisherFromURI(defaultDest, redisClient)
	if err != nil { log.Fatalf(err.Error()) }
	for _, value := range os.Args[4:] {
		splitValues := strings.Split(value, "=")
		addressHex := splitValues[0]
		publisher, err := channels.PublisherFromURI(splitValues[1], redisClient)
		if err != nil { log.Fatalf(err.Error()) }
		addressBytes, err := common.HexToBytes(addressHex)
		if err != nil {
			log.Fatalf(err.Error())
		}
		mapping[*common.BytesToOrAddress(addressBytes)] = publisher
	}
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		concurrency = 5
	}
	exchangeSplitter := splitter.NewExchangeSplitterConsumer(mapping, defaultPublisher, concurrency)
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
