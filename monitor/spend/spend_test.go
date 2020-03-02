package spend_test


import (
	"encoding/json"
	"encoding/hex"
	"math/big"
	"testing"
	"github.com/notegio/openrelay/funds/balance"
	"github.com/notegio/openrelay/monitor/spend"
	"github.com/notegio/openrelay/monitor/blocks"
	"github.com/notegio/openrelay/monitor/blocks/mock"
	orTypes "github.com/notegio/openrelay/types"
	orCommon "github.com/notegio/openrelay/common"
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

func spendLog() *types.Log {
	ctrAddress := common.HexToAddress("0x3495ffcee09012ab7d827abf3e3b3ae428a38443")
	senderAddress, _ := orCommon.HexToAddress("0x34ab4a96678c4de8eb34597dbbcf09c27d9bc79d")
	receiverAddress, _ := orCommon.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093")
	spendTopic := &big.Int{}
	spendTopic.SetString("ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", 16)
	topics := []common.Hash{
		common.BigToHash(spendTopic),
		common.BigToHash(senderAddress.Big()),
		common.BigToHash(receiverAddress.Big()),
	}
	data := common.HexToHash("0x0000000000000000000000000000000000000000000000006f05b59d3b200000")
	return buildLog(ctrAddress, topics, data[:])
}

func TestBloom(t *testing.T) {
	bloomBytes, _ := hex.DecodeString("20000000000000000000000000080000000000000000020000000000000000000000000000000000000000008000000000000000000000000000000000000004000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000001000000001000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000602000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	// bloomBytes, _ := hex.DecodeString("00000000000000000000080000800000000000000000000000000000000000000000000000000000000000000000000000000000000004000000080000200000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000010008000000000000000002000000000000000000000000008000000000000")
	spendTopic := &big.Int{}
	spendTopic.SetString("ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", 16)
	bloom := types.BytesToBloom(bloomBytes)
	if !types.BloomLookup(bloom, spendTopic) {
		t.Errorf("Bloom filter didn't match spend")
	}
	if types.BloomLookup(bloom, &big.Int{}) {
		t.Errorf("Bloom filter shouldn't have matched empty integer")
	}
}

func TestSpendFromBlock(t *testing.T) {
	testLog := spendLog()
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
	tokenProxyAddress := &orTypes.Address{}
	tokenProxyBytes := common.HexToAddress("0x3333333333333333333333333333333333333333")
	copy(tokenProxyAddress[:], tokenProxyBytes[:])
	tokenBytes := common.HexToAddress("0x3495ffcee09012ab7d827abf3e3b3ae428a38443")
	tokenAddress := &orTypes.Address{}
	spenderBytes := common.HexToAddress("0x34ab4a96678c4de8eb34597dbbcf09c27d9bc79d")
	spenderAddress := &orTypes.Address{}
	copy(tokenAddress[:], tokenBytes[:])
	copy(spenderAddress[:], spenderBytes[:])
	balanceMap := make(map[string]map[orTypes.Address]*big.Int)
	balanceMap[string(orCommon.ToERC20AssetData(tokenAddress))] = make(map[orTypes.Address]*big.Int)
	balanceMap[string(orCommon.ToERC20AssetData(tokenAddress))][*spenderAddress] = big.NewInt(0)
	consumerChannel.AddConsumer(spend.NewSpendBlockConsumer(tokenProxyAddress,
		"0x4444444444444444444444444444444444444444",
		mock.NewMockLogFilterer([]types.Log{*testLog}),
		destPublisher,
		balance.NewMockBalanceChecker(balanceMap),
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
	if sr.TokenAddress != "0x3495ffcee09012ab7d827abf3e3b3ae428a38443" {
		t.Errorf("Unexpected token address, got '%v'", sr.TokenAddress)
	}
	if sr.SpenderAddress != "0x34ab4a96678c4de8eb34597dbbcf09c27d9bc79d" {
		t.Errorf("Unexpected spender address, got '%v'", sr.TokenAddress)
	}
	if sr.ZrxToken != "0x4444444444444444444444444444444444444444" {
		t.Errorf("Unexpected zrx token address, got '%v'", sr.TokenAddress)
	}
	balance, _ := new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
	if sr.Balance != balance.String() {
		t.Errorf("Unexpected token address, got '%v'", sr.TokenAddress)
	}
}
func TestNoSpendInBlock(t *testing.T) {
	testLog := spendLog()
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
	tokenProxyAddress := &orTypes.Address{}
	tokenProxyBytes := common.HexToAddress("0x3333333333333333333333333333333333333333")
	copy(tokenProxyAddress[:], tokenProxyBytes[:])
	consumerChannel.AddConsumer(spend.NewSpendBlockConsumer(tokenProxyAddress,
		"0x4444444444444444444444444444444444444444",
		mock.NewMockLogFilterer([]types.Log{*testLog}),
		destPublisher,
		balance.NewMockBalanceChecker(make(map[string]map[orTypes.Address]*big.Int)),
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
