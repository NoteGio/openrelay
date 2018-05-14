package mock

import (
	"context"
	"bytes"
	"errors"
	"math/big"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum"
	"log"
)

func GenerateBlockHeader(parentHash common.Hash, number int64, topics []common.Hash) (*types.Header){
	bloom := types.Bloom{}
	for _, topic := range topics {
		bloom.Add(new(big.Int).SetBytes(topic[:]))
	}
	return &types.Header{
		ParentHash: parentHash,
		UncleHash: common.Hash{},
		Coinbase: common.Address{},
		Root: common.Hash{},
		TxHash: common.Hash{},
		ReceiptHash: common.Hash{},
		Bloom: bloom,
		Difficulty: new(big.Int),
		Number: big.NewInt(number),
		GasLimit: new(big.Int),
		GasUsed: new(big.Int),
		Time: big.NewInt(number),
		Extra: []byte{},
		MixDigest: common.Hash{},
		Nonce: types.BlockNonce{},
	}
}

func GenerateHeaderChain(n int64) []*types.Header {
	return GenerateChainSplit(0, n, common.Hash{}, []byte{})
}

func GenerateChainSplit(start, n int64, parentHash common.Hash, extra []byte) []*types.Header {
	header := GenerateBlockHeader(parentHash, start, []common.Hash{})
	header.Extra = extra
	headers := []*types.Header{header}
	for i := start + 1; i < start + n; i++ {
		header = GenerateBlockHeader(header.Hash(), int64(i), []common.Hash{})
		headers = append(headers, header)
	}
	return headers
}

type MockLogFilterer struct {
	logs []types.Log
}

func (lf *MockLogFilterer) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	results := []types.Log{}
	LOG_LOOP:
	for _, testLog := range lf.logs {
		match := len(q.Addresses) == 0
		for _, address := range q.Addresses {
			if bytes.Equal(address[:], testLog.Address[:]) {
				match = true
				break
			}
		}
		if !match {
			continue
		}
		for i, topics := range q.Topics {
			match = len(topics) == 0 // If there are no topics listed, everything matches
			for _, topic := range topics {
				if bytes.Equal(topic[:], testLog.Topics[i][:]) {
					match = true
					break
				}
			}
			if !match {
				log.Printf("Topic %v does not match", i)
				continue LOG_LOOP
			}
		}
		results = append(results, testLog)
	}
	return results, nil
}

func (lf *MockLogFilterer) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return nil, errors.New("Not Implemented")
}

func NewMockLogFilterer(logs []types.Log) (ethereum.LogFilterer) {
	return &MockLogFilterer{logs}
}
