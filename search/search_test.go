package search_test

import (
	"crypto/rand"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/notegio/openrelay/blockhash"
	"github.com/notegio/openrelay/channels"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/search"
	"github.com/notegio/openrelay/types"
	"github.com/ethereum/go-ethereum/crypto"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"encoding/json"
	"os"
	// "reflect"
	"testing"
)

func sampleOrder(t *testing.T) *dbModule.Order {
	order := &types.Order{}
	if orderData, err := ioutil.ReadFile("../formatted_transaction.json"); err == nil {
		if err := json.Unmarshal(orderData, order); err != nil {
			t.Fatalf(err.Error())
		}
	}
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *order
	return dbOrder
}


func saltedSampleOrder(t *testing.T) *dbModule.Order {
	// TODO: Sign this order
	key, _ := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	order := sampleOrder(t)
	address := crypto.PubkeyToAddress(key.PublicKey)
	copy(order.Maker[:], address[:])
	rand.Read(order.Salt[:])

	hashedBytes := append([]byte("\x19Ethereum Signed Message:\n32"), order.Hash()...)
	signedBytes := crypto.Keccak256(hashedBytes)

	sig, _ := crypto.Sign(signedBytes, key)
	order.Signature[0] = sig[64] + 27
	copy(order.Signature[1:33], sig[0:32])
	copy(order.Signature[33:65], sig[32:64])
	order.Signature[65] = types.SigTypeEthSign
	return order
}

func getTestSearchHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	_, consumerChannel := channels.MockChannel()
	blockHash := blockhash.NewChanneledBlockHash(consumerChannel)
	return search.BlockHashDecorator(blockHash, search.SearchHandler(db))
}

func getTestOrderHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	_, consumerChannel := channels.MockChannel()
	blockHash := blockhash.NewChanneledBlockHash(consumerChannel)
	return search.BlockHashDecorator(blockHash, search.OrderHandler(db))
}

func getTestOrderBookHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	_, consumerChannel := channels.MockChannel()
	blockHash := blockhash.NewChanneledBlockHash(consumerChannel)
	return search.BlockHashDecorator(blockHash, search.OrderBookHandler(db))
}

func getDb() (*gorm.DB, error) {
	connectionString := fmt.Sprintf(
		"postgres://%v@%v",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_HOST"),
	)
	db, err := dbModule.GetDB(connectionString, os.Getenv("POSTGRES_PASSWORD"))
	// db.LogMode(true)
	return db, err
}

func TestFormatResponseJson(t *testing.T) {
	order := sampleOrder(t)
	orders := []dbModule.Order{*order, *order}
	response, contentType, err := search.FormatResponse(orders, "application/json")
	if err != nil {
		t.Errorf("Error getting formatted response: %v", err.Error())
	}
	if contentType != "application/json" {
		t.Errorf("Expected content type application/json, got '%v'", contentType)
	}
	if string(response) != "[{\"makerAddress\":\"0x627306090abab3a6e1400e9345bc60c78a8bef57\",\"takerAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetData\":\"0xf47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba\",\"takerAssetData\":\"0xf47261b000000000000000000000000005d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"feeRecipientAddress\":\"0x0000000000000000000000000000000000000000\",\"exchangeAddress\":\"0x90fe2af704b34e0224bf2299c838e04d4dcf1364\",\"senderAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetAmount\":\"50000000000000000000\",\"takerAssetAmount\":\"1000000000000000000\",\"makerFee\":\"0\",\"takerFee\":\"0\",\"expirationTimeSeconds\":\"5797808836\",\"salt\":\"11065671350908846865864045738088581419204014210814002044381812654087807531\",\"signature\":\"0x1c34f27415dc0177bc4016d48c3ec7eb19ee31124bcf4ca2eb3aba767c24e4712043bf8e49d1e28c6efa5a5e8a6824886700f356a403e0e66c75621e56b184b47b03\"},{\"makerAddress\":\"0x627306090abab3a6e1400e9345bc60c78a8bef57\",\"takerAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetData\":\"0xf47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba\",\"takerAssetData\":\"0xf47261b000000000000000000000000005d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"feeRecipientAddress\":\"0x0000000000000000000000000000000000000000\",\"exchangeAddress\":\"0x90fe2af704b34e0224bf2299c838e04d4dcf1364\",\"senderAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetAmount\":\"50000000000000000000\",\"takerAssetAmount\":\"1000000000000000000\",\"makerFee\":\"0\",\"takerFee\":\"0\",\"expirationTimeSeconds\":\"5797808836\",\"salt\":\"11065671350908846865864045738088581419204014210814002044381812654087807531\",\"signature\":\"0x1c34f27415dc0177bc4016d48c3ec7eb19ee31124bcf4ca2eb3aba767c24e4712043bf8e49d1e28c6efa5a5e8a6824886700f356a403e0e66c75621e56b184b47b03\"}]" {
		t.Errorf("Got '%v'", string(response))
	}
}

// func TestFormatResponseBin(t *testing.T) {
// 	order := sampleOrder(t)
// 	orders := []dbModule.Order{*order, *order}
// 	response, contentType, err := search.FormatResponse(orders, "application/octet-stream")
// 	if err != nil {
// 		t.Errorf("Error getting formatted response: %v", err.Error())
// 	}
// 	if contentType != "application/octet-stream" {
// 		t.Errorf("Expected content type application/octet-stream, got '%v'", contentType)
// 	}
// 	orderBytes := order.Bytes()
// 	orderValue := []byte{}
// 	orderValue = append(orderValue, orderBytes[:]...)
// 	orderValue = append(orderValue, orderBytes[:]...)
// 	if !reflect.DeepEqual(response, orderValue) {
// 		t.Errorf("Got '%#x'", response)
// 	}
// }
func TestFormatSingleResponseJson(t *testing.T) {
	order := sampleOrder(t)
	response, contentType, err := search.FormatSingleResponse(order, "application/json")
	if err != nil {
		t.Errorf("Error getting formatted response: %v", err.Error())
	}
	if contentType != "application/json" {
		t.Errorf("Expected content type application/json, got '%v'", contentType)
	}
	if string(response) != "{\"makerAddress\":\"0x627306090abab3a6e1400e9345bc60c78a8bef57\",\"takerAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetData\":\"0xf47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba\",\"takerAssetData\":\"0xf47261b000000000000000000000000005d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"feeRecipientAddress\":\"0x0000000000000000000000000000000000000000\",\"exchangeAddress\":\"0x90fe2af704b34e0224bf2299c838e04d4dcf1364\",\"senderAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetAmount\":\"50000000000000000000\",\"takerAssetAmount\":\"1000000000000000000\",\"makerFee\":\"0\",\"takerFee\":\"0\",\"expirationTimeSeconds\":\"5797808836\",\"salt\":\"11065671350908846865864045738088581419204014210814002044381812654087807531\",\"signature\":\"0x1c34f27415dc0177bc4016d48c3ec7eb19ee31124bcf4ca2eb3aba767c24e4712043bf8e49d1e28c6efa5a5e8a6824886700f356a403e0e66c75621e56b184b47b03\"}" {
		t.Errorf("Got '%v'", string(response))
	}
}

// func TestFormatSingleResponseBin(t *testing.T) {
// 	order := sampleOrder(t)
// 	response, contentType, err := search.FormatSingleResponse(order, "application/octet-stream")
// 	if err != nil {
// 		t.Errorf("Error getting formatted response: %v", err.Error())
// 	}
// 	if contentType != "application/octet-stream" {
// 		t.Errorf("Expected content type application/json, got '%v'", contentType)
// 	}
// 	orderBytes := order.Bytes()
// 	orderValue := []byte{}
// 	orderValue = append(orderValue, orderBytes[:]...)
// 	if !reflect.DeepEqual(response, orderValue) {
// 		t.Errorf("Got '%#x'", response)
// 	}
// }

func TestBlockhashRedirect(t *testing.T) {
	_, consumerChannel := channels.MockChannel()
	blockHash := blockhash.NewChanneledBlockHash(consumerChannel)
	handler := search.BlockHashDecorator(blockHash, search.SearchHandler(nil))
	request, _ := http.NewRequest("GET", "/v0/orders?makertoken=0x627306090abab3a6e1400e9345bc60c78a8bef57", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 307 {
		t.Errorf("Did not redirect")
	}
	if location := recorder.Header().Get("Location"); location != "/v0/orders?blockhash=initializing&makertoken=0x627306090abab3a6e1400e9345bc60c78a8bef57" {
		t.Errorf("Expected orderHash to be added, got '%v'", location)
	}
	if cacheControl := recorder.Header().Get("Cache-Control"); cacheControl != "max-age=5, public" {
		t.Errorf("Cache-Control header not as expect. '%v'", cacheControl)
	}
}

func filterContractRequest(queryString, emptyQueryString string, t *testing.T) {
	db, err := getDb()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	if err := sampleOrder(t).Save(tx, 0).Error; err != nil {
		t.Fatalf(err.Error())
	}
	handler := getTestSearchHandler(tx)
	request, _ := http.NewRequest("GET", "/v0/orders?"+queryString+"&blockhash=x&_expTime=0", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Unexpected response code '%v'", recorder.Code)
	}
	response := recorder.Body.String()
	if string(response) != "[{\"makerAddress\":\"0x627306090abab3a6e1400e9345bc60c78a8bef57\",\"takerAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetData\":\"0xf47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba\",\"takerAssetData\":\"0xf47261b000000000000000000000000005d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"feeRecipientAddress\":\"0x0000000000000000000000000000000000000000\",\"exchangeAddress\":\"0x90fe2af704b34e0224bf2299c838e04d4dcf1364\",\"senderAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetAmount\":\"50000000000000000000\",\"takerAssetAmount\":\"1000000000000000000\",\"makerFee\":\"0\",\"takerFee\":\"0\",\"expirationTimeSeconds\":\"5797808836\",\"salt\":\"11065671350908846865864045738088581419204014210814002044381812654087807531\",\"signature\":\"0x1c34f27415dc0177bc4016d48c3ec7eb19ee31124bcf4ca2eb3aba767c24e4712043bf8e49d1e28c6efa5a5e8a6824886700f356a403e0e66c75621e56b184b47b03\"}]" {
		t.Errorf("Got '%v'", string(response))
	}
	request, _ = http.NewRequest("GET", "/v0/orders?"+emptyQueryString+"&blockhash=x", nil)
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
	filterContractRequest("takerAssetAddress=0x05d090b51c40b020eab3bfcb6a2dff130df22e9c", "takerAssetAddress=0x90fe2af704b34e0224bf2299c838e04d4dcf1300", t)
}

func TestFilterMakerTokenContract(t *testing.T) {
	filterContractRequest("makerAssetAddress=0x1dad4783cf3fe3085c1426157ab175a6119a04ba", "makerAssetAddress=0x90fe2af704b34e0224bf2299c838e04d4dcf1300", t)
}

func TestFilterTokenContract(t *testing.T) {
	filterContractRequest("assetAddress=0x1dad4783cf3fe3085c1426157ab175a6119a04ba", "assetAddress=0x90fe2af704b34e0224bf2299c838e04d4dcf1300", t)
}

func TestFilterMakerTokenAndTakerTokenContract(t *testing.T) {
	filterContractRequest("makerAssetAddress=0x1dad4783cf3fe3085c1426157ab175a6119a04ba&takerAssetAddress=0x05d090b51c40b020eab3bfcb6a2dff130df22e9c", "makerAssetAddress=0x90fe2af704b34e0224bf2299c838e04d4dcf1300&takerAssetAddress=0x05d090b51c40b020eab3bfcb6a2dff130df22e9c", t)
}

func TestFilterMaker(t *testing.T) {
	filterContractRequest("maker=0x627306090abab3a6e1400e9345bc60c78a8bef57", "maker=0x90fe2af704b34e0224bf2299c838e04d4dcf1300", t)
}

func TestFilterTrader(t *testing.T) {
	filterContractRequest("trader=0x627306090abab3a6e1400e9345bc60c78a8bef57", "trader=0x90fe2af704b34e0224bf2299c838e04d4dcf1300", t)
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
	defer func() {
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	for i := 0; i < 21; i++ {
		saltedSampleOrder(t).Save(tx, 0)
	}
	handler := getTestSearchHandler(tx)
	request, _ := http.NewRequest("GET", "/v0/orders?&blockhash=x&_expTime=0", nil)
	request.Header.Set("Accept", "application/json")
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Unexpected response code '%v'", recorder.Code)
	}
	orders := []dbModule.Order{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &orders); err != nil {
		t.Fatalf(err.Error())
	}
	if length := len(orders); length != 20 {
			t.Errorf("Expected 20 items, got '%v'", (length))
	}
	// if length := recorder.Body.Len(); length != (20 * 441) {
	// 	t.Errorf("Expected 20 items, got '%v'", (length / 441))
	// }
	request, _ = http.NewRequest("GET", "/v0/orders?page=2&blockhash=x&_expTime=0", nil)
	request.Header.Set("Accept", "application/json")
	recorder = httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Unexpected response code '%v'", recorder.Code)
	}
	orders = []dbModule.Order{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &orders); err != nil {
		t.Fatalf(err.Error())
	}
	if length := len(orders); length != 1 {
			t.Errorf("Expected 1 items, got '%v'", (length))
	}
}

func TestOrderLookup(t *testing.T) {
	db, err := getDb()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	order := sampleOrder(t)
	order.Save(tx, 0)
	orderHash := order.Hash()
	orderHashHex := hex.EncodeToString(orderHash)
	handler := getTestOrderHandler(tx)
	request, _ := http.NewRequest("GET", "/v0/order/0x"+orderHashHex+"?blockhash=x", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Unexpected response code '%v'", recorder.Code)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Unexpected Content-Type: '%v'", contentType)
	}
	response := recorder.Body.String()
	if string(response) != "{\"makerAddress\":\"0x627306090abab3a6e1400e9345bc60c78a8bef57\",\"takerAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetData\":\"0xf47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba\",\"takerAssetData\":\"0xf47261b000000000000000000000000005d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"feeRecipientAddress\":\"0x0000000000000000000000000000000000000000\",\"exchangeAddress\":\"0x90fe2af704b34e0224bf2299c838e04d4dcf1364\",\"senderAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetAmount\":\"50000000000000000000\",\"takerAssetAmount\":\"1000000000000000000\",\"makerFee\":\"0\",\"takerFee\":\"0\",\"expirationTimeSeconds\":\"5797808836\",\"salt\":\"11065671350908846865864045738088581419204014210814002044381812654087807531\",\"signature\":\"0x1c34f27415dc0177bc4016d48c3ec7eb19ee31124bcf4ca2eb3aba767c24e4712043bf8e49d1e28c6efa5a5e8a6824886700f356a403e0e66c75621e56b184b47b03\"}" {
		t.Errorf("Got '%v'", string(response))
	}
}

func TestPairLookup(t *testing.T) {
	db, err := getDb()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	order := sampleOrder(t)
	order.Save(tx, 0)
	handler := search.PairHandler(tx)
	request, _ := http.NewRequest("GET", "/v0/token_pairs", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Unexpected response code '%v'", recorder.Code)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Unexpected Content-Type: '%v'", contentType)
	}
	response := recorder.Body.String()
	if string(response) != "[{\"tokenA\":{\"address\":\"0x05d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"minAmount\":\"1\",\"maxAmount\":\"115792089237316195423570985008687907853269984665640564039457584007913129639935\",\"precision\":5},\"tokenB\":{\"address\":\"0x1dad4783cf3fe3085c1426157ab175a6119a04ba\",\"minAmount\":\"1\",\"maxAmount\":\"115792089237316195423570985008687907853269984665640564039457584007913129639935\",\"precision\":5}}]" {
		t.Errorf("Got unexpected JSON response '%v'", string(response))
	}
}

func TestOrderBookLookup(t *testing.T) {
	db, err := getDb()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	order := sampleOrder(t)
	order.Save(tx, 0)
	handler := getTestOrderBookHandler(tx)
	request, _ := http.NewRequest("GET", "/v0/orderbook?blockhash=x&quoteTokenAddress=0x1dad4783cf3fe3085c1426157ab175a6119a04ba&baseTokenAddress=0x05d090b51c40b020eab3bfcb6a2dff130df22e9c&_expTime=0", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Unexpected response code '%v'", recorder.Code)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Unexpected Content-Type: '%v'", contentType)
	}
	response := recorder.Body.String()
	if string(response) != "{\"asks\":[],\"bids\":[{\"makerAddress\":\"0x627306090abab3a6e1400e9345bc60c78a8bef57\",\"takerAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetData\":\"0xf47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba\",\"takerAssetData\":\"0xf47261b000000000000000000000000005d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"feeRecipientAddress\":\"0x0000000000000000000000000000000000000000\",\"exchangeAddress\":\"0x90fe2af704b34e0224bf2299c838e04d4dcf1364\",\"senderAddress\":\"0x0000000000000000000000000000000000000000\",\"makerAssetAmount\":\"50000000000000000000\",\"takerAssetAmount\":\"1000000000000000000\",\"makerFee\":\"0\",\"takerFee\":\"0\",\"expirationTimeSeconds\":\"5797808836\",\"salt\":\"11065671350908846865864045738088581419204014210814002044381812654087807531\",\"signature\":\"0x1c34f27415dc0177bc4016d48c3ec7eb19ee31124bcf4ca2eb3aba767c24e4712043bf8e49d1e28c6efa5a5e8a6824886700f356a403e0e66c75621e56b184b47b03\"}]}" {
		t.Errorf("Got '%v'", string(response))
	}
}
