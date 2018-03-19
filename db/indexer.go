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
)

type FillRecord struct {
	OrderHash                 string `json:"orderHash"`
	FilledTakerTokenAmount    string `json:"filledTakerTokenAmount"`
	CancelledTakerTokenAmount string `json:"cancelledTakerTokenAmount"`
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
	if fillRecord.FilledTakerTokenAmount == "" {
		fillRecord.FilledTakerTokenAmount = "0"
	}
	if fillRecord.CancelledTakerTokenAmount == "" {
		fillRecord.CancelledTakerTokenAmount = "0"
	}
	amountFilled, ok := new(big.Int).SetString(fillRecord.FilledTakerTokenAmount, 10)
	if !ok {
		return fmt.Errorf("FilledTakerTokenAmount could not be parsed as intger: '%v'", fillRecord.FilledTakerTokenAmount)
	}
	amountCancelled, ok := new(big.Int).SetString(fillRecord.CancelledTakerTokenAmount, 10)
	if !ok {
		return fmt.Errorf("CancelledTakerTokenAmount could not be parsed as intger: '%v'", fillRecord.CancelledTakerTokenAmount)
	}
	dbOrder := &Order{}
	dbOrder.Initialize()
	indexer.db.Model(&Order{}).Where("order_hash = ?", hashBytes).First(dbOrder)
	totalFilled := new(big.Int).SetBytes(dbOrder.TakerTokenAmountFilled[:])
	totalCancelled := new(big.Int).SetBytes(dbOrder.TakerTokenAmountCancelled[:])
	copy(dbOrder.TakerTokenAmountFilled[:], abi.U256(totalFilled.Add(totalFilled, amountFilled)))
	copy(dbOrder.TakerTokenAmountCancelled[:], abi.U256(totalCancelled.Add(totalCancelled, amountCancelled)))
	return dbOrder.Save(indexer.db, dbOrder.Status).Error
}

// RecordSpend takes information about a token transfer, and updates any
// orders that might have become unfillable as a result of the transfer.
func (indexer *Indexer) RecordSpend(makerAddress, tokenAddress, zrxAddress *types.Address, balance *types.Uint256) error {
	// NOTE: Right now we're doing this as a single check/update. Eventually it
	// might make sense to do a check against a read replica, and the update
	// against the write node if the check passes. It's more work over-all, but
	// if the write node is a major bottleneck, it could probably take a good
	// bit of pressure off.
	query := indexer.db.Model(&Order{}).Where("status = ? AND maker_token = ? AND maker = ? AND ? < maker_token_remaining", StatusOpen, tokenAddress, makerAddress, balance)
	if(bytes.Equal(tokenAddress[:], zrxAddress[:])) {
		query = query.Or("maker = ? AND ? < maker_fee_remaining", makerAddress, balance)
	}
	return query.Update("status", indexer.status).Error
}

func NewIndexer(db *gorm.DB, status int64) *Indexer {
	return &Indexer{db, status}
}
