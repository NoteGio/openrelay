package subscriptions_test

import (
	"context"
	"testing"
	"github.com/notegio/openrelay/channels"
	"github.com/jinzhu/gorm"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/subscriptions"
	"github.com/notegio/openrelay/types"
	"github.com/gorilla/websocket"
	"time"
	"log"
	"fmt"
	"os"
)

func getDb(t *testing.T) (*gorm.DB, error) {
	connectionString := fmt.Sprintf(
		"postgres://%v@%v",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_HOST"),
	)
	if connectionString == "postgres://@" {
		log.Printf("DB Not configured. Skipping.")
		t.Skip()
		return nil, fmt.Errorf("No DB Configured")
	}
	db, err := dbModule.GetDB(connectionString, os.Getenv("POSTGRES_PASSWORD"))
	// db.LogMode(true)
	return db, err
}

func TestWebsocketSubscriptionConsumer(t *testing.T) {
	incomingPublisher, incomingConsumerChannel := channels.MockChannel()
	manager := subscriptions.NewWebsocketSubscriptionManager()
	incomingConsumerChannel.AddConsumer(manager)
	incomingConsumerChannel.StartConsuming()
	db, err := getDb(t)
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
	quit, err := manager.ListenForSubscriptions(4321, tx)
	defer quit()
	if err != nil {
		t.Fatalf(err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	c, resp, err := websocket.DefaultDialer.DialContext(ctx, "ws://localhost:4321/v3/", nil)

	defer cancel()
	if err != nil {
		content := []byte{}
		statusCode := 0
		if resp != nil {
			resp.Body.Read(content[:])
			statusCode = resp.StatusCode
		}
		t.Fatalf("%v - (%v) %v", err.Error(), statusCode, content)
	}

	c.WriteMessage(websocket.TextMessage, []byte(`{
	    "type": "subscribe",
	    "channel": "orders",
	    "requestId": "123e4567-e89b-12d3-a456-426655440000"
	}`))

	time.Sleep(50 * time.Millisecond)

	torder := &types.Order{}
	torder.Initialize()
	incomingPublisher.Publish(string(torder.Bytes()))
	channels.MockFinish(incomingConsumerChannel, 1)

	mtype, p, err := c.ReadMessage()
	if err != nil {
		t.Errorf(err.Error())
	}
	if mtype != websocket.TextMessage {
		t.Errorf("Unexpected message type %v", mtype)
	}

	if string(p) != `{"type":"update","channel":"orders","requestId":"123e4567-e89b-12d3-a456-426655440000","payload":[{"order":{"chainId":0,"makerAddress":"0x0000000000000000000000000000000000000000","takerAddress":"0x0000000000000000000000000000000000000000","makerAssetData":"","takerAssetData":"","makerFeeAssetData":"","takerFeeAssetData":"","feeRecipientAddress":"0x0000000000000000000000000000000000000000","exchangeAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"0","takerAssetAmount":"0","makerFee":"0","takerFee":"0","expirationTimeSeconds":"0","salt":"0","signature":""},"metaData":{"hash":"0xa9782b83f1c408f8f4fe049c666f253bd1e6a8f5a8bba27b1b7eb57e49759620","feeRate":0,"status":1,"takerAssetAmountRemaining":"0"}}]}` {
		t.Errorf("Unexpected value: %v", string(p))
	}
}
