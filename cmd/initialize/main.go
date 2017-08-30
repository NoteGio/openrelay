package main

import (
	"encoding/hex"
	"gopkg.in/redis.v3"
	"github.com/notegio/openrelay/config"
	"github.com/notegio/openrelay/affiliates"
	"math/big"
	"os"
	"fmt"
)

func main() {
	redisURL := os.Args[1]
	baseFeeString := os.Args[2]
	authorizedAddresses := os.Args[2:]
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	baseFeeService := config.NewBaseFee(redisClient)
	baseFeeInt := new(big.Int)
	baseFeeInt.SetString(baseFeeString, 10)
	baseFeeService.Set(baseFeeInt)

	affiliateService := affiliates.NewRedisAffiliateService(redisClient)
	for _, address := range(authorizedAddresses) {
		if addressBytes, err := hex.DecodeString(address); err == nil {
			addressArray := [20]byte{}
			copy(addressArray[:], addressBytes[:])
			affiliate := affiliates.NewAffiliate(baseFeeInt, 100)
			affiliateService.Set(addressArray, affiliate)
			fmt.Printf("Added address '%v'\n", hex.EncodeToString(addressArray[:]))
		}
	}
}
