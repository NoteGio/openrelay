package types_test

import (
	"encoding/hex"
	"encoding/json"
	"github.com/notegio/openrelay/types"
	"io/ioutil"
	"reflect"
	"testing"
)

func checkOrder(order *types.Order, t *testing.T) {
	if hex.EncodeToString(order.MakerToken[:]) != "1dad4783cf3fe3085c1426157ab175a6119a04ba" {
		t.Errorf("Unexpected MakerToken")
	}
	if hex.EncodeToString(order.Maker[:]) != "324454186bb728a3ea55750e0618ff1b18ce6cf8" {
		t.Errorf("Unexpected Maker")
	}
	if hex.EncodeToString(order.Taker[:]) != "0000000000000000000000000000000000000000" {
		t.Errorf("Unexpected Taker")
	}
	if hex.EncodeToString(order.FeeRecipient[:]) != "0000000000000000000000000000000000000000" {
		t.Errorf("Unexpected FeeRecipient")
	}
	if hex.EncodeToString(order.TakerToken[:]) != "05d090b51c40b020eab3bfcb6a2dff130df22e9c" {
		t.Errorf("Unexpected TakerToken")
	}
	if hex.EncodeToString(order.ExchangeAddress[:]) != "90fe2af704b34e0224bf2299c838e04d4dcf1364" {
		t.Errorf("Unexpected TakerToken")
	}
	if hex.EncodeToString(order.MakerTokenAmount[:]) != "000000000000000000000000000000000000000000000002b5e3af16b1880000" {
		t.Errorf("Unexpected MakerTokenAmount")
	}
	if hex.EncodeToString(order.TakerTokenAmount[:]) != "0000000000000000000000000000000000000000000000000de0b6b3a7640000" {
		t.Errorf("Unexpected MakerTokenAmount")
	}
	if hex.EncodeToString(order.MakerFee[:]) != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Errorf("Unexpected MakerFee")
	}
	if hex.EncodeToString(order.TakerFee[:]) != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Errorf("Unexpected TakerFee")
	}
	if hex.EncodeToString(order.ExpirationTimestampInSec[:]) != "0000000000000000000000000000000000000000000000000000000059938ac4" {
		t.Errorf("Unexpected ExpirationTimestampInSec")
	}
	if hex.EncodeToString(order.Salt[:]) != "000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b" {
		t.Errorf("Unexpected Salt")
	}
	if order.Signature.V != byte(0) {
		t.Errorf("Unexpected sig.v %v", order.Signature.V)
	}
	if hex.EncodeToString(order.Signature.R[:]) != "021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e" {
		t.Errorf("Unexpected sig.r")
	}
	if hex.EncodeToString(order.Signature.S[:]) != "12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1" {
		t.Errorf("Unexpected sig.s")
	}
	if hex.EncodeToString(order.Hash()) != "731319211689ccf0327911a0126b0af0854570c1b6cdfeb837b0127e29fe9fd5" {
		t.Errorf("Hashes not equal")
	}
	if !order.Signature.Verify(order.Maker) {
		t.Errorf("Signature not valid")
	}
	testOrderBytes, _ := hex.DecodeString("324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c000000000000000000000000000000000000000090fe2af704b34e0224bf2299c838e04d4dcf1364000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000059938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b00021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	calculatedBytes := order.Bytes()
	if !reflect.DeepEqual(calculatedBytes[:], testOrderBytes[:]) {
		t.Errorf("Unexpected byte stream")
	}
}

func TestByteDeserialize(t *testing.T) {
	testOrderBytes, _ := hex.DecodeString("324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c000000000000000000000000000000000000000090fe2af704b34e0224bf2299c838e04d4dcf1364000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000059938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b00021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	var testOrderByteArray [377]byte
	copy(testOrderByteArray[:], testOrderBytes[:])
	newOrder := types.OrderFromBytes(testOrderByteArray)
	checkOrder(newOrder, t)
}
func TestJSONDeserialize(t *testing.T) {
	newOrder := types.Order{}
	if orderData, err := ioutil.ReadFile("../formatted_transaction.json"); err == nil {
		if err := json.Unmarshal(orderData, &newOrder); err != nil {
			t.Errorf(err.Error())
			return
		}
	}
	checkOrder(&newOrder, t)
}
