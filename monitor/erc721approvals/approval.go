package erc721approvals

import (
	"encoding/json"
	"math/big"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	coreTypes "github.com/ethereum/go-ethereum/core/types"
	// "github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/types"
	orCommon "github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/exchangecontract"
	"github.com/notegio/openrelay/monitor/blocks"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"log"
	"strings"
	"fmt"
)

type approvalBlockConsumer struct {
	tokenProxyAddress *big.Int
	approvalTopic     *big.Int
	approveAllTopic   *big.Int
	logFilter         ethereum.LogFilterer
	publisher         channels.Publisher
}

func (consumer *approvalBlockConsumer) Consume(delivery channels.Delivery) {
	block := &blocks.MiniBlock{}
	err := json.Unmarshal([]byte(delivery.Payload()), block)
	if err != nil {
		log.Printf("Error parsing payload: %v\n", err.Error())
	}
	if coreTypes.BloomLookup(block.Bloom, consumer.approvalTopic) || (coreTypes.BloomLookup(block.Bloom, consumer.approveAllTopic) && coreTypes.BloomLookup(block.Bloom, consumer.tokenProxyAddress)){
		// TODO: This test is errantly failing. Not sure why.
		log.Printf("Block %#x bloom filter indicates approval event for %#x", block.Hash, consumer.tokenProxyAddress)
		query := ethereum.FilterQuery{
			FromBlock: block.Number,
			ToBlock: block.Number,
			Addresses: nil,
			Topics: [][]common.Hash{
				[]common.Hash{common.BigToHash(consumer.approvalTopic), common.BigToHash(consumer.approveAllTopic)},
				nil,
				nil,
				nil, // Having 3 nils excludes ERC20 token approval logs
			},
		}
		logs, err := consumer.logFilter.FilterLogs(context.Background(), query)
		if err != nil {
			delivery.Return()
			log.Fatalf("Failed to filter logs on block %v - aborting: %v", block.Number, err.Error())
		}
		log.Printf("Found %v approval logs", len(logs))
		for _, approvalLog := range logs {
			if new(big.Int).SetBytes(approvalLog.Topics[0][:]).Cmp(consumer.approvalTopic) == 0 {
				if topicCount := len(approvalLog.Topics); topicCount != 4 {
					// Not enough topics, probably an ERC20 event
					log.Printf("Expected 4 topics, got %v - %v", topicCount, approvalLog.Address.String())
					continue
				}
				// Approve
				if new(big.Int).SetBytes(approvalLog.Topics[2][:]).Cmp(consumer.tokenProxyAddress) == 0 {
					// They've just set the tokenProxy as the approved sender. This can't
					// make an order invalid, so we don't need to send anything.
					continue
				} else {
					// The new approved sender is not the tokenProxy. We're going to
					// assume this means the tokenProxy is not approved at all. It's
					// technically possible that the tokenProxy is an approvedForAll
					// operator, but having to check that for every ERC721 approval ever
					// could get pretty cumbersome, so we're going to the less precise
					// more efficient approach of removing all orders where the target
					// asset is Approved for someone other than the token proxy.

					// This will prune all invalidated orders, but may also prune some
					// valid ones.
					sr := &db.SpendRecord{
						TokenAddress: strings.ToLower(approvalLog.Address.String()),
						SpenderAddress: hexutil.Encode(approvalLog.Topics[1][12:]),
						AssetData: fmt.Sprintf("%#x", orCommon.ToERC721AssetData(orCommon.BytesToOrAddress(approvalLog.Address), orCommon.BytesToUint256(approvalLog.Topics[3]))),
						Balance: "0",
					}
					msg, err := json.Marshal(sr)
					if err != nil {
						delivery.Return()
						log.Fatalf("Failed to encode SpendRecord on block %v", block.Number)
					}
					consumer.publisher.Publish(string(msg))
					continue
				}

			} else {
				if topicCount := len(approvalLog.Topics); topicCount != 3 {
					// Not enough topics, probably an ERC20 event
					log.Printf("Expected 3 topics, got %v - %v", topicCount, approvalLog.Address.String())
					continue
				}
				// ApproveForAll
				if new(big.Int).SetBytes(approvalLog.Topics[2][:]).Cmp(consumer.tokenProxyAddress) == 0 {
					if new(big.Int).SetBytes(approvalLog.Data[:]).Cmp(big.NewInt(0)) == 0 {
						// If they've revoked the tokenProxy as an operator, we're going to
						// remove all orders for tokens of this asset type. We don't have a
						// mechanism to easily check whether specific assets have
						// asset-specific approval, so we're going to remove all of them.
						// This may lead to some orders improperly pruned, but won't leave
						// any invalid orders on the orderbook.
						sr := &db.SpendRecord{
							AssetData: "",
							TokenAddress: strings.ToLower(approvalLog.Address.String()),
							SpenderAddress: hexutil.Encode(approvalLog.Topics[1][12:]),
							Balance: "0",
						}
						msg, err := json.Marshal(sr)
						if err != nil {
							delivery.Return()
							log.Fatalf("Failed to encode SpendRecord on block %v", block.Number)
						}
						consumer.publisher.Publish(string(msg))
					}
				}
			}
		}
	} else {
		log.Printf("Block %v shows no approval events", block.Hash)
	}
	delivery.Ack()
}

func NewAllowanceBlockConsumer(tp *big.Int, lf ethereum.LogFilterer, publisher channels.Publisher) (channels.Consumer) {
	approvalTopic := &big.Int{}
	approveAllTopic := &big.Int{}
	approvalTopic.SetString("8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", 16)
	approveAllTopic.SetString("17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31", 16)
	return &approvalBlockConsumer{tp, approvalTopic, approveAllTopic, lf, publisher}
}

func NewRPCAllowanceBlockConsumer(rpcURL string, exchangeAddress string, publisher channels.Publisher) (channels.Consumer, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	exchange, err := exchangecontract.NewExchangecontract(common.HexToAddress(exchangeAddress), client)
	if err != nil {
		log.Printf("Error intializing exchange contract '%v': '%v'", exchangeAddress, err.Error())
		return nil, err
	}
	tokenProxyAddress, err := exchange.GetAssetProxy(nil, types.ERC721ProxyID)
	if err != nil {
		log.Printf("error getting tokenProxyAddress")
		return nil, err
	}
	log.Printf("TP: %#x - %v", tokenProxyAddress[:], exchangeAddress)
	return NewAllowanceBlockConsumer(big.NewInt(0).SetBytes(tokenProxyAddress[:]), client, publisher), nil
}
