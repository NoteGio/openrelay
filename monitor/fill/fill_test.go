package fill_test


import (
	"encoding/json"
	"encoding/hex"
	"math/big"
	"testing"
	"github.com/notegio/openrelay/monitor/fill"
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

func fillLog() *types.Log {
	ctrAddress := common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093")
	fillTopic := &big.Int{}
	fillTopic.SetString("0d0b9391970d9a25552f37d436d2aae2925e2bfe1b2a923754bada030c498cb3", 16)
	topics := []common.Hash{
		common.BigToHash(fillTopic),
		common.HexToHash("0xe793cdaea79d85bd9f19d2c4c52e24fe2caff78c"),
		common.HexToHash("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		common.HexToHash("0x3ef8470e12a03517b9a5827402c04b1bc2c997abed40c75cf40e56462d9cad6e"),
	}
	data, _ := hex.DecodeString("000000000000000000000000f21ec009e3a156f94c2d5ae4353fe361c27661c20000000000000000000000006c6ee5e31d828de241282b9606c8e98ea48526e2000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000176b852995b0d203d37a6000000000000000000000000000000000000000000000000387b8b6fabd7000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f4daf69e0a30af58851c303d47f9627c66d59666fb51f9b39f1a8f92ab900de4")
	return buildLog(ctrAddress, topics, data[:])
}

func TestFillFromBlock(t *testing.T) {
	testLog := fillLog()
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
	consumerChannel.AddConsumer(fill.NewFillBlockConsumer(
		common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093").Big(),
		mock.NewMockLogFilterer([]types.Log{*testLog}),
		destPublisher,
	))
	consumerChannel.StartConsuming()
	defer consumerChannel.StopConsuming()
	srcPublisher.Publish(string(data))
	payload := <-tc.channel
	fr := &db.FillRecord{}
	err = json.Unmarshal([]byte(payload), fr)
	if err != nil {
		t.Errorf(err.Error())
	}
	if fr.OrderHash != "0xf4daf69e0a30af58851c303d47f9627c66d59666fb51f9b39f1a8f92ab900de4" {
		t.Errorf("Unexpected order hash, got '%v'", fr.OrderHash)
	}
	if fr.FilledTakerTokenAmount != "4070000000000000000" {
		t.Errorf("Unexpected filled amount, got '%v'", fr.FilledTakerTokenAmount)
	}
	if fr.CancelledTakerTokenAmount != "0" {
		t.Errorf("Unexpected cancelled amount, got '%v'", fr.CancelledTakerTokenAmount)
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
	consumerChannel.AddConsumer(fill.NewFillBlockConsumer(
		common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093").Big(),
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
