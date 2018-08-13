package affiliate

import (
	"encoding/json"
	"math/big"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/affiliates"
	"github.com/notegio/openrelay/monitor/blocks"
	"gopkg.in/redis.v3"
	"log"
)

type affiliateBlockConsumer struct {
	affiliateSignupAddress   *big.Int
	affiliateService  affiliates.AffiliateService
	logFilter         ethereum.LogFilterer
}

func (consumer *affiliateBlockConsumer) Consume(delivery channels.Delivery) {
	block := &blocks.MiniBlock{}
	err := json.Unmarshal([]byte(delivery.Payload()), block)
	if err != nil {
		log.Printf("Error parsing payload: %v\n", err.Error())
	}
	affiliateTopic := &big.Int{}
	affiliateTopic.SetString("60dad0d232381238c031553102e3a2d779bda5a9507ec806820542b3da2801eb", 16)
	if block.Bloom.Test(consumer.affiliateSignupAddress) && block.Bloom.Test(affiliateTopic) {
		query := ethereum.FilterQuery{
			FromBlock: block.Number,
			ToBlock: block.Number,
			Addresses: []common.Address{common.BigToAddress(consumer.affiliateSignupAddress)},
			Topics: [][]common.Hash{
				[]common.Hash{common.BigToHash(affiliateTopic)},
			},
		}
		log.Printf("Block %v - %#x bloom filter indicates affiliate logs", block.Number, block.Hash)
		logs, err := consumer.logFilter.FilterLogs(context.Background(), query)
		if err != nil {
			delivery.Return()
			log.Fatalf("Failed to filter logs on block %v - aborting: %v", block.Number, err.Error())
		}
		log.Printf("Found %v affiliate logs", len(logs))
		for _, affiliateLog := range logs {
			affiliate := affiliates.NewAffiliate(nil, 100)
			affiliateAddress := &types.Address{}
			copy(affiliateAddress[:], affiliateLog.Data[12:32])
			if err := consumer.affiliateService.Set(affiliateAddress, affiliate); err != nil {
				delivery.Return()
				log.Fatalf("Error registering affiliate: %#x", affiliateAddress[:])
			}
			log.Printf("Added affiliate: %#x", affiliateAddress[:])
		}
	}
	delivery.Ack()
}

func NewAffiliateBlockConsumer(affiliateSignupAddress *big.Int, lf ethereum.LogFilterer, affiliateService affiliates.AffiliateService) (channels.Consumer) {
	return &affiliateBlockConsumer{affiliateSignupAddress, affiliateService, lf}
}

func NewRPCAffiliateBlockConsumer(rpcURL string, affiliateSignupAddress string, redisClient *redis.Client) (channels.Consumer, error) {

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return NewAffiliateBlockConsumer(common.HexToAddress(affiliateSignupAddress).Big(), client, affiliates.NewRedisAffiliateService(redisClient)), nil
}
