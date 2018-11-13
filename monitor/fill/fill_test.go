package fill_test


import (
	"encoding/json"
	"encoding/hex"
	"math/big"
	"math/rand"
	"os"
	"testing"
	"time"
	"github.com/notegio/openrelay/fillbloom"
	"github.com/notegio/openrelay/monitor/fill"
	"github.com/notegio/openrelay/monitor/blocks"
	"github.com/notegio/openrelay/monitor/blocks/mock"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/db"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	// "log"
	"fmt"
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
	fillTopic.SetString("0bcc4c97732e47d9946f229edb95f5b6323f601300e4690de719993f3c371129", 16)
	topics := []common.Hash{
		common.BigToHash(fillTopic),
		common.HexToHash("0x5409ed021d9299bf6814279a6a1411a7e866a631"),
		common.HexToHash("0x00"),
		common.HexToHash("0x91b419e1cc29695dd4da477967c1b529eaad1591692566778eaf2d4baec3c593"),
	}
	data, _ := hex.DecodeString("000000000000000000000000e36ea790bc9d7ab70c55260c66d52b1eca985f84000000000000000000000000e36ea790bc9d7ab70c55260c66d52b1eca985f84000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001600000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000006ff6c0ff1d68b964901f986d4c9fa3ac68346570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024f47261b0000000000000000000000000653e49e301e508a13237c0ddc98ae7d4cd2667a100000000000000000000000000000000000000000000000000000000")
	return buildLog(ctrAddress, topics, data[:])
}

func TestFillFromBlock(t *testing.T) {
	directory := fmt.Sprintf("/tmp/test-%v", rand.Int())
	os.Mkdir(directory, 0755)
	itemURL := fmt.Sprintf("file://%v/test", directory)
	testLog := fillLog()
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
	if err != nil { t.Fatalf(err.Error()) }
	fillBloom, err := fillbloom.NewFillBloom(itemURL)
	if err != nil { t.Fatalf(err.Error()) }
	consumerChannel.AddConsumer(fill.NewFillBlockConsumer(
		common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093").Big(),
		mock.NewMockLogFilterer([]types.Log{*testLog}),
		destPublisher,
		fillBloom,
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
	if fr.OrderHash != "0x91b419e1cc29695dd4da477967c1b529eaad1591692566778eaf2d4baec3c593" {
		t.Errorf("Unexpected order hash, got '%v'", fr.OrderHash)
	}
	if fr.FilledTakerAssetAmount != "4" {
		t.Errorf("Unexpected filled amount, got '%v'", fr.FilledTakerAssetAmount)
	}
	if fr.Cancel != false {
		t.Errorf("Unexpected cancelled amount, got '%v'", fr.Cancel)
	}
	time.Sleep(1000 * time.Millisecond)
	fb, err := fillbloom.NewFillBloom(itemURL)
	if err != nil { t.Errorf(err.Error()) }
	fb.Initialize(
		mock.NewMockLogFilterer([]types.Log{*testLog}),
		0,
		[]common.Address{common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093")},
	)
	orderHash := common.HexToHash(fr.OrderHash)
	if !fb.Test(orderHash[:]) {
		t.Errorf("Expected to find orderHash '%#x' in bloom filter", orderHash[:])
	}
}
func TestNoAllowanceInBlock(t *testing.T) {
	directory := fmt.Sprintf("/tmp/test-%v", rand.Int())
	os.Mkdir(directory, 0755)
	itemURL := fmt.Sprintf("file://%v/test", directory)
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
	fillBloom, err := fillbloom.NewFillBloom(itemURL)
	if err != nil { t.Fatalf(err.Error()) }
	consumerChannel.AddConsumer(fill.NewFillBlockConsumer(
		common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093").Big(),
		mock.NewMockLogFilterer([]types.Log{}),
		destPublisher,
		fillBloom,
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
