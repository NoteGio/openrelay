package affiliates

import (
	"encoding/json"
	"github.com/notegio/openrelay/config"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/common"
	"gopkg.in/redis.v3"
	"math/big"
	"fmt"
	"log"
	"strings"
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

func (affiliateService *redisAffiliateService) List() ([]types.Address, error) {
	addresses := []types.Address{}
	addressesHex, err := affiliateService.redisClient.Keys("affiliate::*").Result()
	if err != nil {
		return addresses, err
	}
	for _, addressHex := range addressesHex {
		if address, err := common.HexToAddress(strings.TrimPrefix(addressHex, "affiliate::")); err == nil {
			addresses = append(addresses, *address)
		} else {
			log.Printf("Invalid affiliate address: %v", addressHex)
		}
	}
	return addresses, nil
}

func NewRedisAffiliateService(redisClient *redis.Client) AffiliateService {
	return &redisAffiliateService{
		redisClient,
		config.NewBaseFee(redisClient),
	}
}
