package ws_test

import (
	"context"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/channels/ws"
	"github.com/gorilla/websocket"
	"testing"
	"time"
	// "log"
)

type TestConsumer struct {
	channel *ws.WebsocketChannel
}

func (consumer *TestConsumer) Consume(delivery channels.Delivery) {
	if delivery.Payload() == "quit" {
		consumer.channel.StopConsuming()
		delivery.Ack()
	}
	consumer.channel.Publish(delivery.Payload())
	delivery.Ack()
}

func TestGetChannels(t *testing.T) {
	clean := false
	channels, quit := ws.GetChannels(4321, nil, func() { clean = true })
	go func() {
		for channel := range channels {
			channel.AddConsumer(&TestConsumer{channel})
			channel.StartConsuming()
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	c, resp, err := websocket.DefaultDialer.DialContext(ctx, "ws://localhost:4321/v2/", nil)
	defer cancel()
	if err != nil {
		content := []byte{}
		resp.Body.Read(content[:])
		t.Fatalf("%v - (%v) %v", err.Error(), resp.StatusCode, content)
	}
	if err := c.WriteMessage(websocket.BinaryMessage, []byte("ping")); err != nil {
		t.Errorf(err.Error())
	}
	mtype, p, err := c.ReadMessage()
	if err != nil {
		t.Errorf(err.Error())
	}
	if mtype != websocket.BinaryMessage {
		t.Errorf("Unexpected message type %v", mtype)
	}
	if string(p) != "ping" {
		t.Errorf("Unexpected message: %v", string(p))
	}
	c.Close()
	if err := quit(); err != nil {
		t.Errorf(err.Error())
	}
	time.Sleep(50 * time.Millisecond)
	if clean != true {
		t.Errorf("Should have cleaned up")
	}
}
