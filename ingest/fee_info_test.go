package ingest_test

import (
	"github.com/notegio/openrelay/ingest"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	// "reflect"
	// "io/ioutil"
	// "fmt"
)

func TestFeeRecipientAndMakerProvided(t *testing.T) {
	publisher := TestPublisher{}
	handler := ingest.FeeHandler(&publisher, &TestAccountService{false, new(big.Int)}, &TestAffiliateService{new(big.Int), nil}, [20]byte{})
	reader := TestReader{
		[]byte("{\"maker\": \"0x0000000000000000000000000000000000000000\", \"feeRecipient\": \"0000000000000000000000000000000000000000\", \"takerTokenAmount\": \"100\", \"makerTokenAmount\": \"100\"}"),
		nil,
	}
	request, _ := http.NewRequest("POST", "/v0/fees", reader)
	request.Header["Content-Type"] = []string{"application/json"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Expected error code 200, got '%v'", recorder.Code)
		t.Errorf("Body: '%v'", recorder.Body.String())
	}
	body := recorder.Body.String()
	if body != "{\"makerFee\":\"0\",\"takerFee\":\"0\",\"feeRecipient\":\"0x0000000000000000000000000000000000000000\",\"takerToSpecify\":\"0x0000000000000000000000000000000000000000\"}" {
		t.Errorf("Unexpected body: '%v'", body)
	}
}
func TestFeeRecipientAndMakerDefault(t *testing.T) {
	publisher := TestPublisher{}
	handler := ingest.FeeHandler(&publisher, &TestAccountService{false, new(big.Int)}, &TestAffiliateService{new(big.Int), nil}, [20]byte{})
	reader := TestReader{
		[]byte("{\"takerTokenAmount\": \"100\", \"makerTokenAmount\": \"100\"}"),
		nil,
	}
	request, _ := http.NewRequest("POST", "/v0/fees", reader)
	request.Header["Content-Type"] = []string{"application/json"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Expected error code 200, got '%v'", recorder.Code)
		t.Errorf("Body: '%v'", recorder.Body.String())
	}
	body := recorder.Body.String()
	if body != "{\"makerFee\":\"0\",\"takerFee\":\"0\",\"feeRecipient\":\"0x0000000000000000000000000000000000000000\",\"takerToSpecify\":\"0x0000000000000000000000000000000000000000\"}" {
		t.Errorf("Unexpected body: '%v'", body)
	}
}
