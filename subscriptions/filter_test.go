package subscriptions_test

import (
	"errors"
	"testing"
	"github.com/notegio/openrelay/common"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/subscriptions"
	"github.com/notegio/openrelay/types"
	"fmt"
)

type MockExchangeLookup struct {
	mapping map[int64][]*types.Address
}

func (lookup *MockExchangeLookup) GetExchangesByNetwork(network int64) ([]*types.Address, error) {
	if exchanges, ok := lookup.mapping[network]; ok {
		return exchanges, nil
	}
	return []*types.Address{}, errors.New("Bad network id")
}


type sample struct {
	query string
	match bool
}

func TestFilterFunctions(t *testing.T) {
	torder := &types.Order{}
	torder.Initialize()
	order := &dbModule.Order{Order: *torder}
	order.Maker, _ = common.HexToAddress("0x1111111111111111111111111111111111111111")
	order.Taker, _ = common.HexToAddress("0x2222222222222222222222222222222222222222")
	order.MakerAssetData, _ = common.HexToAssetData("0xf47261b00000000000000000000000003333333333333333333333333333333333333333")
	order.MakerAssetAddress = order.MakerAssetData.Address()
	order.TakerAssetData, _ = common.HexToAssetData("0x025717920000000000000000000000004444444444444444444444444444444444444444000000000000000000000000000000000000000000000000000000000000000F")
	order.TakerAssetAddress = order.TakerAssetData.Address()
	order.TakerAssetAmount = common.Int64ToUint256(1)
	order.ExchangeAddress, _ = common.HexToAddress("0x5555555555555555555555555555555555555555")
	order.Populate()
	mapping := make(map[int64][]*types.Address)
	mapping[1] = []*types.Address{order.ExchangeAddress}
	sampleExchange, _ := common.HexToAddress("0x6666666666666666666666666666666666666666")
	mapping[2] = []*types.Address{sampleExchange}
	lookup := &MockExchangeLookup{mapping}
	samples := []sample{
		sample{"makerAddress=0x1111111111111111111111111111111111111111&takerAddress=0x2222222222222222222222222222222222222222&makerAssetProxyId=0xf47261b0&takerAssetProxyId=0x02571792", true},
		sample{"makerAssetAddress=0x3333333333333333333333333333333333333333", true},
		sample{"takerAssetAddress=0x3333333333333333333333333333333333333333", false},
		sample{"networkId=1", true},
		sample{"networkId=2", false},
		sample{"traderAddress=0x1111111111111111111111111111111111111111", true},
		sample{"traderAddress=0x3333333333333333333333333333333333333333", false},
		sample{"assetAddress=0x3333333333333333333333333333333333333333", true},
		sample{"assetAddress=0x1111111111111111111111111111111111111111", false},
		sample{"_poolName=nil", true},
		sample{"_poolName=foo", false},
		sample{fmt.Sprintf("_poolId=%#x", dbModule.DefaultSha3()), true},
		sample{"_poolId=0x0000000000000000000000000000000000000000000000000000000000000000", false},
		sample{"exchangeContractAddress=0x5555555555555555555555555555555555555555", true},
		sample{"exchangeContractAddress=0x6666666666666666666666666666666666666666", false},
		sample{"feeRecipient=0x6666666666666666666666666666666666666666&senderAddress=0x0000000000000000000000000000000000000000", false},
		sample{"senderAddress=0x0000000000000000000000000000000000000000", true},
		sample{"makerAssetData=0xf47261b00000000000000000000000003333333333333333333333333333333333333333", true},
		sample{"takerAssetData=0xf47261b00000000000000000000000003333333333333333333333333333333333333333", false},
		sample{"traderAssetData=0xf47261b00000000000000000000000003333333333333333333333333333333333333333", true},
		sample{"", true},
	}
	for _, sample := range samples {
		filter, err := subscriptions.FilterFromQueryString(sample.query)
		if err != nil {
			t.Errorf(err.Error())
			continue
		}
		filterFn, err := filter.GetFilter(lookup)
		if err != nil {
			t.Errorf(err.Error())
			continue
		}
		if filterFn(order) != sample.match{
			t.Errorf("Expected filter %v to match order (%v)", sample.query, sample.match)
		}
	}

}
