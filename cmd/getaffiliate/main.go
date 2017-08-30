package main

import (
	"encoding/hex"
	"gopkg.in/redis.v3"
	"github.com/notegio/openrelay/affiliates"
	"os"
	"fmt"
)

func main() {
	redisURL := os.Args[1]
	address := os.Args[2]
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	affiliateService := affiliates.NewRedisAffiliateService(redisClient)
	if addressBytes, err := hex.DecodeString(address); err == nil {
		addressArray := [20]byte{}
		copy(addressArray[:], addressBytes[:])
		affiliate := affiliates.Get(addressArray)
		data, _ := json.Marshal(affiliate)
		fmt.Println(string(data))
	}
}
