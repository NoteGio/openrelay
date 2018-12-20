package db_test

import (
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/types"
	// "math/big"
	"bytes"
	"testing"
	// "log"
)

func TestExchangesByNetwork(t *testing.T) {
	db, err := getDb()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Exchange{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	address := &types.Address{}
	tx.Model(&dbModule.Exchange{}).Create(&dbModule.Exchange{address, 1})
	lookup := dbModule.NewExchangeLookup(tx)
	exchanges, err := lookup.GetExchangesByNetwork(1)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(exchanges) != 1 {
		t.Fatalf("Unexpected exchange count")
	}
	if !bytes.Equal(exchanges[0][:], address[:]) {
		t.Errorf("Unexpected exchange value: %v", exchanges[0])
	}
}

func TestNetworkByExchange(t *testing.T) {
	db, err := getDb()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Exchange{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	address := &types.Address{}
	tx.Model(&dbModule.Exchange{}).Create(&dbModule.Exchange{address, 1})
	lookup := dbModule.NewExchangeLookup(tx)
	networkID, err := lookup.GetNetworkByExchange(address)
	if err != nil {
		t.Errorf(err.Error())
	}
	if networkID != 1 {
		t.Errorf("Unexpected exchange id")
	}
}
func TestExchangeIsKnown(t *testing.T) {
	db, err := getDb()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Exchange{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	address := &types.Address{}
	lookup := dbModule.NewExchangeLookup(tx)
	if (<-lookup.ExchangeIsKnown(address) != 0) {
		t.Errorf("Expected exchange to be unknown")
	}
	tx.Model(&dbModule.Exchange{}).Create(&dbModule.Exchange{address, 1})
	if (<-lookup.ExchangeIsKnown(address) == 0) {
		t.Errorf("Expected exchange to be known")
	}
}
