package types_test

import (
	// "database/sql"
	// "database/sql/driver"
	"encoding/hex"
	// "encoding/json"
	"github.com/notegio/openrelay/types"
	// "io/ioutil"
	// "reflect"
	"testing"
	"bytes"
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
	// if order.Signature.V != byte(27) {
	// 	t.Errorf("Unexpected sig.v %v", order.Signature.V)
	// }
	// if hex.EncodeToString(order.Signature.R[:]) != "37adbc51c87a2f4c8c40c25fab5a73c65d078322f1db5739ee6fd49f18ce4463" {
	// 	t.Errorf("Unexpected sig.r: %x", order.Signature.R[:])
	// }
	// if hex.EncodeToString(order.Signature.S[:]) != "7382de9b4cf7ceaf602f221132c9ddf41b83fb9666839022703da852d4ed88af" {
	// 	t.Errorf("Unexpected sig.s: %x", order.Signature.S[:])
	// }
	if hex.EncodeToString(order.Hash()) != "367ad7730eb8b5feab8a9c9f47c6fcba77a2d4df125ee6a59cc26ac955710f7e" {
		t.Errorf("Hashes not equal %x", order.Hash())
	}
	// if !order.Signature.Verify(order.Maker) {
	// 	t.Errorf("Signature not valid")
	// }
	// testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b37adbc51c87a2f4c8c40c25fab5a73c65d078322f1db5739ee6fd49f18ce44637382de9b4cf7ceaf602f221132c9ddf41b83fb9666839022703da852d4ed88af")
	// calculatedBytes := order.Bytes()
	// if !reflect.DeepEqual(calculatedBytes[:377], testOrderBytes[:377]) {
	// 	t.Errorf("Unexpected byte stream")
	// }
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
	if err != nil {
		t.Errorf(err.Error())
	}
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

//
// func TestByteDeserialize(t *testing.T) {
// 	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b37adbc51c87a2f4c8c40c25fab5a73c65d078322f1db5739ee6fd49f18ce44637382de9b4cf7ceaf602f221132c9ddf41b83fb9666839022703da852d4ed88af")
// 	var testOrderByteArray [441]byte
// 	copy(testOrderByteArray[:], testOrderBytes[:])
// 	newOrder := types.OrderFromBytes(testOrderByteArray)
// 	checkOrder(newOrder, t)
// }
// func TestJSONDeserialize(t *testing.T) {
// 	newOrder := types.Order{}
// 	if orderData, err := ioutil.ReadFile("../formatted_transaction.json"); err == nil {
// 		if err := json.Unmarshal(orderData, &newOrder); err != nil {
// 			t.Errorf(err.Error())
// 			return
// 		}
// 	}
// 	checkOrder(&newOrder, t)
// }
//
// func value(valuer driver.Valuer) (interface{}, error) {
// 	return valuer.Value()
// }
//
// func scan(scanner sql.Scanner, data []byte) error {
// 	return scanner.Scan(data)
// }
//
// func TestValuerInterfaceAddress(t *testing.T) {
// 	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
// 	var testOrderByteArray [441]byte
// 	copy(testOrderByteArray[:], testOrderBytes[:])
// 	newOrder := types.OrderFromBytes(testOrderByteArray)
// 	if ExchangeValue, _ := value(newOrder.ExchangeAddress); !reflect.DeepEqual(ExchangeValue, testOrderBytes[:20]) {
// 		t.Errorf("Unexpected ExchangeAddress")
// 	}
// }
//
// func TestScannerInterfaceAddress(t *testing.T) {
// 	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
// 	order := types.Order{}
// 	order.Initialize()
// 	scan(order.ExchangeAddress, testOrderBytes[:20])
// 	if !reflect.DeepEqual(order.ExchangeAddress[:], testOrderBytes[:20]) {
// 		t.Errorf("Failed to load exchange address: %v", hex.EncodeToString(order.ExchangeAddress[:]))
// 	}
// }
//
// func TestValuerInterfaceUint256(t *testing.T) {
// 	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
// 	var testOrderByteArray [441]byte
// 	copy(testOrderByteArray[:], testOrderBytes[:])
// 	newOrder := types.OrderFromBytes(testOrderByteArray)
// 	if MakerTokenAmount, _ := value(newOrder.MakerTokenAmount); !reflect.DeepEqual(MakerTokenAmount, testOrderBytes[120:152]) {
// 		t.Errorf("Unexpected MakerTokenAmount")
// 	}
// }
//
// func TestScannerInterfaceUint256(t *testing.T) {
// 	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
// 	order := types.Order{}
// 	order.Initialize()
// 	scan(order.MakerTokenAmount, testOrderBytes[120:152])
// 	if !reflect.DeepEqual(order.MakerTokenAmount[:], testOrderBytes[120:152]) {
// 		t.Errorf("Failed to load MakerTokenAmount: %v", hex.EncodeToString(order.MakerTokenAmount[:]))
// 	}
// }
//
// func TestValuerInterfaceSignature(t *testing.T) {
// 	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
// 	var testOrderByteArray [441]byte
// 	copy(testOrderByteArray[:], testOrderBytes[:])
// 	newOrder := types.OrderFromBytes(testOrderByteArray)
// 	sigValue, _ := value(newOrder.Signature)
// 	sigBytes := sigValue.([]byte)
//
// 	if !reflect.DeepEqual(sigBytes[:32], testOrderBytes[313:345]) {
// 		t.Errorf("Unexpected Sig R")
// 	}
// 	if !reflect.DeepEqual(sigBytes[32:64], testOrderBytes[345:377]) {
// 		t.Errorf("Unexpected Sig S")
// 	}
// 	if sigBytes[64] != byte(int(testOrderBytes[312])-27) {
// 		t.Errorf("Unexpected Sig V")
// 	}
//
// }
//
// func TestScannerInterfaceSignature(t *testing.T) {
// 	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
// 	var testOrderByteArray [441]byte
// 	copy(testOrderByteArray[:], testOrderBytes[:])
//
// 	sigBytes := make([]byte, 65)
// 	copy(sigBytes[:64], testOrderBytes[313:377])
// 	sigBytes[64] = byte(int(testOrderBytes[312]) - 27)
// 	signature := &types.Signature{}
// 	scan(signature, sigBytes)
// 	if !reflect.DeepEqual(signature.R[:], sigBytes[:32]) {
// 		t.Errorf("Unexpected Sig R")
// 	}
// 	if !reflect.DeepEqual(signature.S[:], sigBytes[32:64]) {
// 		t.Errorf("Unexpected Sig S")
// 	}
// 	if signature.V != testOrderBytes[312] {
// 		t.Errorf("Unexpected Sig V")
// 	}
// }
//
// func TestJsonMarshal(t *testing.T) {
// 	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
// 	var testOrderByteArray [441]byte
// 	copy(testOrderByteArray[:], testOrderBytes[:])
// 	newOrder := types.OrderFromBytes(testOrderByteArray)
// 	data, err := json.Marshal(newOrder)
// 	if err != nil {
// 		t.Errorf("Got error marshalling: %v", err.Error())
// 		return
// 	}
// 	if string(data) != "{\"maker\":\"0x324454186bb728a3ea55750e0618ff1b18ce6cf8\",\"taker\":\"0x0000000000000000000000000000000000000000\",\"makerTokenAddress\":\"0x1dad4783cf3fe3085c1426157ab175a6119a04ba\",\"takerTokenAddress\":\"0x05d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"feeRecipient\":\"0x0000000000000000000000000000000000000000\",\"exchangeContractAddress\":\"0x90fe2af704b34e0224bf2299c838e04d4dcf1364\",\"makerTokenAmount\":\"50000000000000000000\",\"takerTokenAmount\":\"1000000000000000000\",\"makerFee\":\"0\",\"takerFee\":\"0\",\"expirationUnixTimestampSec\":\"5797808836\",\"salt\":\"11065671350908846865864045738088581419204014210814002044381812654087807531\",\"ecSignature\":{\"v\":27,\"r\":\"0x021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e\",\"s\":\"0x12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1\"}}" {
// 		t.Errorf("Got unexpected JSON value: %v", string(data))
// 	}
// }
//
// func TestJsonMarshalSlice(t *testing.T) {
// 	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
// 	var testOrderByteArray [441]byte
// 	copy(testOrderByteArray[:], testOrderBytes[:])
// 	newOrder := types.OrderFromBytes(testOrderByteArray)
// 	orderList := []types.Order{*newOrder}
// 	data, err := json.Marshal(orderList)
// 	if err != nil {
// 		t.Errorf("Got error marshalling: %v", err.Error())
// 		return
// 	}
// 	if string(data) != "[{\"maker\":\"0x324454186bb728a3ea55750e0618ff1b18ce6cf8\",\"taker\":\"0x0000000000000000000000000000000000000000\",\"makerTokenAddress\":\"0x1dad4783cf3fe3085c1426157ab175a6119a04ba\",\"takerTokenAddress\":\"0x05d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"feeRecipient\":\"0x0000000000000000000000000000000000000000\",\"exchangeContractAddress\":\"0x90fe2af704b34e0224bf2299c838e04d4dcf1364\",\"makerTokenAmount\":\"50000000000000000000\",\"takerTokenAmount\":\"1000000000000000000\",\"makerFee\":\"0\",\"takerFee\":\"0\",\"expirationUnixTimestampSec\":\"5797808836\",\"salt\":\"11065671350908846865864045738088581419204014210814002044381812654087807531\",\"ecSignature\":{\"v\":27,\"r\":\"0x021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e\",\"s\":\"0x12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1\"}}]" {
// 		t.Errorf("Got unexpected JSON value: %v", string(data))
// 	}
// }
