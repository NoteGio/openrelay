package allowance_test


import (
	"encoding/json"
	"encoding/hex"
	"math/big"
	"testing"
	"github.com/notegio/openrelay/monitor/allowance"
	"github.com/notegio/openrelay/monitor/blocks"
	"github.com/notegio/openrelay/monitor/blocks/mock"
	"github.com/notegio/openrelay/channels"
	orCommon "github.com/notegio/openrelay/common"
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
	ctrAddress := common.HexToAddress("0x1d7022f5b17d2f8b695918fb48fa1089c9f85401")
	senderAddress, _ := orCommon.HexToAddress("0x5409ed021d9299bf6814279a6a1411a7e866a631")
	tokenProxyAddress, _ := orCommon.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48")
	approvalTopic := &big.Int{}
	approvalTopic.SetString("8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", 16)
	topics := []common.Hash{
		common.BigToHash(approvalTopic),
		common.BigToHash(senderAddress.Big()),
		common.BigToHash(tokenProxyAddress.Big()),
	}
	data := common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	return buildLog(ctrAddress, topics, data[:])
}

func TestBloom(t *testing.T) {
	bloomBytes, _ := hex.DecodeString("00000000000000000000080000800000000000000000000000000000000000000000000000000000000000000000000000000000000004000000080000200000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000010008000000000000000002000000000000000000000000008000000000000")
	approvalTopic := &big.Int{}
	approvalTopic.SetString("8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", 16)
	proxyAddress, _ := big.NewInt(0).SetString("0x1dc4c1cefef38a777b15aa20260a54e584b16c48", 0)
	bloom := types.BytesToBloom(bloomBytes)
	if !types.BloomLookup(bloom, approvalTopic) {
		t.Errorf("Bloom filter didn't match approval")
	}
	if !types.BloomLookup(bloom, proxyAddress) {
		t.Errorf("Bloom filter didn't match proxy")
	}
	if types.BloomLookup(bloom, &big.Int{}) {
		t.Errorf("Bloom filter shouldn't have matched empty integer")
	}
}

func TestAllowanceFromBlock(t *testing.T) {
	testLog := allowanceLog()
	bloom := types.BytesToBloom(types.LogsBloom([]*types.Log{testLog}).Bytes())
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
	addr, _ := big.NewInt(0).SetString("0x1dc4c1cefef38a777b15aa20260a54e584b16c48", 0)
	consumerChannel.AddConsumer(allowance.NewAllowanceBlockConsumer(
		addr,
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
	if sr.TokenAddress != "0x1d7022f5b17d2f8b695918fb48fa1089c9f85401" {
		t.Errorf("Unexpected token address, got '%v'", sr.TokenAddress)
	}
	if sr.SpenderAddress != "0x5409ed021d9299bf6814279a6a1411a7e866a631" {
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
	addr, _ := big.NewInt(0).SetString("0x1dc4c1cefef38a777b15aa20260a54e584b16c48", 0)
	consumerChannel.AddConsumer(allowance.NewAllowanceBlockConsumer(
		addr,
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
