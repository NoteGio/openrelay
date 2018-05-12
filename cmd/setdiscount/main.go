package main

import (
	"os"
	"encoding/hex"
	"gopkg.in/redis.v3"
	"github.com/notegio/openrelay/accounts"
	"github.com/notegio/openrelay/types"
	"log"
	"strconv"
)

func main() {
	redisURL := os.Args[1]
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	addrBytes, err := hex.DecodeString(os.Args[2])
	if err != nil {
		log.Fatalf(err.Error())
	}
	address := &types.Address{}
	copy(address[:], addrBytes[:])
	expiration, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = accounts.NewRedisAccountService(redisClient).Set(address, accounts.NewAccount(
		false,
		nil,
		100,
		int64(expiration),
	))
	if err != nil {
		log.Fatalf(err.Error())
	}
}
