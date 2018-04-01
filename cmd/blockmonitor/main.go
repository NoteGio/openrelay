package main

import (
	"github.com/notegio/openrelay/monitor/blocks"
	"github.com/notegio/openrelay/channels"
	"gopkg.in/redis.v3"
	"os/signal"
	"os"
	"log"
	"strconv"
	"time"
	"fmt"
	"strings"
)

func main() {
	redisURL := os.Args[1]
	rpcURL := os.Args[2]
	dst := os.Args[3]
	brbSize := 200
	pollInterval := 3*time.Second
	var err error
	if len(os.Args) > 4 {
		brbSize, err = strconv.Atoi(os.Args[4])
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	if len(os.Args) > 5 {
		pollIntervalInt, err := strconv.Atoi(os.Args[5])
		if err != nil {
			log.Fatalf(err.Error())
		}
		pollInterval = time.Duration(pollIntervalInt) * time.Second
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	publisher, err := channels.PublisherFromURI(dst, redisClient)
	if err != nil {
		log.Fatalf("Error constructing publisher: %v", err.Error())
	}
	blockRecorder := blocks.NewRedisBlockRecorder(redisClient, fmt.Sprintf("%v::blocknumber", strings.Split(dst, "://")[1]))
	monitor, err := blocks.NewRPCBlockMonitor(rpcURL, publisher, pollInterval, blockRecorder, brbSize)
	if err != nil {
		log.Fatalf("Error constructing monitor: %v", err.Error())
	}
	go func() {
		err := monitor.Process()
		if err != nil {
			log.Fatalf("Processing error: %v", err.Error())
		}
	}()
	log.Printf("Block Monitor: Started block monitor. RPC Host: '%v'. Queue: '%v'", rpcURL, dst)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	monitor.Stop()

}
