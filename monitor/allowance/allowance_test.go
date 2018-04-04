package allowance_test


import (
	"encoding/json"
	"math/big"
	"testing"
	"github.com/notegio/openrelay/monitor/allowance"
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

func allowanceLog() *types.Log {
	ctrAddress := common.HexToAddress("0x1111111111111111111111111111111111111111")
	senderAddress := common.HexToAddress("0x2222222222222222222222222222222222222222")
	tokenProxyAddress := common.HexToAddress("0x3333333333333333333333333333333333333333")
	approvalTopic := &big.Int{}
	approvalTopic.SetString("8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", 16)
	topics := []common.Hash{
		common.BigToHash(approvalTopic),
		common.BytesToHash(senderAddress[:]),
		common.BigToHash(tokenProxyAddress.Big()),
	}
	data := common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	return buildLog(ctrAddress, topics, data[:])
}

func TestAllowanceFromBlock(t *testing.T) {
	testLog := allowanceLog()
	bloom := types.Bloom{}
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
	consumerChannel.AddConsumer(allowance.NewAllowanceBlockConsumer(
		common.HexToAddress("0x3333333333333333333333333333333333333333").Big(),
		"0x4444444444444444444444444444444444444444",
		mock.NewMockLogFilterer([]types.Log{*testLog}),
		destPublisher,
	))
	consumerChannel.StartConsuming()
	defer consumerChannel.StopConsuming()
	srcPublisher.Publish(string(data))
	payload := <-tc.channel
	sr := &db.SpendRecord{}
	err = json.Unmarshal([]byte(payload), sr)
	if err != nil {
		t.Errorf(err.Error())
	}
	if sr.TokenAddress != "0x1111111111111111111111111111111111111111" {
		t.Errorf("Unexpected token address, got '%v'", sr.TokenAddress)
	}
	if sr.SpenderAddress != "0x2222222222222222222222222222222222222222" {
		t.Errorf("Unexpected token address, got '%v'", sr.TokenAddress)
	}
	if sr.ZrxToken != "0x4444444444444444444444444444444444444444" {
		t.Errorf("Unexpected token address, got '%v'", sr.TokenAddress)
	}
	balance, _ := new(big.Int).SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	if sr.Balance != balance.String() {
		t.Errorf("Unexpected token address, got '%v'", sr.TokenAddress)
	}
}
func TestNoAllowanceInBlock(t *testing.T) {
	mb := &blocks.MiniBlock{
		common.Hash{},
		big.NewInt(0),
		types.Bloom{},
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
	consumerChannel.AddConsumer(allowance.NewAllowanceBlockConsumer(
		common.HexToAddress("0x3333333333333333333333333333333333333333").Big(),
		"0x4444444444444444444444444444444444444444",
		mock.NewMockLogFilterer([]types.Log{}),
		destPublisher,
	))
	consumerChannel.StartConsuming()
	defer consumerChannel.StopConsuming()
	srcPublisher.Publish(string(data))
	select {
		case _ = <-tc.channel:
			t.Errorf("Should not have gotten a message")
		default:
	}
}
