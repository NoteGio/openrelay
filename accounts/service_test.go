package accounts_test

import (
	"encoding/hex"
	"github.com/notegio/openrelay/accounts"
	"github.com/notegio/openrelay/config"
	"github.com/notegio/openrelay/types"
	"gopkg.in/redis.v3"
	"math/big"
	"os"
	"testing"
	"time"
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

func TestGetDefaultAccount(t *testing.T) {
	redisClient := getRedisClient(t)
	if redisClient == nil {
		return
	}
	service := accounts.NewRedisAccountService(redisClient)
	account := service.Get(&types.Address{})
	if account.Blacklisted() {
		t.Errorf("Default account should not be blacklisted")
		return
	}
}

func TestSetAccount(t *testing.T) {
	redisClient := getRedisClient(t)
	if redisClient == nil {
		return
	}
	baseFee := config.NewBaseFee(redisClient)
	if err := baseFee.Set(big.NewInt(10000)); err != nil {
		t.Errorf(err.Error())
		return
	}
	service := accounts.NewRedisAccountService(redisClient)
	account := accounts.NewAccount(false, new(big.Int), 0, time.Now().Unix()+5)
	address, _ := hex.DecodeString("0000000000000000000000000000000000000000")
	addressArray := &types.Address{}
	copy(addressArray[:], address[:])
	err := service.Set(addressArray, account)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	lookedUpAccount := service.Get(addressArray)
	if lookedUpAccount.Blacklisted() {
		t.Errorf("Expected blacklisted to be false")
	}
}
