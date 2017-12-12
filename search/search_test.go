package search_test

import (
	"testing"
	"encoding/hex"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/notegio/openrelay/search"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/blockhash"
	"github.com/notegio/openrelay/types"
	dbModule "github.com/notegio/openrelay/db"
	"net/http"
	"net/http/httptest"
	"reflect"
	"os"
	"fmt"
	"crypto/rand"
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

func saltedSampleOrder() *dbModule.Order {
	order := &types.Order{}
	order.FromBytes(getTestOrderBytes())
	rand.Read(order.Salt[:])
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *order
	return dbOrder
}

func getTestHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	_, consumerChannel := channels.MockChannel()
	blockHash := blockhash.NewChanneledBlockHash(consumerChannel)
	return search.Handler(db, blockHash)
}

func getDb() (*gorm.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%v sslmode=disable user=%v password=%v",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
	)
	db, err := gorm.Open("postgres", connectionString)
	// db.LogMode(true)
	return db, err
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
	_, consumerChannel := channels.MockChannel()
	blockHash := blockhash.NewChanneledBlockHash(consumerChannel)
	handler := search.Handler(nil, blockHash)
	request, _ := http.NewRequest("GET", "/v0/orders?makertoken=0x324454186bb728a3ea55750e0618ff1b18ce6cf8", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 307 {
		t.Errorf("Did not redirect")
	}
	if location := recorder.Header().Get("Location"); location != "/v0/orders?blockhash=initializing&makertoken=0x324454186bb728a3ea55750e0618ff1b18ce6cf8" {
		t.Errorf("Expected orderHash to be added, got '%v'", location)
	}
}

func filterContractRequest(queryString, emptyQueryString string, t *testing.T) {
	db, err := getDb()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	tx := db.Begin()
	defer func(){
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	sampleOrder().Save(tx, 0)
	handler := getTestHandler(tx)
	request, _ := http.NewRequest("GET", "/v0/orders?" + queryString + "&blockhash=x", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Unexpected response code '%v'", recorder.Code)
	}
	response := recorder.Body.String()
	if string(response) != "[{\"maker\":\"0x324454186bb728a3ea55750e0618ff1b18ce6cf8\",\"taker\":\"0x0000000000000000000000000000000000000000\",\"makerTokenAddress\":\"0x1dad4783cf3fe3085c1426157ab175a6119a04ba\",\"takerTokenAddress\":\"0x05d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"feeRecipient\":\"0x0000000000000000000000000000000000000000\",\"exchangeContractAddress\":\"0x90fe2af704b34e0224bf2299c838e04d4dcf1364\",\"makerTokenAmount\":\"50000000000000000000\",\"takerTokenAmount\":\"1000000000000000000\",\"makerFee\":\"0\",\"takerFee\":\"0\",\"expirationUnixTimestampSec\":\"1502841540\",\"salt\":\"11065671350908846865864045738088581419204014210814002044381812654087807531\",\"ecSignature\":{\"v\":\"27\",\"r\":\"0x021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e\",\"s\":\"0x12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1\"},\"takerTokenAmountFilled\":\"0\",\"takerTokenAmountCancelled\":\"0\"}]" {
		t.Errorf("Got '%v'", string(response))
	}
	request, _ = http.NewRequest("GET", "/v0/orders?" + emptyQueryString + "&blockhash=x", nil)
	recorder = httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Unexpected response code '%v'", recorder.Code)
	}
	response = recorder.Body.String()
	if string(response) != "[]" {
		t.Errorf("Got '%v'", string(response))
	}
}

func TestFilterExchangeContract(t *testing.T) {
	filterContractRequest("exchangeContractAddress=0x90fe2af704b34e0224bf2299c838e04d4dcf1364", "exchangeContractAddress=0x90fe2af704b34e0224bf2299c838e04d4dcf1300", t)
}

func TestFilterTakerTokenContract(t *testing.T) {
	filterContractRequest("takerTokenAddress=0x05d090b51c40b020eab3bfcb6a2dff130df22e9c", "takerTokenAddress=0x90fe2af704b34e0224bf2299c838e04d4dcf1300", t)
}

func TestFilterMakerTokenContract(t *testing.T) {
	filterContractRequest("makerTokenAddress=0x1dad4783cf3fe3085c1426157ab175a6119a04ba", "makerTokenAddress=0x90fe2af704b34e0224bf2299c838e04d4dcf1300", t)
}

func TestFilterTokenContract(t *testing.T) {
	filterContractRequest("tokenAddress=0x1dad4783cf3fe3085c1426157ab175a6119a04ba", "tokenAddress=0x90fe2af704b34e0224bf2299c838e04d4dcf1300", t)
}

func TestFilterMakerTokenAndTakerTokenContract(t *testing.T) {
	filterContractRequest("makerTokenAddress=0x1dad4783cf3fe3085c1426157ab175a6119a04ba&takerTokenAddress=0x05d090b51c40b020eab3bfcb6a2dff130df22e9c", "makerTokenAddress=0x90fe2af704b34e0224bf2299c838e04d4dcf1300&takerTokenAddress=0x05d090b51c40b020eab3bfcb6a2dff130df22e9c", t)
}

func TestFilterMaker(t *testing.T) {
	filterContractRequest("maker=0x324454186bb728a3ea55750e0618ff1b18ce6cf8", "maker=0x90fe2af704b34e0224bf2299c838e04d4dcf1300", t)
}

func TestFilterTrader(t *testing.T) {
	filterContractRequest("trader=0x324454186bb728a3ea55750e0618ff1b18ce6cf8", "trader=0x90fe2af704b34e0224bf2299c838e04d4dcf1300", t)
}

func TestFilterTaker(t *testing.T) {
	filterContractRequest("taker=0x0000000000000000000000000000000000000000", "taker=0x90fe2af704b34e0224bf2299c838e04d4dcf1300", t)
}

func TestFilterFeeRecipientTaker(t *testing.T) {
	filterContractRequest("feeRecipient=0x0000000000000000000000000000000000000000", "feeRecipient=0x90fe2af704b34e0224bf2299c838e04d4dcf1300", t)
}

func TestPagination(t *testing.T) {
	db, err := getDb()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	tx := db.Begin()
	defer func(){
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	for i:= 0; i < 21; i++ {
		saltedSampleOrder().Save(tx, 0)
	}
	handler := getTestHandler(tx)
	request, _ := http.NewRequest("GET", "/v0/orders?&blockhash=x", nil)
	request.Header.Set("Accept", "application/octet-stream")
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Unexpected response code '%v'", recorder.Code)
	}
	if length := recorder.Body.Len(); length != (20 * 441) {
		t.Errorf("Expected 20 items, got '%v'", (length / 441))
	}
	request, _ = http.NewRequest("GET", "/v0/orders?page=2&blockhash=x", nil)
	request.Header.Set("Accept", "application/octet-stream")
	recorder = httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Unexpected response code '%v'", recorder.Code)
	}
	if length := recorder.Body.Len(); length != (1 * 441) {
		t.Errorf("Expected 1 items, got '%v'", (length / 441))
	}
}
