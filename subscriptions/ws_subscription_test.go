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
	quit, err := manager.ListenForSubscriptions(1234, tx)
	defer quit()
	if err != nil {
		t.Fatalf(err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	c, resp, err := websocket.DefaultDialer.DialContext(ctx, "ws://localhost:4321/v2/", nil)

	defer cancel()
	if err != nil {
		content := []byte{}
		resp.Body.Read(content[:])
		t.Fatalf("%v - (%v) %v", err.Error(), resp.StatusCode, content)
	}

	c.WriteMessage(websocket.BinaryMessage, []byte(`{
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
	if mtype != websocket.BinaryMessage {
		t.Errorf("Unexpected message type %v", mtype)
	}

	if string(p) != `{"type":"update","channel":"orders","requestId":"123e4567-e89b-12d3-a456-426655440000","payload":[{"order":{"makerAddress":"0x0000000000000000000000000000000000000000","takerAddress":"0x0000000000000000000000000000000000000000","makerAssetData":"0x0000000000000000000000000000000000000000","takerAssetData":"0x0000000000000000000000000000000000000000","feeRecipientAddress":"0x0000000000000000000000000000000000000000","exchangeAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"0","takerAssetAmount":"1","makerFee":"0","takerFee":"0","expirationTimeSeconds":"0","salt":"0","signature":""},"metaData":{"hash":"0xadfaa1d67d27cc9240ab3a90bf7b1682eb683e16873b6a88f95294d126b5e6c1","feeRate":0,"status":0,"takerAssetAmountRemaining":"1"}}]}` {
		t.Errorf("Unexpected value: %v", string(p))
	}
}
