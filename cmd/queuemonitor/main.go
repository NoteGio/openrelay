package main

import (
	"os"
	"log"
	"fmt"
	"gopkg.in/redis.v3"
	"strconv"
	"time"
)

func main() {
	counts := make(map[string]int64)
	redisURL := os.Args[1]
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	threshInt, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf(err.Error())
	}
	threshold := int64(threshInt)
	for _, queue := range os.Args[3:] {
		for _, modifier := range []string{"", "::unacked", "::rejected"} {
			k := fmt.Sprintf("%v%v", queue, modifier)
			counts[k], err = redisClient.LLen(k).Result()
			if err != nil {
				log.Fatalf(err.Error())
			}
			log.Printf("Initial Queue: %v - %v", k, counts[k])
		}
	}
	counter := 0
	for {
		time.Sleep(20 * time.Second)
		for k, v := range counts {
			counts[k], err = redisClient.LLen(k).Result()
			if err != nil {
				log.Fatalf(err.Error())
			}
			if (counts[k] / threshold) > (v / threshold) {
				log.Printf("Queue Increasing: %v - %v", k, counts[k])
			} else if (counts[k] / threshold) < (v / threshold) {
				log.Printf("Queue Decreasing: %v - %v", k, counts[k])
			} else if counter % 3 == 0 {
				// Print all the queues once a minute
				log.Printf("Queue Steady: %v - %v", k, counts[k])
			}
		}
		counter++
	}
}
