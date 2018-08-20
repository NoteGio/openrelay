package db

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/types"
	"math/big"
	"strings"
	"bytes"
	// "log"
)

type FillRecord struct {
	OrderHash                 string `json:"orderHash"`
	FilledTakerAssetAmount    string `json:"filledTakerAssetAmount"`
	Cancel                    bool   `json:"cancel"`
}

type Indexer struct {
	db     *gorm.DB
	status int64
}

// Index takes an order and saves it to the database
func (indexer *Indexer) Index(order *types.Order) error {
	dbOrder := Order{}
	dbOrder.Order = *order
	return dbOrder.Save(indexer.db, indexer.status).Error
}

// RecordFill takes information about a filled order and updates the corresponding
// database record, if any exists.
func (indexer *Indexer) RecordFill(fillRecord *FillRecord) error {
	hashBytes, err := hex.DecodeString(strings.TrimPrefix(fillRecord.OrderHash, "0x"))
	if err != nil {
		return err
	}
	if fillRecord.FilledTakerAssetAmount == "" {
		fillRecord.FilledTakerAssetAmount = "0"
	}
	amountFilled, ok := new(big.Int).SetString(fillRecord.FilledTakerAssetAmount, 10)
	if !ok {
		return fmt.Errorf("FilledTakerAssetAmount could not be parsed as intger: '%v'", fillRecord.FilledTakerAssetAmount)
	}
	dbOrder := &Order{}
	dbOrder.Initialize()
	indexer.db.Model(&Order{}).Where("order_hash = ?", hashBytes).First(dbOrder)
	totalFilled := dbOrder.TakerAssetAmountFilled.Big()
	copy(dbOrder.TakerAssetAmountFilled[:], abi.U256(totalFilled.Add(totalFilled, amountFilled)))
	dbOrder.Cancelled = dbOrder.Cancelled || fillRecord.Cancel
	return dbOrder.Save(indexer.db, dbOrder.Status).Error
}

// RecordSpend takes information about a token transfer, and updates any
// orders that might have become unfillable as a result of the transfer.
func (indexer *Indexer) RecordSpend(makerAddress, tokenAddress, zrxAddress *types.Address, assetData types.AssetData, balance *types.Uint256) error {
	// NOTE: Right now we're doing this as a single check/update. Eventually it
	// might make sense to do a check against a read replica, and the update
	// against the write node if the check passes. It's more work over-all, but
	// if the write node is a major bottleneck, it could probably take a good
	// bit of pressure off.
	var query *gorm.DB
	if len(assetData) == 0 {
		query = indexer.db.Model(&Order{}).Where("status = ? AND maker_asset_address = ? AND maker = ? AND ? < maker_asset_remaining", StatusOpen, tokenAddress, makerAddress, balance)
	} else {
		query = indexer.db.Model(&Order{}).Where("status = ? AND maker_asset_data = ? AND maker = ? AND ? < maker_asset_remaining", StatusOpen, assetData, makerAddress, balance)
	}
	if(bytes.Equal(tokenAddress[:], zrxAddress[:])) {
		query = query.Or("maker = ? AND ? < maker_fee_remaining", makerAddress, balance)
	}
	return query.Update("status", indexer.status).Error
}

func (indexer *Indexer) RecordCancellation(cancellation *Cancellation) error {
	if err := cancellation.Save(indexer.db).Error; err != nil {
		return err
	}
	return indexer.db.Model(&Order{}).Where(
		"status = ? AND maker = ? AND sender_address = ? AND salt < ?", StatusOpen, cancellation.Maker, cancellation.Sender, cancellation.Epoch,
	).Update("status", indexer.status).Error
}

func NewIndexer(db *gorm.DB, status int64) *Indexer {
	return &Indexer{db, status}
}
