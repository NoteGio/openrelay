package db

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/types"
	"log"
)

type SpendRecord struct {
	TokenAddress    string `json:"tokenAddress"`
	AssetData       string `json:"assetData"`
	SpenderAddress  string `json:"spenderAddress"`
	Balance         string `json:"balance"`
}

type RecordSpendConsumer struct {
	idx *Indexer
	s   common.Semaphore
}

func (consumer *RecordSpendConsumer) Consume(msg channels.Delivery) {
	consumer.s.Acquire()
	go func(){
		defer consumer.s.Release()
		spendRecord := &SpendRecord{}
		err := json.Unmarshal([]byte(msg.Payload()), spendRecord)
		if err != nil {
			log.Printf("Failed to parse JSON: %v", err.Error())
			msg.Reject()
			return
		}
		spenderAddress, err := common.HexToAddress(spendRecord.SpenderAddress)
		if err != nil {
			log.Printf("Failed to record spend: '%v', '%v'", msg.Payload(), err.Error())
			msg.Reject()
			return
		}
		tokenAddress, err := common.HexToAddress(spendRecord.TokenAddress)
		if err != nil {
			log.Printf("Failed to record spend: '%v', '%v'", msg.Payload(), err.Error())
			msg.Reject()
			return
		}
		assetData, err := common.HexToAssetData(spendRecord.AssetData)
		if err != nil {
			log.Printf("Failed to record spend: '%v', '%v'", msg.Payload(), err.Error())
			msg.Reject()
			return
		}
		balance, err := types.IntStringToUint256(spendRecord.Balance)
		if err := consumer.idx.RecordSpend(spenderAddress, tokenAddress, assetData, balance); err == nil {
			msg.Ack()
		} else {
			log.Printf("Failed to record spend: '%v', '%v'", msg.Payload(), err.Error())
			msg.Reject()
			return
		}
	}()
}

func NewRecordSpendConsumer(db *gorm.DB, concurrency int, publisher channels.Publisher) *RecordSpendConsumer {
	return &RecordSpendConsumer{NewIndexer(db, StatusUnfunded, publisher), make(common.Semaphore, concurrency)}
}
