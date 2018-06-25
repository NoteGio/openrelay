package types_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"github.com/notegio/openrelay/types"
	"io/ioutil"
	"reflect"
	"testing"
	"bytes"
	// "log"
)

func checkOrder(order *types.Order, t *testing.T) {
	if hex.EncodeToString(order.MakerAssetData[:]) != "0000000000000000000000000000000000000000" {
		t.Errorf("Unexpected MakerAssetData: %#x", order.MakerAssetData[:])
	}
	if hex.EncodeToString(order.Maker[:]) != "0000000000000000000000000000000000000000" {
		t.Errorf("Unexpected Maker: %#x", order.Maker[:])
	}
	if hex.EncodeToString(order.Taker[:]) != "0000000000000000000000000000000000000000" {
		t.Errorf("Unexpected Taker: %#x", order.Taker[:])
	}
	if hex.EncodeToString(order.FeeRecipient[:]) != "0000000000000000000000000000000000000000" {
		t.Errorf("Unexpected FeeRecipient: %#x", order.FeeRecipient[:])
	}
	if hex.EncodeToString(order.TakerAssetData[:]) != "0000000000000000000000000000000000000000" {
		t.Errorf("Unexpected TakerAssetData: %#x", order.TakerAssetData[:])
	}
	if hex.EncodeToString(order.MakerAssetAmount[:]) != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Errorf("Unexpected MakerAssetAmount: %#x", order.MakerAssetAmount[:])
	}
	if hex.EncodeToString(order.TakerAssetAmount[:]) != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Errorf("Unexpected MakerAssetAmount: %#x", order.TakerAssetAmount[:])
	}
	if hex.EncodeToString(order.MakerFee[:]) != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Errorf("Unexpected MakerFee: %#x", order.MakerFee[:])
	}
	if hex.EncodeToString(order.TakerFee[:]) != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Errorf("Unexpected TakerFee: %#x", order.TakerFee[:])
	}
	if hex.EncodeToString(order.ExpirationTimestampInSec[:]) != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Errorf("Unexpected ExpirationTimestampInSec: %#x", order.ExpirationTimestampInSec[:])
	}
	if hex.EncodeToString(order.ExchangeAddress[:]) != "b69e673309512a9d726f87304c6984054f87a93b" {
		t.Errorf("Unexpected ExchangeAddress: %#x", order.ExchangeAddress[:])
	}
	if hex.EncodeToString(order.Salt[:]) != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Errorf("Unexpected Salt: %#x", order.Salt[:])
	}
	if hex.EncodeToString(order.Hash()) != "367ad7730eb8b5feab8a9c9f47c6fcba77a2d4df125ee6a59cc26ac955710f7e" {
		t.Errorf("Hashes not equal %x", order.Hash())
	}
}

func TestOrderHash(t *testing.T) {
	order := &types.Order{}
	order.Initialize()
	exchangeAddressBytes, err := types.HexStringToBytes("b69e673309512a9d726f87304c6984054f87a93b")
	if err != nil {
		t.Errorf(err.Error())
	}
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	checkOrder(order, t)
}

func TestOrderToFromBytes(t *testing.T) {
	order := &types.Order{}
	order.Initialize()
	exchangeAddressBytes, err := types.HexStringToBytes("b69e673309512a9d726f87304c6984054f87a93b")
	if err != nil { t.Errorf(err.Error()) }
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	order2, err := types.OrderFromBytes(order.Bytes())
	if err != nil {
		t.Errorf(err.Error())
	}
	checkOrder(order2, t)
	if !bytes.Equal(order.Hash(), order2.Hash()) {
		t.Errorf("Unequal hashes: %#x != %#x", order.Hash(), order2.Hash())
	}

}

func value(valuer driver.Valuer) (interface{}, error) {
	return valuer.Value()
}

func scan(scanner sql.Scanner, data []byte) error {
	return scanner.Scan(data)
}

func TestValuerInterfaceAddress(t *testing.T) {
	order := &types.Order{}
	order.Initialize()
	exchangeAddressBytes, err := types.HexStringToBytes("b69e673309512a9d726f87304c6984054f87a93b")
	if err != nil { t.Errorf(err.Error()) }
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	if ExchangeAddress, _ := value(order.ExchangeAddress); !reflect.DeepEqual(ExchangeAddress, exchangeAddressBytes) {
		t.Errorf("Unexpected MakerAssetAmount")
	}
}

func TestValuerInterfaceUint256(t *testing.T) {
	order := &types.Order{}
	order.Initialize()
	exchangeAddressBytes, err := types.HexStringToBytes("b69e673309512a9d726f87304c6984054f87a93b")
	if err != nil { t.Errorf(err.Error()) }
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	if MakerAssetAmount, _ := value(order.MakerAssetAmount); !reflect.DeepEqual(MakerAssetAmount, make([]byte, 32)) {
		t.Errorf("Unexpected MakerAssetAmount")
	}
}

func TestScannerInterfaceUint256(t *testing.T) {
	order := &types.Order{}
	order.Initialize()
	sampleBytes := make([]byte, 32)
	sampleBytes[27] = 255
	scan(order.MakerAssetAmount, sampleBytes)
	if !reflect.DeepEqual(order.MakerAssetAmount[:], sampleBytes) {
		t.Errorf("Failed to load MakerAssetAmount: %v", hex.EncodeToString(order.MakerAssetAmount[:]))
	}
}
func TestJsonMarshal(t *testing.T) {
	order := &types.Order{}
	order.Initialize()
	exchangeAddressBytes, err := types.HexStringToBytes("b69e673309512a9d726f87304c6984054f87a93b")
	if err != nil { t.Errorf(err.Error()) }
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	sigBytes, err := types.HexStringToBytes("006bcc503876436ae6ebddecc16f95fdc74945ba85aa7debabdfa4a708a80b0272520d4f331a50396583db9a06bce884abc82219bfe180ef0093b0534786c996c203")
	if err != nil { t.Errorf(err.Error()) }
	order.Signature = append(order.Signature, sigBytes...)
	data, err := json.Marshal(order)
	if err != nil {
		t.Errorf("Got error marshalling: %v", err.Error())
		return
	}
	if string(data) != "{\"makerAddress\":\"0x0000000000000000000000000000000000000000\",\"takerAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetData\":\"0x0000000000000000000000000000000000000000\",\"takerAssetData\":\"0x0000000000000000000000000000000000000000\",\"feeRecipientAddress\":\"0x0000000000000000000000000000000000000000\",\"exchangeAddress\":\"0xb69e673309512a9d726f87304c6984054f87a93b\",\"senderAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetAmount\":\"0\",\"takerAssetAmount\":\"0\",\"makerFee\":\"0\",\"takerFee\":\"0\",\"expirationTimeSeconds\":\"0\",\"salt\":\"0\",\"signature\":\"0x006bcc503876436ae6ebddecc16f95fdc74945ba85aa7debabdfa4a708a80b0272520d4f331a50396583db9a06bce884abc82219bfe180ef0093b0534786c996c203\"}" {
		t.Errorf("Got unexpected JSON value: %v", string(data))
	}
}
func TestJsonMarshalSlice(t *testing.T) {
	order := &types.Order{}
	order.Initialize()
	exchangeAddressBytes, err := types.HexStringToBytes("b69e673309512a9d726f87304c6984054f87a93b")
	if err != nil { t.Errorf(err.Error()) }
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	sigBytes, err := types.HexStringToBytes("006bcc503876436ae6ebddecc16f95fdc74945ba85aa7debabdfa4a708a80b0272520d4f331a50396583db9a06bce884abc82219bfe180ef0093b0534786c996c203")
	if err != nil { t.Errorf(err.Error()) }
	order.Signature = append(order.Signature, sigBytes...)
	data, err := json.Marshal([]types.Order{*order})
	if err != nil {
		t.Errorf("Got error marshalling: %v", err.Error())
		return
	}
	if string(data) != "[{\"makerAddress\":\"0x0000000000000000000000000000000000000000\",\"takerAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetData\":\"0x0000000000000000000000000000000000000000\",\"takerAssetData\":\"0x0000000000000000000000000000000000000000\",\"feeRecipientAddress\":\"0x0000000000000000000000000000000000000000\",\"exchangeAddress\":\"0xb69e673309512a9d726f87304c6984054f87a93b\",\"senderAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetAmount\":\"0\",\"takerAssetAmount\":\"0\",\"makerFee\":\"0\",\"takerFee\":\"0\",\"expirationTimeSeconds\":\"0\",\"salt\":\"0\",\"signature\":\"0x006bcc503876436ae6ebddecc16f95fdc74945ba85aa7debabdfa4a708a80b0272520d4f331a50396583db9a06bce884abc82219bfe180ef0093b0534786c996c203\"}]" {
		t.Errorf("Got unexpected JSON value: %v", string(data))
	}
}

func TestJsonUnmarshal(t *testing.T) {
	newOrder := types.Order{}
	if orderData, err := ioutil.ReadFile("../formatted_transaction.json"); err == nil {
		if err := json.Unmarshal(orderData, &newOrder); err != nil {
			t.Fatalf(err.Error())
		}
	}
	if hex.EncodeToString(newOrder.MakerAssetData[:]) != "1dad4783cf3fe3085c1426157ab175a6119a04ba01" {
		t.Errorf("Unexpected MakerAssetData: %#x", newOrder.MakerAssetData[:])
	}
	if hex.EncodeToString(newOrder.Maker[:]) != "627306090abab3a6e1400e9345bc60c78a8bef57" {
		t.Errorf("Unexpected Maker: %#x", newOrder.Maker[:])
	}
	if hex.EncodeToString(newOrder.Taker[:]) != "0000000000000000000000000000000000000000" {
		t.Errorf("Unexpected Taker: %#x", newOrder.Taker[:])
	}
	if hex.EncodeToString(newOrder.FeeRecipient[:]) != "0000000000000000000000000000000000000000" {
		t.Errorf("Unexpected FeeRecipient: %#x", newOrder.FeeRecipient[:])
	}
	if hex.EncodeToString(newOrder.TakerAssetData[:]) != "05d090b51c40b020eab3bfcb6a2dff130df22e9c01" {
		t.Errorf("Unexpected TakerAssetData: %#x", newOrder.TakerAssetData[:])
	}
	if hex.EncodeToString(newOrder.MakerAssetAmount[:]) != "000000000000000000000000000000000000000000000002b5e3af16b1880000" {
		t.Errorf("Unexpected MakerAssetAmount: %#x", newOrder.MakerAssetAmount[:])
	}
	if hex.EncodeToString(newOrder.TakerAssetAmount[:]) != "0000000000000000000000000000000000000000000000000de0b6b3a7640000" {
		t.Errorf("Unexpected MakerAssetAmount: %#x", newOrder.TakerAssetAmount[:])
	}
	if hex.EncodeToString(newOrder.MakerFee[:]) != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Errorf("Unexpected MakerFee: %#x", newOrder.MakerFee[:])
	}
	if hex.EncodeToString(newOrder.TakerFee[:]) != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Errorf("Unexpected TakerFee: %#x", newOrder.TakerFee[:])
	}
	if hex.EncodeToString(newOrder.ExpirationTimestampInSec[:]) != "0000000000000000000000000000000000000000000000000000000159938ac4" {
		t.Errorf("Unexpected ExpirationTimestampInSec: %#x", newOrder.ExpirationTimestampInSec[:])
	}
	if hex.EncodeToString(newOrder.ExchangeAddress[:]) != "90fe2af704b34e0224bf2299c838e04d4dcf1364" {
		t.Errorf("Unexpected ExchangeAddress: %#x", newOrder.ExchangeAddress[:])
	}
	if hex.EncodeToString(newOrder.Salt[:]) != "000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b" {
		t.Errorf("Unexpected Salt: %#x", newOrder.Salt[:])
	}
	if hex.EncodeToString(newOrder.Hash()) != "693c35d2cca7d6011b62ecb01aef6763241cec2b2686fa7f3281843b1aeb1c6b" {
		t.Errorf("Hashes not equal %x", newOrder.Hash())
	}

	if !newOrder.Signature.Verify(newOrder.Maker, newOrder.Hash()) {
		t.Errorf("Failed to verify order with signature: %#x", newOrder.Signature)
	}
}
