package subscriptions_test

import (
	"testing"
	"github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/channels"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/subscriptions"
	"github.com/notegio/openrelay/types"
	// "log"
	// "fmt"
)


func TestSubscriptionConsumer(t *testing.T) {
	incomingPublisher, incomingConsumerChannel := channels.MockChannel()
	outgoingPublisher, deliveries := channels.MockPublisher()
	manager := &subscriptions.SubscriptionManager{}
	mapping := make(map[int64][]*types.Address)
	sampleExchange, _ := common.HexToAddress("0x6666666666666666666666666666666666666666")
	mapping[1] = []*types.Address{sampleExchange}
	lookup := &MockExchangeLookup{mapping}
	consumer := subscriptions.NewSubscriptionConsumer(manager, outgoingPublisher, nil, lookup)
	incomingConsumerChannel.AddConsumer(consumer)
	incomingConsumerChannel.StartConsuming()
	defer incomingConsumerChannel.StopConsuming()
	incomingPublisher.Publish(`{
	    "type": "subscribe",
	    "channel": "orders",
	    "requestId": "123e4567-e89b-12d3-a456-426655440000"
	}`)
	torder := &types.Order{}
	torder.Initialize()
	order := &dbModule.Order{Order: *torder}
	order.TakerAssetAmount = common.Int64ToUint256(1)
	order.Populate()
	channels.MockFinish(incomingConsumerChannel, 1)
	manager.Publish(order)
	d := <-deliveries
	if d.Payload() != `{"type":"update","channel":"orders","requestId":"123e4567-e89b-12d3-a456-426655440000","payload":[{"order":{"makerAddress":"0x0000000000000000000000000000000000000000","takerAddress":"0x0000000000000000000000000000000000000000","makerAssetData":"0x0000000000000000000000000000000000000000","takerAssetData":"0x0000000000000000000000000000000000000000","feeRecipientAddress":"0x0000000000000000000000000000000000000000","exchangeAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"0","takerAssetAmount":"1","makerFee":"0","takerFee":"0","expirationTimeSeconds":"0","salt":"0","signature":""},"metaData":{"hash":"0xadfaa1d67d27cc9240ab3a90bf7b1682eb683e16873b6a88f95294d126b5e6c1","feeRate":0,"status":0,"takerAssetAmountRemaining":"1"}}]}` {
		t.Errorf("Unexpected valud: %v", d.Payload())
	}
}
