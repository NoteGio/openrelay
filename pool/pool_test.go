package pool_test

import (
	dbModule "github.com/notegio/openrelay/db"
	poolModule "github.com/notegio/openrelay/pool"
	"github.com/notegio/openrelay/types"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/jinzhu/gorm"
	"os"
	"fmt"
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
	"io/ioutil"
	"encoding/json"
)



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

func sampleOrder(t *testing.T) *dbModule.Order {
	order := &types.Order{}
	if orderData, err := ioutil.ReadFile("../formatted_transaction.json"); err == nil {
		if err := json.Unmarshal(orderData, order); err != nil {
			t.Fatalf(err.Error())
		}
	}
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *order
	dbOrder.Populate()
	return dbOrder
}

func TestDecorator(t *testing.T) {
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
	if err := tx.AutoMigrate(&poolModule.Pool{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	poolID := sha3.NewKeccak256().Sum([]byte("testPool"))
	tx.Create(&poolModule.Pool{ID: poolID, SenderAddress: &types.Address{}, FilterAddress: &types.Address{}})

	handler := poolModule.PoolDecorator(tx, func(w http.ResponseWriter, r *http.Request, pool *poolModule.Pool) {
		if !bytes.Equal(pool.ID, []byte("")) {
			t.Errorf("Unexpected pool id: '%#x'", pool.ID)
		}
	})
	handler2 := poolModule.PoolDecorator(tx, func(w http.ResponseWriter, r *http.Request, pool *poolModule.Pool) {
		if !bytes.Equal(pool.ID, poolID) {
			t.Errorf("Unexpected pool id '%#x' != %#x", pool.ID, poolID)
		}
	})
	handler3 := poolModule.PoolDecorator(tx, func(w http.ResponseWriter, r *http.Request, pool *poolModule.Pool) {
		t.Errorf("This should not be reached")
	})

	request, _ := http.NewRequest("GET", "/v2/content", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	request, _ = http.NewRequest("GET", "/testPool/v2/content", nil)
	recorder = httptest.NewRecorder()
	handler2(recorder, request)
	request, _ = http.NewRequest("GET", "/unknownPool/v2/content", nil)
	recorder = httptest.NewRecorder()
	handler3(recorder, request)
	if recorder.Code != 404 {
		t.Errorf("Unexpected status code")
	}
}

func TestPoolFilter(t *testing.T) {
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
	order := sampleOrder(t)
	order.Save(tx, dbModule.StatusOpen)

	poolID := sha3.NewKeccak256().Sum([]byte("testPool"))
	poolID2 := sha3.NewKeccak256().Sum([]byte("testPool2"))
	tx.Create(&poolModule.Pool{ID: poolID, SenderAddress: &types.Address{}, FilterAddress: &types.Address{}, SearchTerms: "makerAssetData=0xf47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba"})
	tx.Create(&poolModule.Pool{ID: poolID2, SenderAddress: &types.Address{}, FilterAddress: &types.Address{}, SearchTerms: "takerAssetData=0xf47261b00000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba"})

	var poolCount int
	tx.Model(&poolModule.Pool{}).Count(&poolCount)

	handler := poolModule.PoolDecorator(tx, func(w http.ResponseWriter, r *http.Request, pool *poolModule.Pool) {
		query := tx.Model(&dbModule.Order{})
		query, err := pool.Filter(query)
		if err != nil {
			t.Errorf(err.Error())
		}
		var count int
		query.Count(&count)
		if count != 1 {
			t.Errorf("Expected 1 item in handler query, got %v", count)
		}
		w.WriteHeader(200)
	})
	emptyHandler := poolModule.PoolDecorator(tx, func(w http.ResponseWriter, r *http.Request, pool *poolModule.Pool) {
		query := tx.Model(&dbModule.Order{})
		query, err := pool.Filter(query)
		if err != nil {
			t.Errorf(err.Error())
		}
		var count int
		query.Count(&count)
		if count != 0 {
			t.Errorf("Expected 0 item in handler query, got %v", count)
		}
		w.WriteHeader(200)
	})

	request, _ := http.NewRequest("GET", "/v2/content", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Error on handler")
	}
	request, _ = http.NewRequest("GET", "/testPool/v2/content", nil)
	recorder = httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Error on handler2")
	}
	request, _ = http.NewRequest("GET", "/testPool2/v2/content", nil)
	recorder = httptest.NewRecorder()
	emptyHandler(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Error on emptyHandler: %v - %v", recorder.Code, recorder.Body.String())
	}

}
