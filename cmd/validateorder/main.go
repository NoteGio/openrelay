package main

import (
	"os"
	"github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/config"
	"gopkg.in/redis.v3"
	"log"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func hexToBytes(address string) [20]byte {
	slice, err := hex.DecodeString(address)
	if err != nil {
		return [20]byte{}
	}
	output := [20]byte{}
	copy(output[:], slice[:])
	return output
}

func main() {
	rpcURL := os.Args[1]
	redisURL := os.Args[2]
	orderFile := os.Args[3]
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	fundChecker, err := funds.NewRpcOrderValidator(rpcURL, config.NewFeeToken(redisClient), config.NewTokenProxy(redisClient))
	if err != nil {
		log.Fatalf(err.Error())
	}
	newOrder := types.Order{}
	if orderData, err := ioutil.ReadFile(orderFile); err == nil {
		if err := json.Unmarshal(orderData, &newOrder); err != nil {
			log.Fatalf(err.Error())
		}
	}
	fmt.Println(fundChecker.ValidateOrder(&newOrder))
}
