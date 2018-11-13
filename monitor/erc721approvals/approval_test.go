package erc721approvals_test


import (
	"encoding/json"
	"encoding/hex"
	"math/big"
	"testing"
	"github.com/notegio/openrelay/monitor/erc721approvals"
	"github.com/notegio/openrelay/monitor/blocks"
	"github.com/notegio/openrelay/monitor/blocks/mock"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/db"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"time"
	"log"
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

func allowanceLog(tokenProxyHex string) *types.Log {
	ctrAddress := common.HexToAddress("0x1d7022f5b17d2f8b695918fb48fa1089c9f85401")
	senderAddress := common.HexToAddress("0x5409ed021d9299bf6814279a6a1411a7e866a631")
	tokenProxyAddress := common.HexToAddress(tokenProxyHex)
	approvalTopic := &big.Int{}
	approvalTopic.SetString("8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", 16)
	topics := []common.Hash{
		common.BigToHash(approvalTopic),
		common.BigToHash(senderAddress.Big()),
		common.BigToHash(tokenProxyAddress.Big()),
		common.BigToHash(big.NewInt(1)),
	}
	return buildLog(ctrAddress, topics, []byte{})
}

func approveAllLog(tokenProxyHex string, approved int64) *types.Log {
	ctrAddress := common.HexToAddress("0x1d7022f5b17d2f8b695918fb48fa1089c9f85401")
	senderAddress := common.HexToAddress("0x5409ed021d9299bf6814279a6a1411a7e866a631")
	tokenProxyAddress := common.HexToAddress(tokenProxyHex)
	approvalTopic := &big.Int{}
	approvalTopic.SetString("17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31", 16)
	topics := []common.Hash{
		common.BigToHash(approvalTopic),
		common.BigToHash(senderAddress.Big()),
		common.BigToHash(tokenProxyAddress.Big()),
	}
	data := common.BigToHash(big.NewInt(approved))
	return buildLog(ctrAddress, topics, data[:])
}

func TestBloom(t *testing.T) {
	bloomBytes, _ := hex.DecodeString("00000000000000000000080000800000000000000000000000000000000000000000000000000000000000000000000000000000000004000000080000200000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000010008000000000000000002000000000000000000000000008000000000000")
	approvalTopic := &big.Int{}
	approvalTopic.SetString("8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", 16)
	proxyAddress := common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48").Big()
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

func TestAllowanceFromBlockMatch(t *testing.T) {
	testLog := allowanceLog("0x1dc4c1cefef38a777b15aa20260a54e584b16c48")
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
	consumerChannel.AddConsumer(erc721approvals.NewAllowanceBlockConsumer(
		common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48").Big(),
		"0x4444444444444444444444444444444444444444",
		mock.NewMockLogFilterer([]types.Log{*testLog}),
		destPublisher,
	))
	consumerChannel.StartConsuming()
	defer consumerChannel.StopConsuming()
	srcPublisher.Publish(string(data))
	time.Sleep(time.Second)
	select {
	case _ = <-tc.channel:
		t.Errorf("Channel should have been empty")
		default:
	}
}
func TestAllowanceFromBlock(t *testing.T) {
	log.Printf("TestAllowanceFromBlock")
	testLog := allowanceLog("0x000001cefef38a777b15aa20260a54e584b16c48")
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
	consumerChannel.AddConsumer(erc721approvals.NewAllowanceBlockConsumer(
		common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48").Big(),
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
	if sr.TokenAddress != "0x1d7022f5b17d2f8b695918fb48fa1089c9f85401" {
		t.Errorf("Unexpected token address, got '%v'", sr.TokenAddress)
	}
	if sr.AssetData != "0x025717920000000000000000000000001d7022f5b17d2f8b695918fb48fa1089c9f854010000000000000000000000000000000000000000000000000000000000000001" {
		t.Errorf("Unexpected token address, got '%v'", sr.AssetData)
	}
	if sr.SpenderAddress != "0x5409ed021d9299bf6814279a6a1411a7e866a631" {
		t.Errorf("Unexpected token address, got '%v'", sr.SpenderAddress)
	}
	if sr.ZrxToken != "0x4444444444444444444444444444444444444444" {
		t.Errorf("Unexpected token address, got '%v'", sr.ZrxToken)
	}
	balance := big.NewInt(0)
	if sr.Balance != balance.String() {
		t.Errorf("Unexpected token address, got '%v'", sr.Balance)
	}
}
func TestApproveAllFromBlockDisapprove(t *testing.T) {
	testLog := approveAllLog("0x1dc4c1cefef38a777b15aa20260a54e584b16c48", 0)
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
	consumerChannel.AddConsumer(erc721approvals.NewAllowanceBlockConsumer(
		common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48").Big(),
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
	if sr.TokenAddress != "0x1d7022f5b17d2f8b695918fb48fa1089c9f85401" {
		t.Errorf("Unexpected token address, got '%v'", sr.TokenAddress)
	}
	if sr.AssetData != "" {
		t.Errorf("Unexpected token address, got '%v'", sr.AssetData)
	}
	if sr.SpenderAddress != "0x5409ed021d9299bf6814279a6a1411a7e866a631" {
		t.Errorf("Unexpected token address, got '%v'", sr.SpenderAddress)
	}
	if sr.ZrxToken != "0x4444444444444444444444444444444444444444" {
		t.Errorf("Unexpected token address, got '%v'", sr.ZrxToken)
	}
	balance := big.NewInt(0)
	if sr.Balance != balance.String() {
		t.Errorf("Unexpected token address, got '%v'", sr.Balance)
	}
}
func TestApproveAllFromBlock(t *testing.T) {
	log.Printf("TestApproveAllFromBlock")
	testLog := approveAllLog("0x1dc4c1cefef38a777b15aa20260a54e584b16c48", 1)
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
	consumerChannel.AddConsumer(erc721approvals.NewAllowanceBlockConsumer(
		common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48").Big(),
		"0x4444444444444444444444444444444444444444",
		mock.NewMockLogFilterer([]types.Log{*testLog}),
		destPublisher,
	))
	consumerChannel.StartConsuming()
	defer consumerChannel.StopConsuming()
	srcPublisher.Publish(string(data))
	time.Sleep(time.Second)
	select {
	case _ = <-tc.channel:
		t.Errorf("Channel should have been empty")
		default:
	}
}
func TestNoAllowanceInBlock(t *testing.T) {
	log.Printf("TestNoAllowanceInBlock")
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
	consumerChannel.AddConsumer(erc721approvals.NewAllowanceBlockConsumer(
		common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48").Big(),
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
