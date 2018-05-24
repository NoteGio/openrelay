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
	SpenderAddress  string `json:"spenderAddress"`
	ZrxToken        string `json:"zrxToken"`
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
		zrxToken, err := common.HexToAddress(spendRecord.ZrxToken)
		if err != nil {
			log.Printf("Failed to record spend: '%v', '%v'", msg.Payload(), err.Error())
			msg.Reject()
			return
		}
		balance, err := types.IntStringToUint256(spendRecord.Balance)
		if err := consumer.idx.RecordSpend(spenderAddress, tokenAddress, zrxToken, balance); err == nil {
			msg.Ack()
			} else {
				log.Printf("Failed to record spend: '%v', '%v'", msg.Payload(), err.Error())
				msg.Reject()
				return
			}
	}()
}

func NewRecordSpendConsumer(db *gorm.DB, concurrency int) *RecordSpendConsumer {
	return &RecordSpendConsumer{NewIndexer(db, StatusUnfunded), make(common.Semaphore, concurrency)}
}
