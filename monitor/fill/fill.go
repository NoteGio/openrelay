package fill

import (
	"encoding/json"
	"math/big"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	// "github.com/ethereum/go-ethereum/core/types"
	// "github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/fillbloom"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/monitor/blocks"
	"log"
	"fmt"
)

type fillBlockConsumer struct {
	exchangeAddress   *big.Int
	fillTopic         *big.Int // 0x0bcc4c97732e47d9946f229edb95f5b6323f601300e4690de719993f3c371129
	cancelTopic       *big.Int // 0xdc47b3613d9fe400085f6dbdc99453462279057e6207385042827ed6b1a62cf7
	logFilter         ethereum.LogFilterer
	publisher         channels.Publisher
	fillBloom         *fillbloom.FillBloom
}

func (consumer *fillBlockConsumer) Consume(delivery channels.Delivery) {
	block := &blocks.MiniBlock{}
	err := json.Unmarshal([]byte(delivery.Payload()), block)
	if err != nil {
		log.Printf("Error parsing payload: %v\n", err.Error())
	}
	if !consumer.fillBloom.Initialized {
		if err := consumer.fillBloom.Initialize(
			consumer.logFilter,
			block.Number.Int64(),
			[]common.Address{common.BigToAddress(consumer.exchangeAddress)},
		); err != nil {
			log.Fatalf("Failed to initialize bloom filter: %v", err.Error())
		}
	}
	if (block.Bloom.Test(consumer.fillTopic) || block.Bloom.Test(consumer.cancelTopic)) && block.Bloom.Test(consumer.exchangeAddress) {
		log.Printf("Block %#x bloom filter indicates fill event for %#x", block.Hash, consumer.exchangeAddress)
		query := ethereum.FilterQuery{
			FromBlock: block.Number,
			ToBlock: block.Number,
			Addresses: []common.Address{common.BigToAddress(consumer.exchangeAddress)},
			Topics: [][]common.Hash{
				[]common.Hash{common.BigToHash(consumer.fillTopic), common.BigToHash(consumer.cancelTopic)},
				nil,
				nil,
			},
		}
		logs, err := consumer.logFilter.FilterLogs(context.Background(), query)
		if err != nil {
			delivery.Return()
			log.Fatalf("Failed to filter logs on block %v - aborting: %v", block.Number, err.Error())
		}
		log.Printf("Found %v fill logs", len(logs))
		for _, fillLog := range logs {
			var fr *db.FillRecord
			if new(big.Int).SetBytes(fillLog.Topics[0][:]).Cmp(consumer.fillTopic) == 0 {
				takerTokenFilled := big.NewInt(0)
				takerTokenFilled.SetBytes(fillLog.Data[32*3:32*4])
				orderHash := fillLog.Topics[3][:]
				fr = &db.FillRecord{
					OrderHash: fmt.Sprintf("%#x", orderHash),
					FilledTakerAssetAmount: takerTokenFilled.Text(10),
					Cancel: false,
				}
				consumer.fillBloom.Add(orderHash)
			} else {
				orderHash := fillLog.Topics[3][:]
				fr = &db.FillRecord{
					OrderHash: fmt.Sprintf("%#x", orderHash),
					FilledTakerAssetAmount: "0",
					Cancel: true,
				}
				consumer.fillBloom.Add(orderHash)
			}
			msg, err := json.Marshal(fr)
			if err != nil {
				delivery.Return()
				log.Fatalf("Failed to encode FillRecord on block %v", block.Number)
			}
			consumer.publisher.Publish(string(msg))
		}
		if err := consumer.fillBloom.Save(); err != nil {
			log.Printf("error saving bloom filter: %v", err.Error())
		}
	} else {
		log.Printf("Block %#x shows no fill events", block.Hash)
	}
	delivery.Ack()
}

func NewFillBlockConsumer(exchangeAddress *big.Int, lf ethereum.LogFilterer, publisher channels.Publisher, fb *fillbloom.FillBloom) (channels.Consumer) {
	fillTopic := &big.Int{}
	fillTopic.SetString("0bcc4c97732e47d9946f229edb95f5b6323f601300e4690de719993f3c371129", 16)
	cancelTopic := &big.Int{}
	cancelTopic.SetString("dc47b3613d9fe400085f6dbdc99453462279057e6207385042827ed6b1a62cf7", 16)
	return &fillBlockConsumer{exchangeAddress, fillTopic, cancelTopic, lf, publisher, fb}
}

func NewRPCFillBlockConsumer(rpcURL string, exchangeAddress string, publisher channels.Publisher, fb *fillbloom.FillBloom) (channels.Consumer, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return NewFillBlockConsumer(common.HexToAddress(exchangeAddress).Big(), client, publisher, fb), nil
}
