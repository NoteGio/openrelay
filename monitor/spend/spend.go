package spend

import (
	"encoding/json"
	"math/big"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	coreTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/notegio/openrelay/funds/balance"
	"github.com/notegio/openrelay/channels"
	orCommon "github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/exchangecontract"
	"github.com/notegio/openrelay/monitor/blocks"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"log"
	"strings"
	"fmt"
)

type spendBlockConsumer struct {
	tokenProxyAddress *types.Address
	spendTopic        *big.Int
	feeTokenAddress    string  // Needed for the SpendRecord,
	logFilter          ethereum.LogFilterer
	publisher          channels.Publisher
	balanceChecker     balance.BalanceChecker
}

func (consumer *spendBlockConsumer) Consume(delivery channels.Delivery) {
	block := &blocks.MiniBlock{}
	err := json.Unmarshal([]byte(delivery.Payload()), block)
	if err != nil {
		log.Printf("Error parsing payload: %v\n", err.Error())
	}
	if coreTypes.BloomLookup(block.Bloom, consumer.spendTopic) {
		log.Printf("Block %#x bloom filter indicates spend event", block.Hash)
		query := ethereum.FilterQuery{
			FromBlock: block.Number,
			ToBlock: block.Number,
			Addresses: nil,
			Topics: [][]common.Hash{
				[]common.Hash{common.BigToHash(consumer.spendTopic)},
				nil,
				nil,
			},
		}
		logs, err := consumer.logFilter.FilterLogs(context.Background(), query)
		if err != nil {
			delivery.Return()
			log.Fatalf("Failed to filter logs on block %v - aborting: %v", block.Number, err.Error())
		}
		log.Printf("Found %v spend logs", len(logs))
		tradedTokens := make(map[string]struct{})
		for _, spendLog := range logs {
			// if len(spendLog.Topics) < 2 {
			// 	log.Printf("Unexpected log data. Skipping.")
			// 	continue
			// }
			senderAddress := &types.Address{}
			tokenAddress := &types.Address{}
			var tokenAssetData types.AssetData

			copy(tokenAddress[:], spendLog.Address[:])

			if len(spendLog.Topics) == 0  && len(spendLog.Data) >= 96{
				// CryptoKitties Style ERC721
				copy(senderAddress[:], spendLog.Data[12:32])
				tokenID := &types.Uint256{}
				copy(tokenID[:], spendLog.Data[len(spendLog.Data)-32:])
				tokenAssetData = orCommon.ToERC721AssetData(tokenAddress, tokenID)
			} else if len(spendLog.Topics) == 3 {
				// ERC20
				copy(senderAddress[:], spendLog.Topics[1][12:])
				tokenAssetData = orCommon.ToERC20AssetData(tokenAddress)
			} else if len(spendLog.Topics) == 4 {
				// ERC721
				copy(senderAddress[:], spendLog.Topics[1][12:])
				tokenID := &types.Uint256{}
				copy(tokenID[:], spendLog.Topics[3][:])
				tokenAssetData = orCommon.ToERC721AssetData(tokenAddress, tokenID)
			} else {
				log.Printf("Unexpected log data. Skipping.")
				continue
			}

			pairKey := fmt.Sprintf("%#x:%#x", senderAddress, tokenAddress)
			if _, ok := tradedTokens[pairKey]; ok {
				// If the same account sent the same token multiple times in a single
				// block, we already checked their balance as of the end of the block,
				// so we don't need to check it again.
				continue
			}
			tradedTokens[pairKey] = struct{}{}
			var balance *big.Int
			log.Printf("%#x - %#x - %#x", tokenAssetData[:], senderAddress[:], consumer.tokenProxyAddress[:])
			allowance, err := consumer.balanceChecker.GetAllowance(tokenAssetData, senderAddress, consumer.tokenProxyAddress)
			if err != nil && err.Error() == "VM Exception while processing transaction: revert" {
				// Some ERC721 tokens have the ERC20 signature, but implement ERC721
				// allowance / balance signatures
				log.Printf("ERC20 token with allowance mismatch: %#x", tokenAddress[:])
				tokenID := &types.Uint256{}
				copy(tokenID[:], spendLog.Data[:])
				tokenAssetData = orCommon.ToERC721AssetData(tokenAddress, tokenID)
				allowance, err = consumer.balanceChecker.GetAllowance(tokenAssetData, senderAddress, consumer.tokenProxyAddress)
			}
			if err != nil {
				if strings.Contains(err.Error(), "abi") || err.Error() == "no contract code at given address" || err.Error() == "VM Exception while processing transaction: revert" {
					log.Printf("balance checker gave error: %v -- using 0 balance", err.Error())
					allowance = big.NewInt(0)
				} else {
					delivery.Return()
					log.Fatalf("Failed to get allowance for '%v' - '%v': %v", tokenAddress, senderAddress, err.Error())
				}
			}
			if allowance.Cmp(big.NewInt(0)) == 0 {
				// If the allowance is 0, then the balance we want to report is 0, and
				// we don't need to get the actual balance.
				balance = allowance
			} else {
				balance, err = consumer.balanceChecker.GetBalance(tokenAssetData, senderAddress)
				if err != nil {
					if strings.Contains(err.Error(), "abi") || err.Error() == "no contract code at given address" || err.Error() == "VM Exception while processing transaction: revert" {
						log.Printf("balance checker gave error: %v -- using 0 balance", err.Error())
						balance = big.NewInt(0)
						} else {
							delivery.Return()
							log.Fatalf("Failed to get balance for '%v' - '%v': %v", tokenAddress, senderAddress, err.Error())
						}
					}
					if allowance.Cmp(balance) < 0 {
						// If allowance < balance, we should use that as our removal criteria
						balance = allowance
					}
			}

			sr := &db.SpendRecord{
				TokenAddress: strings.ToLower(spendLog.Address.String()),
				AssetData: hexutil.Encode(tokenAssetData[:]),
				SpenderAddress: hexutil.Encode(spendLog.Topics[1][12:]),
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
	} else {
		log.Printf("Block %v shows no spend events", block.Hash)
	}
	delivery.Ack()
}

func NewSpendBlockConsumer(tp *types.Address, feeToken string, lf ethereum.LogFilterer, publisher channels.Publisher, bc balance.BalanceChecker) (channels.Consumer) {
	spendTopic := &big.Int{}
	spendTopic.SetString("ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", 16)
	return &spendBlockConsumer{tp, spendTopic, feeToken, lf, publisher, bc}
}

func NewRPCSpendBlockConsumer(rpcURL string, exchangeAddress string, publisher channels.Publisher) (channels.Consumer, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	exchange, err := exchangecontract.NewExchangecontract(common.HexToAddress(exchangeAddress), client)
	if err != nil {
		log.Printf("Error intializing exchange contract '%v': '%v'", exchangeAddress, err.Error())
		return nil, err
	}
	feeTokenAssetData, err := exchange.ZRX_ASSET_DATA(nil)
	if err != nil {
		log.Printf("Error getting fee token address for exchange %v", exchangeAddress)
		return nil, err
	}
	feeTokenAsset := make(types.AssetData, len(feeTokenAssetData))
	copy(feeTokenAsset[:], feeTokenAssetData[:])
	feeTokenAddress := feeTokenAsset.Address()
	tokenProxyAddress, err := exchange.GetAssetProxy(nil, types.ERC20ProxyID)
	if err != nil {
		log.Printf("error getting tokenProxyAddress")
		return nil, err
	}
	tokenProxyAddressOr := &types.Address{}
	copy(tokenProxyAddressOr[:], tokenProxyAddress[:])
	balanceChecker, err := balance.NewRpcRoutingBalanceChecker(rpcURL)
	if err != nil {
		log.Printf("Error getting balance checker")
		return nil, err
	}
	return NewSpendBlockConsumer(tokenProxyAddressOr, feeTokenAddress.String(), client, publisher, balanceChecker), nil
}
