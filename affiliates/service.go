package affiliates

import (
	"encoding/json"
	"github.com/notegio/openrelay/config"
	"github.com/notegio/openrelay/types"
	"gopkg.in/redis.v3"
	"math/big"
	"fmt"
)

type redisAffiliateService struct {
	redisClient *redis.Client
	baseFee     config.BaseFee
}

func (affiliateService *redisAffiliateService) Get(address *types.Address) (Affiliate, error) {
	acct := &affiliate{new(big.Int), 100}
	acctJSON, err := affiliateService.redisClient.Get(fmt.Sprintf("affiliate::%x", address[:])).Result()
	if err != nil {
		// Affiliate not found, return the default value
		return nil, err
	}
	fee, err := affiliateService.baseFee.Get()
	if err != nil {
		// If we can't get the base fee, we can't calculate a discount, so
		// we'll return the default affiliate.
		return nil, err
	}
	json.Unmarshal([]byte(acctJSON), acct)
	acct.BaseFee = fee
	return acct, nil
}

func (affiliateService *redisAffiliateService) Set(address *types.Address, acct Affiliate) error {
	data, err := json.Marshal(acct)
	if err != nil {
		return err
	}
	return affiliateService.redisClient.Set(fmt.Sprintf("affiliate::%x", address[:]), string(data), 0).Err()
}

func NewRedisAffiliateService(redisClient *redis.Client) AffiliateService {
	return &redisAffiliateService{
		redisClient,
		config.NewBaseFee(redisClient),
	}
}
