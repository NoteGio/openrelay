package ingest_test

import (
	"bytes"
	"github.com/notegio/openrelay/ingest"
	"github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/types"
	poolModule "github.com/notegio/openrelay/pool"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	// "reflect"
	// "io/ioutil"
	// "fmt"
)

var feeBaseUnits = big.NewInt(1000000000000000000)

type mockBaseFee struct {
	fee *big.Int
}

func (bf *mockBaseFee) Get() (*big.Int, error) {
	return bf.fee, nil
}

func (bf *mockBaseFee) Set(fee *big.Int) error {
	bf.fee = fee
	return nil
}

func mockPoolDecorator(fn func(http.ResponseWriter, *http.Request, *poolModule.Pool)) func(http.ResponseWriter, *http.Request) {
	baseFee := &mockBaseFee{big.NewInt(0)}
	return func(w http.ResponseWriter, r *http.Request) {
		sender, _:= common.HexToAddress("0x0000000000000000000000000000000000000000")
		pool := &poolModule.Pool{SearchTerms: "", ID: []byte("default"), SenderAddresses: types.NetworkAddressMap{1: sender}, FeeTokenAddress: types.NetworkAddressMap{1: sender}}
		pool.SetBaseFee(baseFee)
		fn(w, r, pool)
	}
}

func TestFeeRecipientAndMakerProvided(t *testing.T) {
	publisher := TestPublisher{}
	handler := mockPoolDecorator(ingest.FeeHandler(&publisher, &TestAccountService{false, new(big.Int)}, &TestAffiliateService{new(big.Int), nil}, [20]byte{}, &TestExchangeLookup{1}))
	reader := TestReader{
		bytes.NewReader([]byte("{\"chainId\":1, \"maker\": \"0x0000000000000000000000000000000000000000\", \"feeRecipientAddress\": \"0000000000000000000000000000000000000000\", \"takerTokenAmount\": \"100\", \"makerTokenAmount\": \"100\"}")),
		nil,
	}
	request, _ := http.NewRequest("POST", "/v0.0/fees", reader)
	request.Header["Content-Type"] = []string{"application/json"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Expected error code 200, got '%v'", recorder.Code)
		t.Errorf("Body: '%v'", recorder.Body.String())
	}
	body := recorder.Body.String()
	if body != "{\"makerFee\":\"0\",\"takerFee\":\"0\",\"feeRecipientAddress\":\"0x0000000000000000000000000000000000000000\",\"senderAddress\":\"0x0000000000000000000000000000000000000000\",\"takerToSpecify\":\"0x0000000000000000000000000000000000000000\",\"makerFeeAssetData\":\"0xf47261b00000000000000000000000000000000000000000000000000000000000000000\",\"takerFeeAssetData\":\"0xf47261b00000000000000000000000000000000000000000000000000000000000000000\"}" {
		t.Errorf("Unexpected body: '%v'", body)
	}
}
func TestFeeRecipientAndMakerDefault(t *testing.T) {
	publisher := TestPublisher{}
	handler := mockPoolDecorator(ingest.FeeHandler(&publisher, &TestAccountService{false, new(big.Int)}, &TestAffiliateService{new(big.Int), nil}, [20]byte{}, &TestExchangeLookup{1}))
	reader := TestReader{
		bytes.NewReader([]byte("{\"takerTokenAmount\": \"100\", \"makerTokenAmount\": \"100\"}")),
		nil,
	}
	request, _ := http.NewRequest("POST", "/v0.0/fees", reader)
	request.Header["Content-Type"] = []string{"application/json"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Expected error code 200, got '%v'", recorder.Code)
		t.Errorf("Body: '%v'", recorder.Body.String())
	}
	body := recorder.Body.String()
	if body != "{\"makerFee\":\"0\",\"takerFee\":\"0\",\"feeRecipientAddress\":\"0x0000000000000000000000000000000000000000\",\"senderAddress\":\"0x0000000000000000000000000000000000000000\",\"takerToSpecify\":\"0x0000000000000000000000000000000000000000\",\"makerFeeAssetData\":\"0xf47261b00000000000000000000000000000000000000000000000000000000000000000\",\"takerFeeAssetData\":\"0xf47261b00000000000000000000000000000000000000000000000000000000000000000\"}" {
		t.Errorf("Unexpected body: '%v'", body)
	}
}
