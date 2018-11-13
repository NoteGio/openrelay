package fill

import (
	"encoding/json"
	"math/big"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	coreTypes "github.com/ethereum/go-ethereum/core/types"
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
	fillTopic         *big.Int // 0x0d0b9391970d9a25552f37d436d2aae2925e2bfe1b2a923754bada030c498cb3
	cancelTopic       *big.Int // 0x67d66f160bc93d925d05dae1794c90d2d6d6688b29b84ff069398a9b04587131
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
	if (coreTypes.BloomLookup(block.Bloom, consumer.fillTopic) || coreTypes.BloomLookup(block.Bloom, consumer.cancelTopic)) && coreTypes.BloomLookup(block.Bloom, consumer.exchangeAddress) {
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
			if len(fillLog.Data) < 256 {
				log.Printf("Unexpected log data. Skipping.")
				continue
			}
			var fr *db.FillRecord
			if new(big.Int).SetBytes(fillLog.Topics[0][:]).Cmp(consumer.fillTopic) == 0 {
				takerTokenFilled := big.NewInt(0)
				takerTokenFilled.SetBytes(fillLog.Data[32*4:32*5])
				orderHash := fillLog.Data[32*7:32*8]
				fr = &db.FillRecord{
					OrderHash: fmt.Sprintf("%#x", orderHash),
					FilledTakerTokenAmount: takerTokenFilled.Text(10),
					CancelledTakerTokenAmount: "0",
				}
				consumer.fillBloom.Add(orderHash)
			} else {
				takerTokenCancelled := big.NewInt(0)
				takerTokenCancelled.SetBytes(fillLog.Data[32*3:32*4])
				orderHash := fillLog.Data[32*4:32*5]
				fr = &db.FillRecord{
					OrderHash: fmt.Sprintf("%#x", orderHash),
					FilledTakerTokenAmount: "0",
					CancelledTakerTokenAmount: takerTokenCancelled.Text(10),
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
	fillTopic.SetString("0d0b9391970d9a25552f37d436d2aae2925e2bfe1b2a923754bada030c498cb3", 16)
	cancelTopic := &big.Int{}
	cancelTopic.SetString("67d66f160bc93d925d05dae1794c90d2d6d6688b29b84ff069398a9b04587131", 16)
	return &fillBlockConsumer{exchangeAddress, fillTopic, cancelTopic, lf, publisher, fb}
}

func NewRPCFillBlockConsumer(rpcURL string, exchangeAddress string, publisher channels.Publisher, fb *fillbloom.FillBloom) (channels.Consumer, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return NewFillBlockConsumer(common.HexToAddress(exchangeAddress).Big(), client, publisher, fb), nil
}
