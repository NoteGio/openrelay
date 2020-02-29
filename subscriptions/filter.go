package subscriptions

import (
	"encoding/hex"
	"github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/common"
	"golang.org/x/crypto/sha3"
	"net/url"
	"bytes"
	"strings"
	"strconv"
)

type OrderFilter struct {
	NetworkID           int64  `json:"networkId"`
	MakerAssetProxyID   string `json:"makerAssetProxyId"`
	TakerAssetProxyID   string `json:"takerAssetProxyId"`
	MakerAssetAddress   string `json:"makerAssetAddress"`
	TakerAssetAddress   string `json:"takerAssetAddress"`
	AssetAddress        string `json:"assetAddress"`
	ExchangeAddress     string `json:"exchangeContractAddress"`
	SenderAddress       string `json:"senderAddress"`
	MakerAddress        string `json:"makerAddress"`
	TakerAddress        string `json:"takerAddress"`
	TraderAddress       string `json:"traderAddress"`
	FeeRecipientAddress string `json:"feeRecipient"`
	MakerAssetData      string `json:"makerAssetData"`
	TakerAssetData      string `json:"takerAssetData"`
	TraderAssetData     string `json:"traderAssetData"`
	PoolID              string `json:"_poolId"`
	PoolName            string `json:"_poolName"`
	TakerFee            string `json:"_takerFee"`
}

type ExchangeLookup interface {
	GetExchangesByNetwork(network int64) ([]*types.Address, error)
}

// GetFilter returns a function that can be quickly used to evaluate whether a
// given order matches the criteria of this filter. This is sort of a
// just-in-time compiler, in that it produces an efficient fuction for testing
// an order. While this function is fairly complicated, the function it returns
// will test only the attributes that the filter is interested in, and does not
// re-evaluate filter parameters when checking each order.
func (ofilter *OrderFilter) GetFilter(lookup ExchangeLookup) (func(*db.Order) (bool), error) {
	predicates := []func(*db.Order)(bool){}
	// TODO: Once we have some data on the combinations of attributes people
	// filter on, we can potentially order these to maximize filtering efficiency.
	if ofilter.MakerAssetProxyID != "" {
		data := [4]byte{}
		raw, err := hex.DecodeString(strings.TrimPrefix(ofilter.MakerAssetProxyID, "0x"))
		if err != nil { return nil, err }
		copy(data[:], raw[:])
		predicates = append(predicates, func(order *db.Order) (bool) { return order.MakerAssetData.IsType(data) } )
	}
	if ofilter.TakerAssetProxyID != "" {
		data := [4]byte{}
		raw, err := hex.DecodeString(strings.TrimPrefix(ofilter.TakerAssetProxyID, "0x"))
		if err != nil { return nil, err }
		copy(data[:], raw[:])
		predicates = append(predicates, func(order *db.Order) (bool) { return order.TakerAssetData.IsType(data) } )
	}
	if ofilter.MakerAssetAddress != "" {
		address, err  := common.HexToAddress(ofilter.MakerAssetAddress)
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.MakerAssetAddress[:], address[:]) } )
	}
	if ofilter.TakerAssetAddress != "" {
		address, err := common.HexToAddress(ofilter.TakerAssetAddress)
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.TakerAssetAddress[:], address[:]) } )
	}
	if ofilter.AssetAddress != "" {
		address, err := common.HexToAddress(ofilter.AssetAddress)
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.TakerAssetAddress[:], address[:]) || bytes.Equal(order.MakerAssetAddress[:], address[:]) } )
	}
	if ofilter.ExchangeAddress != "" {
		address, err := common.HexToAddress(ofilter.ExchangeAddress)
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.ExchangeAddress[:], address[:]) } )
	}
	if ofilter.FeeRecipientAddress != "" {
		address, err := common.HexToAddress(ofilter.FeeRecipientAddress)
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.FeeRecipient[:], address[:]) } )
	}
	if ofilter.SenderAddress != "" {
		address, err := common.HexToAddress(ofilter.SenderAddress)
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.SenderAddress[:], address[:]) } )
	}
	if ofilter.MakerAddress != "" {
		address, err := common.HexToAddress(ofilter.MakerAddress)
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.Maker[:], address[:]) } )
	}
	if ofilter.TakerAddress != "" {
		address, err := common.HexToAddress(ofilter.TakerAddress)
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.Taker[:], address[:]) } )
	}
	if ofilter.TraderAddress != "" {
		address, err := common.HexToAddress(ofilter.TraderAddress)
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.Maker[:], address[:]) || bytes.Equal(order.Taker[:], address[:]) } )
	}
	if ofilter.MakerAssetData != "" {
		assetData, err := common.HexToAssetData(ofilter.MakerAssetData)
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.MakerAssetData, assetData) } )
	}
	if ofilter.TakerAssetData != "" {
		assetData, err := common.HexToAssetData(ofilter.TakerAssetData)
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.TakerAssetData, assetData) } )
	}
	if ofilter.TraderAssetData != "" {
		assetData, err := common.HexToAssetData(ofilter.TraderAssetData)
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.TakerAssetData, assetData) || bytes.Equal(order.MakerAssetData, assetData) } )
	}
	if ofilter.PoolID != "" {
		poolID, err := hex.DecodeString(strings.TrimPrefix(ofilter.PoolID, "0x"))
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.PoolID, poolID) } )
	}
	if ofilter.PoolName != "" {
		dataSha := sha3.NewLegacyKeccak256()
		dataSha.Write([]byte(ofilter.PoolName))
		poolID := dataSha.Sum(nil)
		predicates = append(predicates, func(order *db.Order) (bool) { return bytes.Equal(order.PoolID, poolID) } )
	}
	if ofilter.NetworkID != 0 {
		exchangeAddresses, err := lookup.GetExchangesByNetwork(ofilter.NetworkID)
		if err != nil { return nil, err }
		predicates = append(predicates, func(order *db.Order) (bool) {
			for _, address := range exchangeAddresses {
				if bytes.Equal(order.ExchangeAddress[:], address[:]) {
					return true
				}
			}
			return false
		})
	}
	return func(order *db.Order) (bool) {
		for _, predicate := range predicates {
			if !predicate(order) {
				return false
			}
		}
		return true
	}, nil
}

func FilterFromQueryString(queryString string) (*OrderFilter, error) {
	query, err := url.ParseQuery(queryString)
	if err != nil { return nil, err}
	ofilter := &OrderFilter{}
	if query.Get("networkId") != "" {
		networkID, err := strconv.Atoi(query.Get("networkId"))
		if err != nil { return nil, err}
		ofilter.NetworkID = int64(networkID)
	} else {
		ofilter.NetworkID = 0
	}
	ofilter.MakerAssetProxyID = query.Get("makerAssetProxyId")
	ofilter.TakerAssetProxyID = query.Get("takerAssetProxyId")
	ofilter.MakerAssetAddress = query.Get("makerAssetAddress")
	ofilter.TakerAssetAddress = query.Get("takerAssetAddress")
	ofilter.AssetAddress = query.Get("assetAddress")
	ofilter.ExchangeAddress = query.Get("exchangeContractAddress")
	ofilter.SenderAddress = query.Get("senderAddress")
	ofilter.MakerAddress = query.Get("makerAddress")
	ofilter.TakerAddress = query.Get("takerAddress")
	ofilter.TraderAddress = query.Get("traderAddress")
	ofilter.FeeRecipientAddress = query.Get("feeRecipient")
	ofilter.MakerAssetData = query.Get("makerAssetData")
	ofilter.TakerAssetData = query.Get("takerAssetData")
	ofilter.TraderAssetData = query.Get("traderAssetData")
	ofilter.PoolID = query.Get("_poolId")
	ofilter.PoolName = query.Get("_poolName")
	ofilter.TakerFee = query.Get("_takerFee")
	return ofilter, nil
}
