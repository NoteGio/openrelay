package ingest_test

import (
	"encoding/hex"
	accountsModule "github.com/notegio/openrelay/accounts"
	affiliatesModule "github.com/notegio/openrelay/affiliates"
	"github.com/notegio/openrelay/ingest"
	"github.com/notegio/openrelay/types"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	// "reflect"
	"errors"
	// "io/ioutil"
	// "fmt"
)

type TestPublisher struct {
	channel  string
	messages []string
}

func (pub *TestPublisher) Publish(message string) bool {
	pub.messages = append(pub.messages, message)
	return true
}

type TestAccount struct {
	blacklist bool
	discount  *big.Int
}

func (acct *TestAccount) Blacklisted() bool {
	return acct.blacklist
}
func (acct *TestAccount) Discount() *big.Int {
	return acct.discount
}

type TestAffiliate struct {
	fee *big.Int
}

func (affiliate *TestAffiliate) Fee() *big.Int {
	return affiliate.fee
}

type TestAffiliateService struct {
	fee *big.Int
	err error
}

func (service *TestAffiliateService) Get(address *types.Address) (affiliatesModule.Affiliate, error) {
	if service.err == nil {
		return &TestAffiliate{service.fee}, nil
	}
	return nil, service.err
}

// Set must be provided to satisfy the interface, but we don't need it for these tests.
func (service *TestAffiliateService) Set(address *types.Address, affiliate affiliatesModule.Affiliate) error {
	return nil
}

type TestAccountService struct {
	blacklist bool
	discount  *big.Int
}

type TestReader struct {
	body []byte
	err  error
}

func (reader TestReader) Read(p []byte) (n int, err error) {
	copy(p[:], reader.body[:])
	return len(reader.body), reader.err
}

// Get makes up an account deterministically based on the provided address
func (service *TestAccountService) Get(address *types.Address) accountsModule.Account {
	account := &TestAccount{
		service.blacklist,
		service.discount,
	}
	return account
}

// Set must be provided to satisfy the interface, but we don't need it for these tests.
func (service *TestAccountService) Set(address *types.Address, account accountsModule.Account) error {
	return nil
}

func TestTooLongBytes(t *testing.T) {
	publisher := TestPublisher{}
	handler := ingest.Handler(&publisher, &TestAccountService{}, &TestAffiliateService{})
	reader := TestReader{
		[]byte("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		nil,
	}
	request, _ := http.NewRequest("POST", "/", reader)
	request.Header["Content-Type"] = []string{"application/octet-stream"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 400 {
		t.Errorf("Expected error code 400, got '%v'", recorder.Code)
	}
	if recorder.HeaderMap["Content-Type"][0] != "application/json" {
		t.Errorf("Got unexpected content type", recorder.HeaderMap["Content-Type"])
	}
	body := recorder.Body.String()
	if body != "{\"code\":100,\"reason\":\"Orders should be exactly 377 bytes\"}" {
		t.Errorf("Got unexpected body: '%v' - %v", body, len(body))
	}
	if len(publisher.messages) != 0 {
		t.Errorf("Unexpected message count '%v'", len(publisher.messages))
	}
}
func TestBadRead(t *testing.T) {
	publisher := TestPublisher{}
	handler := ingest.Handler(&publisher, &TestAccountService{}, &TestAffiliateService{})
	reader := TestReader{
		[]byte("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		errors.New("Fail!"),
	}
	request, _ := http.NewRequest("POST", "/", reader)
	request.Header["Content-Type"] = []string{"application/octet-stream"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 500 {
		t.Errorf("Expected error code 500, got '%v'", recorder.Code)
	}
	if recorder.HeaderMap["Content-Type"][0] != "application/json" {
		t.Errorf("Got unexpected content type", recorder.HeaderMap["Content-Type"])
	}
	body := recorder.Body.String()
	if body != "{\"code\":100,\"reason\":\"Error reading content\"}" {
		t.Errorf("Got unexpected body: '%v' - %v", body, len(body))
	}
	if len(publisher.messages) != 0 {
		t.Errorf("Unexpected message count '%v'", len(publisher.messages))
	}
}
func TestBadJSON(t *testing.T) {
	publisher := TestPublisher{}
	handler := ingest.Handler(&publisher, &TestAccountService{}, &TestAffiliateService{})
	reader := TestReader{
		[]byte("bad json"),
		nil,
	}
	request, _ := http.NewRequest("POST", "/", reader)
	request.Header["Content-Type"] = []string{"application/json"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 400 {
		t.Errorf("Expected error code 400, got '%v'", recorder.Code)
	}
	if contentType, ok := recorder.HeaderMap["Content-Type"]; !ok || contentType[0] != "application/json" {
		t.Errorf("Got unexpected content type %v", contentType)
	}
	body := recorder.Body.String()
	if body != "{\"code\":101,\"reason\":\"Malformed JSON\"}" {
		t.Errorf("Got unexpected body: '%v' - %v", body, len(body))
	}
	if len(publisher.messages) != 0 {
		t.Errorf("Unexpected message count '%v'", len(publisher.messages))
	}
}
func TestJSONBadRead(t *testing.T) {
	publisher := TestPublisher{}
	handler := ingest.Handler(&publisher, &TestAccountService{}, &TestAffiliateService{})
	reader := TestReader{
		[]byte("bad json"),
		errors.New("Sample Error"),
	}
	request, _ := http.NewRequest("POST", "/", reader)
	request.Header["Content-Type"] = []string{"application/json"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 500 {
		t.Errorf("Expected error code 500, got '%v'", recorder.Code)
	}
	if contentType, ok := recorder.HeaderMap["Content-Type"]; !ok || contentType[0] != "application/json" {
		t.Errorf("Got unexpected content type %v", contentType)
	}
	body := recorder.Body.String()
	if body != "{\"code\":100,\"reason\":\"Error reading content\"}" {
		t.Errorf("Got unexpected body: '%v' - %v", body, len(body))
	}
	if len(publisher.messages) != 0 {
		t.Errorf("Unexpected message count '%v'", len(publisher.messages))
	}
}
func TestNoContentType(t *testing.T) {
	publisher := TestPublisher{}
	handler := ingest.Handler(&publisher, &TestAccountService{}, &TestAffiliateService{})
	reader := TestReader{
		[]byte(""),
		nil,
	}
	request, _ := http.NewRequest("POST", "/", reader)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 415 {
		t.Errorf("Expected error code 400, got '%v'", recorder.Code)
	}
	if contentType, ok := recorder.HeaderMap["Content-Type"]; !ok || contentType[0] != "application/json" {
		t.Errorf("Got unexpected content type %v", contentType)
	}
	body := recorder.Body.String()
	if body != "{\"code\":100,\"reason\":\"Unsupported content-type\"}" {
		t.Errorf("Got unexpected body: '%v' - %v", body, len(body))
	}
	if len(publisher.messages) != 0 {
		t.Errorf("Unexpected message count '%v'", len(publisher.messages))
	}
}
func TestBadSignature(t *testing.T) {
	publisher := TestPublisher{}
	handler := ingest.Handler(&publisher, &TestAccountService{}, &TestAffiliateService{})
	data, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000059938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cf3139501007f1")
	reader := TestReader{
		data,
		nil,
	}
	request, _ := http.NewRequest("POST", "/", reader)
	request.Header["Content-Type"] = []string{"application/octet-stream"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 400 {
		t.Errorf("Expected error code 400, got '%v'", recorder.Code)
	}
	if contentType, ok := recorder.HeaderMap["Content-Type"]; !ok || contentType[0] != "application/json" {
		t.Errorf("Got unexpected content type %v", contentType)
	}
	body := recorder.Body.String()
	if body != "{\"code\":100,\"reason\":\"Validation Failed\",\"validationErrors\":[{\"field\":\"ecSignature\",\"code\":1005,\"reason\":\"Signature validation failed\"}]}" {
		t.Errorf("Got unexpected body: '%v' - %v", body, len(body))
	}
	if len(publisher.messages) != 0 {
		t.Errorf("Unexpected message count '%v'", len(publisher.messages))
	}
}
func TestInsufficientFee(t *testing.T) {
	publisher := TestPublisher{}
	fee := new(big.Int)
	fee.SetInt64(1000)
	handler := ingest.Handler(&publisher, &TestAccountService{false, new(big.Int)}, &TestAffiliateService{fee, nil})
	data, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000059938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	reader := TestReader{
		data,
		nil,
	}
	request, _ := http.NewRequest("POST", "/", reader)
	request.Header["Content-Type"] = []string{"application/octet-stream"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 402 {
		t.Errorf("Expected error code 402, got '%v'", recorder.Code)
	}
	if contentType, ok := recorder.HeaderMap["Content-Type"]; !ok || contentType[0] != "application/json" {
		t.Errorf("Got unexpected content type %v", contentType)
	}
	body := recorder.Body.String()
	if body != "{\"code\":100,\"reason\":\"Validation Failed\",\"validationErrors\":[{\"field\":\"makerFee\",\"code\":1004,\"reason\":\"Total fee must be at least: 1000\"},{\"field\":\"takerFee\",\"code\":1004,\"reason\":\"Total fee must be at least: 1000\"}]}" {
		t.Errorf("Got unexpected body: '%v' - %v", body, len(body))
	}
	if len(publisher.messages) != 0 {
		t.Errorf("Unexpected message count '%v'", len(publisher.messages))
	}
}
func TestBlacklisted(t *testing.T) {
	publisher := TestPublisher{}
	handler := ingest.Handler(&publisher, &TestAccountService{true, new(big.Int)}, &TestAffiliateService{new(big.Int), nil})
	data, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000059938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	reader := TestReader{
		data,
		nil,
	}
	request, _ := http.NewRequest("POST", "/", reader)
	request.Header["Content-Type"] = []string{"application/octet-stream"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 202 {
		t.Errorf("Expected code 202, got '%v'", recorder.Code)
	}
	if len(publisher.messages) != 0 {
		t.Errorf("Unexpected message count '%v'", len(publisher.messages))
	}
}
func TestNotFeeRecipient(t *testing.T) {
	publisher := TestPublisher{}
	handler := ingest.Handler(&publisher, &TestAccountService{true, new(big.Int)}, &TestAffiliateService{nil, errors.New("Fee Recipient must be an authorized address")})
	data, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000059938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	reader := TestReader{
		data,
		nil,
	}
	request, _ := http.NewRequest("POST", "/", reader)
	request.Header["Content-Type"] = []string{"application/octet-stream"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 402 {
		t.Errorf("Expected error code 402, got '%v'", recorder.Code)
	}
	body := recorder.Body.String()
	if body != "{\"code\":100,\"reason\":\"Validation Failed\",\"validationErrors\":[{\"field\":\"feeRecipient\",\"code\":1002,\"reason\":\"Invalid fee recipient\"}]}" {
		t.Errorf("Got unexpected body: '%v' - %v", body, len(body))
	}
	if len(publisher.messages) != 0 {
		t.Errorf("Unexpected message count '%v'", len(publisher.messages))
	}
}
func TestValid(t *testing.T) {
	publisher := TestPublisher{}
	handler := ingest.Handler(&publisher, &TestAccountService{false, new(big.Int)}, &TestAffiliateService{new(big.Int), nil})
	data, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000059938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	reader := TestReader{
		data[:377],
		nil,
	}
	request, _ := http.NewRequest("POST", "/", reader)
	request.Header["Content-Type"] = []string{"application/octet-stream"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 202 {
		t.Errorf("Expected error code 202, got '%v'", recorder.Code)
		t.Errorf("Body: '%v'", recorder.Body.String())
	}
	if len(publisher.messages) != 1 {
		t.Errorf("Unexpected message count '%v'", len(publisher.messages))
		return
	}
	if publisher.messages[0] != string(data) {
		t.Errorf("Unexpected message data")
	}
}
func TestBadExchange(t *testing.T) {
	publisher := TestPublisher{}
	handler := ingest.Handler(&publisher, &TestAccountService{false, new(big.Int)}, &TestAffiliateService{new(big.Int), nil})
	data, _ := hex.DecodeString("80fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000059938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	reader := TestReader{
		data[:377],
		nil,
	}
	request, _ := http.NewRequest("POST", "/", reader)
	request.Header["Content-Type"] = []string{"application/octet-stream"}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 400 {
		t.Errorf("Expected error code 202, got '%v'", recorder.Code)
	}
	body := recorder.Body.String()
	if body != "{\"code\":100,\"reason\":\"Validation Failed\",\"validationErrors\":[{\"field\":\"exchangeContractAddress\",\"code\":1002,\"reason\":\"Unknown exchangeContractAddress\"}]}" {
		t.Errorf("Got unexpected body: '%v' - %v", body, len(body))
	}
	if len(publisher.messages) != 0 {
		t.Errorf("Unexpected message count '%v'", len(publisher.messages))
		return
	}
}
