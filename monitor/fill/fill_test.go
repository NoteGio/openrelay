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
	fillTopic.SetString("6869791f0a34781b29882982cc39e882768cf2c96995c2a110c577c53bc932d5", 16)
	topics := []common.Hash{
		common.BigToHash(fillTopic),
		common.HexToHash("0x5409ed021d9299bf6814279a6a1411a7e866a631"),
		common.HexToHash("0x00"),
		common.HexToHash("0x91b419e1cc29695dd4da477967c1b529eaad1591692566778eaf2d4baec3c593"),
	}
	data, _ := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000220000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000042e2e925a5febbbf9da7e7a632e877fe5c000000000000000000000000591e9f22e2e925a5febbbf9da7e7a632e877fe5c000000000000000000000000000000000000000000000001655a6092164f0000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000672bb9e46b9f00000000000000000000000000000000000000000000000000000000000000024f47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000006b175474e89094c44da98b954eedeac495271d0f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
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
	addr, _ := big.NewInt(0).SetString("0x12459c951127e0c374ff9105dda097662a027093", 0)
	consumerChannel.AddConsumer(fill.NewFillBlockConsumer(
		addr,
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
	addr, _ := big.NewInt(0).SetString("0x12459c951127e0c374ff9105dda097662a027093", 0)
	consumerChannel.AddConsumer(fill.NewFillBlockConsumer(
		addr,
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
