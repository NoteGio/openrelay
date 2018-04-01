package blocks

import (
	"context"
	"math/big"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"reflect"
)

type MockHeaderGetter struct {
	// We're using a map for the headers instead of a list because we need to be
	// able to represent chain reorgs where the header at a given number can
	// change.
	headers map[int64]*types.Header
}


func (hg *MockHeaderGetter) HeaderByNumber(ctx context.Context, bigIdx *big.Int) (*types.Header, error){
	index := bigIdx.Int64()
	if index >= int64(len(hg.headers)) {
		return nil, ethereum.NotFound
	}
	return hg.headers[index], nil

}
func (hg *MockHeaderGetter) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error){
	for _, header := range hg.headers {
		if reflect.DeepEqual(header.Hash(), hash) {
			return header, nil
		}
	}
	return nil, ethereum.NotFound
}

func (hg *MockHeaderGetter) AddHeader(header *types.Header) {
	hg.headers[header.Number.Int64()] = header
}


func NewMockHeaderGetter(headers []*types.Header) *MockHeaderGetter {
	headerMap := make(map[int64]*types.Header)
	for _, header := range headers {
		headerMap[header.Number.Int64()] = header
	}
	return &MockHeaderGetter{headerMap}
}
