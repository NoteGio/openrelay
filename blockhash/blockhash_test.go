package blockhash_test

import (
	"github.com/notegio/openrelay/blockhash"
	"github.com/notegio/openrelay/channels"
	"testing"
	"time"
)

func TestMockChannelGetBlockhash(t *testing.T) {
	publisher, consumerChannel := channels.MockChannel()
	blockHash := blockhash.NewChanneledBlockHash(consumerChannel)
	if value := blockHash.Get(); value != "initializing" {
		t.Errorf("Expected blockHash to be initializing")
	}
	publisher.Publish("new value")
	time.Sleep(100 * time.Millisecond)
	if value := blockHash.Get(); value != "new value" {
		t.Errorf("Expected blockHash to be updated value, got '%v'", value)
	}
}
