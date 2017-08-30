package affiliates_test

import (
	"encoding/hex"
	"github.com/notegio/0xrelay/affiliates"
	"github.com/notegio/0xrelay/config"
	"gopkg.in/redis.v3"
	"math/big"
	"os"
	"testing"
	// "time"
)

func getRedisClient(t *testing.T) *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		t.Errorf("Please set the REDIS_URL environment variable")
		return nil
	}
	return redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
}

func TestGetMissingAffiliate(t *testing.T) {
	redisClient := getRedisClient(t)
	if redisClient == nil {
		return
	}
	service := affiliates.NewRedisAffiliateService(redisClient)
	address, _ := hex.DecodeString("1000000000000000000000000000000000000000")
	var addressArray [20]byte
	copy(addressArray[:], address[:])
	_, err := service.Get(addressArray)
	if err == nil {
		t.Errorf("Missing affiliate should return error")
		return
	}
}

func TestSetAffiliate(t *testing.T) {
	redisClient := getRedisClient(t)
	if redisClient == nil {
		return
	}
	baseFee := config.NewBaseFee(redisClient)
	if err := baseFee.Set(big.NewInt(10000)); err != nil {
		t.Errorf(err.Error())
		return
	}
	service := affiliates.NewRedisAffiliateService(redisClient)
	affiliate := affiliates.NewAffiliate(new(big.Int), 100)
	address, _ := hex.DecodeString("0000000000000000000000000000000000000000")
	var addressArray [20]byte
	copy(addressArray[:], address[:])
	err := service.Set(addressArray, affiliate)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	lookedUpAffiliate, err := service.Get(addressArray)
	fee, _ := baseFee.Get()
	if lookedUpAffiliate.Fee().Cmp(fee) != 0 {
		t.Errorf("Fee should be equal to base fee")
	}
}
