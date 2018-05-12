package accounts

import (
	"encoding/json"
	"github.com/notegio/openrelay/config"
	"github.com/notegio/openrelay/types"
	"gopkg.in/redis.v3"
	"math/big"
	"log"
	"fmt"
)

type redisAccountService struct {
	redisClient *redis.Client
	baseFee     config.BaseFee
}

func (accountService *redisAccountService) Get(address *types.Address) Account {
	log.Printf("Getting account")
	acct := &account{false, new(big.Int), 0, 0}
	acctJSON, err := accountService.redisClient.Get(fmt.Sprintf("account::%x", address[:])).Result()
	if err != nil {
		log.Printf("Error getting account: %v", err.Error())
		// Account not found, return the default value
		return acct
	}
	fee, err := accountService.baseFee.Get()
	if err != nil {
		log.Printf("Error getting base fee: %v", err.Error())
		// If we can't get the base fee, we can't calculate a discount, so
		// we'll return the default account.
		return acct
	}
	json.Unmarshal([]byte(acctJSON), acct)
	acct.BaseFee = fee
	return acct
}

func (accountService *redisAccountService) Set(address *types.Address, acct Account) error {
	data, err := json.Marshal(acct)
	if err != nil {
		return err
	}
	return accountService.redisClient.Set(fmt.Sprintf("account::%x", address[:]), string(data), 0).Err()
}

func NewRedisAccountService(redisClient *redis.Client) AccountService {
	return &redisAccountService{
		redisClient,
		config.NewBaseFee(redisClient),
	}
}
