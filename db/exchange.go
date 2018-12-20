package db

import (
	"github.com/notegio/openrelay/types"
	"github.com/jinzhu/gorm"
)

type Exchange struct {
	Address  *types.Address `gorm:"primary_key"`
	Network  int64          `gorm:"index"`
}

type ExchangeLookup struct {
	byAddressCache map[types.Address]int64
	byNetworkCache map[int64][]*types.Address
	db *gorm.DB
}

func (lookup *ExchangeLookup) GetExchangesByNetwork(network int64) ([]*types.Address, error) {
	if addresses, ok := lookup.byNetworkCache[network]; ok {
		return addresses, nil
	}
	addresses := []*types.Address{}
	exchanges := []Exchange{}
	if err := lookup.db.Model(&Exchange{}).Where("network = ?", network).Find(&exchanges).Error; err != nil {
		return nil, err
	}
	for _, exchange := range exchanges {
		addresses = append(addresses, exchange.Address)
	}
	lookup.byNetworkCache[network] = addresses
	return addresses, nil
}

func (lookup *ExchangeLookup) GetNetworkByExchange(address *types.Address) (int64, error) {
	if network, ok := lookup.byAddressCache[*address]; ok {
		return network, nil
	}
	exchange := &Exchange{}
	if err := lookup.db.Model(&Exchange{}).Where("address = ?", address).First(exchange).Error; err != nil {
		return 0, err
	}
	lookup.byAddressCache[*address] = exchange.Network
	return exchange.Network, nil
}

func (lookup *ExchangeLookup) ExchangeIsKnown(address *types.Address) (<-chan uint) {
	result := make(chan uint)
	go func(address *types.Address, result chan uint) {
		networkid, _ := lookup.GetNetworkByExchange(address)
		result <- uint(networkid)
	}(address, result)
	return result
}

func NewExchangeLookup(db *gorm.DB) (*ExchangeLookup) {
	return &ExchangeLookup{
		make(map[types.Address]int64),
		make(map[int64][]*types.Address),
		db,
	}
}
