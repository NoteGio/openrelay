package search_test

import (
	"testing"
	"encoding/hex"
	"github.com/notegio/openrelay/search"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/blockhash"
	"github.com/notegio/openrelay/types"
	dbModule "github.com/notegio/openrelay/db"
	"net/http"
	"net/http/httptest"
	"reflect"
	"time"
)

func getTestOrderBytes() [441]byte {
	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000059938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	var testOrderByteArray [441]byte
	copy(testOrderByteArray[:], testOrderBytes[:])
	return testOrderByteArray
}

func sampleOrder() *dbModule.Order {
	order := &types.Order{}
	order.FromBytes(getTestOrderBytes())
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *order
	return dbOrder
}

func TestFormatResponseJson(t *testing.T) {
	order := sampleOrder()
	orders := []dbModule.Order{*order, *order}
	response, contentType, err := search.FormatResponse(orders, "application/json")
	if err != nil {
		t.Errorf("Error getting formatted response: %v", err.Error())
	}
	if contentType != "application/json" {
		t.Errorf("Expected content type application/json, got '%v'", contentType)
	}
	if string(response) != "[{\"maker\":\"0x324454186bb728a3ea55750e0618ff1b18ce6cf8\",\"taker\":\"0x0000000000000000000000000000000000000000\",\"makerTokenAddress\":\"0x1dad4783cf3fe3085c1426157ab175a6119a04ba\",\"takerTokenAddress\":\"0x05d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"feeRecipient\":\"0x0000000000000000000000000000000000000000\",\"exchangeContractAddress\":\"0x90fe2af704b34e0224bf2299c838e04d4dcf1364\",\"makerTokenAmount\":\"50000000000000000000\",\"takerTokenAmount\":\"1000000000000000000\",\"makerFee\":\"0\",\"takerFee\":\"0\",\"expirationUnixTimestampSec\":\"1502841540\",\"salt\":\"11065671350908846865864045738088581419204014210814002044381812654087807531\",\"ecSignature\":{\"v\":\"27\",\"r\":\"0x021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e\",\"s\":\"0x12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1\"},\"takerTokenAmountFilled\":\"0\",\"takerTokenAmountCancelled\":\"0\"},{\"maker\":\"0x324454186bb728a3ea55750e0618ff1b18ce6cf8\",\"taker\":\"0x0000000000000000000000000000000000000000\",\"makerTokenAddress\":\"0x1dad4783cf3fe3085c1426157ab175a6119a04ba\",\"takerTokenAddress\":\"0x05d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"feeRecipient\":\"0x0000000000000000000000000000000000000000\",\"exchangeContractAddress\":\"0x90fe2af704b34e0224bf2299c838e04d4dcf1364\",\"makerTokenAmount\":\"50000000000000000000\",\"takerTokenAmount\":\"1000000000000000000\",\"makerFee\":\"0\",\"takerFee\":\"0\",\"expirationUnixTimestampSec\":\"1502841540\",\"salt\":\"11065671350908846865864045738088581419204014210814002044381812654087807531\",\"ecSignature\":{\"v\":\"27\",\"r\":\"0x021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e\",\"s\":\"0x12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1\"},\"takerTokenAmountFilled\":\"0\",\"takerTokenAmountCancelled\":\"0\"}]" {
		t.Errorf("Got '%v'", string(response))
	}
}
func TestFormatResponseBin(t *testing.T) {
	order := sampleOrder()
	orders := []dbModule.Order{*order, *order}
	response, contentType, err := search.FormatResponse(orders, "application/octet-stream")
	if err != nil {
		t.Errorf("Error getting formatted response: %v", err.Error())
	}
	if contentType != "application/octet-stream" {
		t.Errorf("Expected content type application/json, got '%v'", contentType)
	}
	orderBytes := order.Bytes()
	orderValue := []byte{}
	orderValue = append(orderValue, orderBytes[:]...)
	orderValue = append(orderValue, orderBytes[:]...)
	if !reflect.DeepEqual(response, orderValue) {
		t.Errorf("Got '%#x'", response)
	}
}

func TestBlockhashRedirect(t *testing.T) {
	publisher, consumerChannel := channels.MockChannel()
	blockHash := blockhash.NewChanneledBlockHash(consumerChannel)
	publisher.Publish("hashValue")
	time.Sleep(300 * time.Millisecond)
	handler := search.Handler(nil, blockHash)
	request, _ := http.NewRequest("GET", "/v0/orders?makertoken=0x324454186bb728a3ea55750e0618ff1b18ce6cf8", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 307 {
		t.Errorf("Did not redirect")
	}
	if location := recorder.Header().Get("Location"); location != "/v0/orders?blockhash=hashValue&makertoken=0x324454186bb728a3ea55750e0618ff1b18ce6cf8" {
		t.Errorf("Expected orderHash to be added, got '%v'", location)
	}
}
