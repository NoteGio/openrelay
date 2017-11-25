package db

import (
	"encoding/hex"
	"math/big"
	"github.com/notegio/openrelay/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/jinzhu/gorm"
	"fmt"
	"strings"
)

type FillRecord struct {
	OrderHash string `json:"orderHash"`
	FilledTakerTokenAmount  string `json:"filledTakerTokenAmount"`
	CancelledTakerTokenAmount string `json:"cancelledTakerTokenAmount"`
}

type Indexer struct {
	db *gorm.DB
	status int64
}

func (indexer *Indexer) Index(order *types.Order) (error) {
	dbOrder := Order{}
	dbOrder.Order = *order
	return dbOrder.Save(indexer.db, indexer.status).Error
}

func (indexer *Indexer) RecordFill(fillRecord *FillRecord) (error) {
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

func NewIndexer(db *gorm.DB, status int64) (*Indexer) {
	return &Indexer{db, status}
}
