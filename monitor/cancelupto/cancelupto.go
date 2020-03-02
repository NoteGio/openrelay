package cancelupto

import (
	"encoding/json"
	"math/big"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	coreTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/monitor/blocks"
	"log"
)

type cancelBlockConsumer struct {
	exchangeAddress   *big.Int
	cancelUpToTopic   *big.Int // 0x82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0.
	logFilter         ethereum.LogFilterer
	publisher         channels.Publisher
}

func (consumer *cancelBlockConsumer) Consume(delivery channels.Delivery) {
	block := &blocks.MiniBlock{}
	err := json.Unmarshal([]byte(delivery.Payload()), block)
	if err != nil {
		log.Printf("Error parsing payload: %v\n", err.Error())
	}
	if (coreTypes.BloomLookup(block.Bloom, consumer.cancelUpToTopic)) && coreTypes.BloomLookup(block.Bloom, consumer.exchangeAddress) {
		log.Printf("Block %#x bloom filter indicates cancelUpTo event for %#x", block.Hash, consumer.exchangeAddress)
		query := ethereum.FilterQuery{
			FromBlock: block.Number,
			ToBlock: block.Number,
			Addresses: []common.Address{common.BigToAddress(consumer.exchangeAddress)},
			Topics: [][]common.Hash{
				[]common.Hash{common.BigToHash(consumer.cancelUpToTopic)},
				nil,
				nil,
			},
		}
		logs, err := consumer.logFilter.FilterLogs(context.Background(), query)
		if err != nil {
			delivery.Return()
			log.Fatalf("Failed to filter logs on block %v - aborting: %v", block.Number, err.Error())
		}
		log.Printf("Found %v cancellation logs", len(logs))
		for _, cancelLog := range logs {
			if len(cancelLog.Topics) < 3 || len(cancelLog.Data) != 32 {
				log.Printf("Unexpected log data. Skipping.")
				continue
			}
			cancellation := &db.Cancellation{&types.Address{}, &types.Address{}, &types.Uint256{}}
			copy(cancellation.Maker[:], cancelLog.Topics[1][12:])
			copy(cancellation.Sender[:], cancelLog.Topics[2][12:])
			copy(cancellation.Epoch[:], cancelLog.Data[:])
			msg, err := json.Marshal(cancellation)
			if err != nil {
				delivery.Return()
				log.Fatalf("Failed to encode Cancellation on block %v: %v", block.Number, err.Error())
			}
			consumer.publisher.Publish(string(msg))
		}
	} else {
		log.Printf("Block %#x shows no cancelUpTo events", block.Hash)
	}
	delivery.Ack()
}

func NewCancelUpToBlockConsumer(exchangeAddress *big.Int, lf ethereum.LogFilterer, publisher channels.Publisher) (channels.Consumer) {
	cancelUpToTopic := &big.Int{}
	cancelUpToTopic.SetString("82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0", 16)
	return &cancelBlockConsumer{exchangeAddress, cancelUpToTopic, lf, publisher}
}

func NewRPCCancelUpToBlockConsumer(rpcURL string, exchangeAddress string, publisher channels.Publisher) (channels.Consumer, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	exchangeAddressBig, _ := big.NewInt(0).SetString(exchangeAddress, 0)
	return NewCancelUpToBlockConsumer(exchangeAddressBig, client, publisher), nil
}
