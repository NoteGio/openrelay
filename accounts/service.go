package accounts

import (
	"github.com/notegio/0xrelay/config"
	"gopkg.in/redis.v3"
	"encoding/json"
	"math/big"
)

type redisAccountService struct {
	redisClient *redis.Client
	baseFee config.BaseFee
}

func (accountService *redisAccountService) Get(address [20]byte) Account {
	acct := &account{false, new(big.Int), 0, 0}
	acctJSON, err := accountService.redisClient.Get("account::" + string(address[:])).Result()
	if err != nil {
		// Account not found, return the default value
		return acct
	}
	fee, err := accountService.baseFee.Get()
	if err != nil {
		// If we can't get the base fee, we can't calculate a discount, so
		// we'll return the default account.
		return acct
	}
	json.Unmarshal([]byte(acctJSON), acct)
	acct.baseFee = fee
	return acct
}

func (accountService *redisAccountService) Set(address [20]byte, acct Account) error {
	data, err := json.Marshal(acct)
	if err != nil {
		return err
	}
	return accountService.redisClient.Set("account::" + string(address[:]), string(data), 0).Err()
}

func NewRedisAccountService(redisClient *redis.Client) (AccountService) {
	return &redisAccountService{
		redisClient,
		config.NewBaseFee(redisClient),
	}
}
