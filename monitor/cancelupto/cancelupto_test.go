package cancelupto_test


import (
	"encoding/json"
	"encoding/hex"
	"math/big"
	"testing"
	"github.com/notegio/openrelay/monitor/cancelupto"
	"github.com/notegio/openrelay/monitor/blocks"
	"github.com/notegio/openrelay/monitor/blocks/mock"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/db"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	// "log"
)

type testConsumer struct {
	channel chan string
}

func (consumer *testConsumer) Consume(msg channels.Delivery) {
	consumer.channel <- msg.Payload()
}

func newTestConsumer() *testConsumer {
	return &testConsumer{make(chan string, 5)}
}

func buildLog(address common.Address, topics []common.Hash, data []byte) *types.Log {
	return &types.Log{
		Address: address,
		Topics: topics,
		Data: data,
	}
}

func cancelLog() *types.Log {
	ctrAddress := common.HexToAddress("0xb65619b82c4d385de0c5b4005452c2fdee0f86d1")
	cancelUpToTopic := &big.Int{}
	cancelUpToTopic.SetString("82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0", 16)
	topics := []common.Hash{
		common.BigToHash(cancelUpToTopic),
		common.HexToHash("0x324454186bb728a3ea55750e0618ff1b18ce6cf8"),
		common.HexToHash("0x00"),
	}
	data, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000002")
	return buildLog(ctrAddress, topics, data[:])
}

func TestCancelUpToFromBlock(t *testing.T) {
	testLog := cancelLog()
	bloom := types.Bloom{}
	bloom.Add(new(big.Int).SetBytes(testLog.Address[:]))
	for _, topic := range testLog.Topics {
		bloom.Add(new(big.Int).SetBytes(topic[:]))
	}
	mb := &blocks.MiniBlock{
		common.Hash{},
		big.NewInt(0),
		bloom,
	}
	srcPublisher, consumerChannel := channels.MockChannel()
	destPublisher, destConsumerChannel := channels.MockChannel()
	data, err := json.Marshal(mb)
	if err != nil {
		t.Errorf(err.Error())
	}
	tc := newTestConsumer()
	destConsumerChannel.AddConsumer(tc)
	destConsumerChannel.StartConsuming()
	defer destConsumerChannel.StopConsuming()
	if err != nil { t.Fatalf(err.Error()) }
	consumerChannel.AddConsumer(cancelupto.NewCancelUpToBlockConsumer(
		common.HexToAddress("0xb65619b82c4d385de0c5b4005452c2fdee0f86d1").Big(),
		mock.NewMockLogFilterer([]types.Log{*testLog}),
		destPublisher,
	))
	consumerChannel.StartConsuming()
	defer consumerChannel.StopConsuming()
	srcPublisher.Publish(string(data))
	payload := <-tc.channel
	cancellation := &db.Cancellation{}
	err = json.Unmarshal([]byte(payload), cancellation)
	if err != nil {
		t.Errorf(err.Error())
	}
	if cancellation.Maker.String() != "0x324454186bb728a3ea55750e0618ff1b18ce6cf8" {
		t.Errorf("Unexpected maker, got '%#x'", cancellation.Maker)
	}
	if cancellation.Sender.String() != "0x0000000000000000000000000000000000000000" {
		t.Errorf("Unexpected Sender, got '%#x'", cancellation.Sender)
	}
	if cancellation.Epoch.String() != "2" {
		t.Errorf("Unexpected epoch, got '%v'", cancellation.Epoch)
	}
}
