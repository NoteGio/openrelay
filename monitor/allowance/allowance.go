package allowance

import (
	"encoding/json"
	"math/big"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	// "github.com/ethereum/go-ethereum/core/types"
	// "github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/exchangecontract"
	"github.com/notegio/openrelay/monitor/blocks"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"log"
)

type allowanceBlockConsumer struct {
	tokenProxyAddress *big.Int
	approvalTopic     *big.Int
	feeTokenAddress   string  // Needed for the SpendRecord,
	logFilter         ethereum.LogFilterer
	publisher         channels.Publisher
}

func (consumer *allowanceBlockConsumer) Consume(delivery channels.Delivery) {
	log.Printf("Start")
	block := &blocks.MiniBlock{}
	err := json.Unmarshal([]byte(delivery.Payload()), block)
	if err != nil {
		log.Printf("Error parsing payload: %v\n", err.Error())
	}
	if block.Bloom.Test(consumer.approvalTopic) && block.Bloom.Test(consumer.approvalTopic) {
		query := ethereum.FilterQuery{
			FromBlock: block.Number,
			ToBlock: block.Number,
			Addresses: nil,
			Topics: [][]common.Hash{
				[]common.Hash{common.BigToHash(consumer.approvalTopic)},
				[]common.Hash{},
				[]common.Hash{common.BigToHash(consumer.tokenProxyAddress)},
			},
		}
		logs, err := consumer.logFilter.FilterLogs(context.Background(), query)
		if err != nil {
			delivery.Return()
			log.Fatalf("Failed to filter logs on block %v - aborting", block.Number)
		}
		for _, approvalLog := range logs {
			balance := big.NewInt(0)
			balance.SetBytes(approvalLog.Data)
			sr := &db.SpendRecord{
				TokenAddress: approvalLog.Address.String(),
				SpenderAddress: hexutil.Encode(approvalLog.Topics[1][12:]),
				ZrxToken: consumer.feeTokenAddress,
				Balance: balance.String(),
			}
			msg, err := json.Marshal(sr)
			if err != nil {
				delivery.Return()
				log.Fatalf("Failed to encode SpendRecord on block %v", block.Number)
			}
			consumer.publisher.Publish(string(msg))
		}
		delivery.Ack()
	}
}

/* TODO:
* Tests
* Command
* Get feeToken from exchange address
* Get tokenProxy from exchange address
*/
func NewAllowanceBlockConsumer(tp *big.Int, feeToken string, lf ethereum.LogFilterer, publisher channels.Publisher) (channels.Consumer) {
	approvalTopic := &big.Int{}
	approvalTopic.SetString("8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", 16)
	return &allowanceBlockConsumer{tp, approvalTopic, feeToken, lf, publisher}
}

func NewRPCAllowanceBlockConsumer(rpcURL string, exchangeAddress string, publisher channels.Publisher) (channels.Consumer, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	exchange, err := exchangecontract.NewExchange(common.StringToAddress(exchangeAddress), client)
	if err != nil {
		log.Printf("Error intializing exchange contract '%v': '%v'", exchangeAddress, err.Error())
		return nil, err
	}
	feeTokenAddress, err := exchange.ZRX_TOKEN_CONTRACT(nil)
	if err != nil {
		log.Printf("error getting feeTokenAddress")
		return nil, err
	}
	tokenProxyAddress, err := exchange.TOKEN_TRANSFER_PROXY_CONTRACT(nil)
	if err != nil {
		log.Printf("error getting tokenProxyAddress")
		return nil, err
	}
	return NewAllowanceBlockConsumer(tokenProxyAddress.Big(), feeTokenAddress.Hex(), client, publisher), nil
}
