package pool_test

import (
	"net/http"
	"net/http/httptest"
	dbModule "github.com/notegio/openrelay/db"
	poolModule "github.com/notegio/openrelay/pool"
	"github.com/notegio/openrelay/types"
	"golang.org/x/crypto/sha3"
	"crypto/hmac"
	"crypto/sha256"
	"testing"
	"bytes"
	"fmt"
	"time"
)

type TestPublisher struct {
	channel  string
	messages []string
}

func (pub *TestPublisher) Publish(message string) bool {
	pub.messages = append(pub.messages, message)
	return true
}


func TestCancellation(t *testing.T) {
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
	if err := tx.AutoMigrate(&dbModule.Exchange{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	if err := tx.AutoMigrate(&poolModule.Pool{}).Error; err != nil {
		t.Errorf(err.Error())
	}

	poolHash := sha3.NewLegacyKeccak256()
	poolHash.Write([]byte("testPool"))
	poolID := poolHash.Sum(nil)
	poolHash2 := sha3.NewLegacyKeccak256()
	poolHash2.Write([]byte("testPool2"))
	poolID2 := poolHash2.Sum(nil)
	tx.Create(&poolModule.Pool{ID: poolID, SenderAddresses: types.NetworkAddressMap{}, FilterAddresses: types.NetworkAddressMap{}, SearchTerms: "makerAssetData=0xf47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba", SecretKey: []byte("secret")})
	tx.Create(&poolModule.Pool{ID: poolID2, SenderAddresses: types.NetworkAddressMap{}, FilterAddresses: types.NetworkAddressMap{}, SearchTerms: "takerAssetData=0xf47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba", SecretKey: []byte("secret2")})

	order := sampleOrder(t)
	order.PoolID = poolID
	order.Save(tx, dbModule.StatusOpen, nil)

	var poolCount int
	sourcePublisher := &TestPublisher{"cancellations", []string{}}
	tx.Model(&poolModule.Pool{}).Count(&poolCount)
	handler := poolModule.PoolDecoratorBaseFee(tx, nil, poolModule.PoolAdminHandler(tx, sourcePublisher))
	message := []byte(fmt.Sprintf(`{"id":1,"method":"cancellation","params":["%#x"],"expiration":%v}`, order.Hash(), time.Now().Add(10 * time.Second).Unix()))
	h := hmac.New(sha256.New, []byte("secret"))
	checksum := h.Sum(message)
	request, _ := http.NewRequest("POST", "/testPool/v3/admin", bytes.NewReader(message))
	request.Header["Authorization"] = []string{fmt.Sprintf("%#x", checksum)}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf(recorder.Body.String())
	}
	if len(sourcePublisher.messages) != 1 {
		t.Errorf("Expected cancellation message")
	}
}

func TestCancellationWrongPool(t *testing.T) {
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
	if err := tx.AutoMigrate(&dbModule.Exchange{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	if err := tx.AutoMigrate(&poolModule.Pool{}).Error; err != nil {
		t.Errorf(err.Error())
	}

	poolHash := sha3.NewLegacyKeccak256()
	poolHash.Write([]byte("testPool"))
	poolID := poolHash.Sum(nil)
	poolHash2 := sha3.NewLegacyKeccak256()
	poolHash2.Write([]byte("testPool2"))
	poolID2 := poolHash2.Sum(nil)
	tx.Create(&poolModule.Pool{ID: poolID, SenderAddresses: types.NetworkAddressMap{}, FilterAddresses: types.NetworkAddressMap{}, SearchTerms: "makerAssetData=0xf47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba", SecretKey: []byte("secret")})
	tx.Create(&poolModule.Pool{ID: poolID2, SenderAddresses: types.NetworkAddressMap{}, FilterAddresses: types.NetworkAddressMap{}, SearchTerms: "takerAssetData=0xf47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba", SecretKey: []byte("secret2")})

	order := sampleOrder(t)
	order.PoolID = poolID2
	order.Save(tx, dbModule.StatusOpen, nil)

	var poolCount int
	sourcePublisher := &TestPublisher{"cancellations", []string{}}
	tx.Model(&poolModule.Pool{}).Count(&poolCount)
	handler := poolModule.PoolDecoratorBaseFee(tx, nil, poolModule.PoolAdminHandler(tx, sourcePublisher))
	message := []byte(fmt.Sprintf(`{"id":1,"method":"cancellation","params":["%#x"],"expiration":%v}`, order.Hash(), time.Now().Add(10 * time.Second).Unix()))
	h := hmac.New(sha256.New, []byte("secret"))
	checksum := h.Sum(message)
	request, _ := http.NewRequest("POST", "/testPool/v3/admin", bytes.NewReader(message))
	request.Header["Authorization"] = []string{fmt.Sprintf("%#x", checksum)}
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf(recorder.Body.String())
	}
	if len(sourcePublisher.messages) != 0 {
		t.Errorf("Expected cancellation message")
	}
}
