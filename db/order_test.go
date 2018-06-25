package db_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/types"
	"os"
	"reflect"
	"testing"
	"bytes"
	"io/ioutil"
)

func sampleOrder(t *testing.T) *types.Order {
	order := &types.Order{}
	if orderData, err := ioutil.ReadFile("../formatted_transaction.json"); err == nil {
		if err := json.Unmarshal(orderData, order); err != nil {
			t.Fatalf(err.Error())
		}
	}
	return order
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

func TestSaveOrder(t *testing.T) {
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
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *order
	if err := dbOrder.Save(tx, dbModule.StatusOpen).Error; err != nil {
		t.Errorf(err.Error())
	}
}

func TestFailOnEmptyOrder(t *testing.T) {
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
	dbOrder := &dbModule.Order{}
	dbOrder.Initialize()
	if err := dbOrder.Save(tx, dbModule.StatusOpen).Error; err == nil {
		t.Errorf("Expected error saving empty order")
	}
}

func TestQueryOrder(t *testing.T) {
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
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *order
	if err := dbOrder.Save(tx, dbModule.StatusOpen).Error; err != nil {
		t.Errorf(err.Error())
	}
	dbOrders := []dbModule.Order{}
	tx.Model(&dbModule.Order{}).Where("order_hash = ?", order.Hash()).Find(&dbOrders)
	dbOrder = &dbOrders[0]
	if !reflect.DeepEqual(dbOrder.Bytes(), order.Bytes()) {
		dbBytes := dbOrder.Bytes()
		orderBytes := order.Bytes()
		t.Errorf(
			"Queried order not equal to saved order; '%v' != '%v'",
			hex.EncodeToString(dbBytes[:]),
			hex.EncodeToString(orderBytes[:]),
		)
	}
	if dbOrder.Status != dbModule.StatusOpen {
		t.Errorf("Order unexpectedly not open - %v", dbOrder.Status)
	}
	if dbOrder.Price != 0.02 {
		t.Errorf("Expected price '0.02' got '%v'", dbOrder.Price)
	}
	if dbOrder.FeeRate != 0 {
		t.Errorf("Expected FeeRate '0' got '%v'", dbOrder.FeeRate)
	}
	fmt.Printf("Filled: %#x", dbOrder.TakerAssetAmountFilled)
	if !bytes.Equal(dbOrder.MakerAssetRemaining[:], dbOrder.MakerAssetAmount[:]) {
		t.Errorf("Unexpected MakerAssetRemaining, expected %#x got : %#x", dbOrder.MakerAssetAmount, dbOrder.MakerAssetRemaining)
	}
	if !bytes.Equal(dbOrder.MakerFeeRemaining[:], dbOrder.MakerFee[:]) {
		t.Errorf("Unexpected MakerAssetRemaining, expected %#x got : %#x", dbOrder.MakerFee, dbOrder.MakerAssetRemaining)
	}
}

func checkPairs(t *testing.T, tokenPairs []dbModule.Pair, sOrder *types.Order) {
	if len(tokenPairs) != 1 {
		t.Errorf("Expected 1 value, got %v", len(tokenPairs))
		return
	}
	if !reflect.DeepEqual(tokenPairs[0].TokenA, sOrder.TakerAssetData.Address()) {
		t.Errorf("Expected %#x, got %#x", sOrder.TakerAssetData.Address(), tokenPairs[0].TokenA[:])
	}
	if !reflect.DeepEqual(tokenPairs[0].TokenB, sOrder.MakerAssetData.Address()) {
		t.Errorf("Expected %#x, got %#x", sOrder.MakerAssetData.Address(), tokenPairs[0].TokenB[:])
	}
}

func TestQueryPairs(t *testing.T) {
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
	sOrder := sampleOrder(t)
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *sOrder
	if err := dbOrder.Save(tx, dbModule.StatusOpen).Error; err != nil {
		t.Errorf(err.Error())
	}
	tokenPairs, err := dbModule.GetAllTokenPairs(tx, 0, 10)
	if err != nil {
		t.Errorf(err.Error())
	}
	checkPairs(t, tokenPairs, sOrder)
}

func TestQueryPairsTokenAFilter(t *testing.T) {
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
	sOrder := sampleOrder(t)
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *sOrder
	if err := dbOrder.Save(tx, dbModule.StatusOpen).Error; err != nil {
		t.Errorf(err.Error())
	}
	tokenPairs, err := dbModule.GetTokenAPairs(tx, sOrder.TakerAssetData.Address(), 0, 10)
	if err != nil {
		t.Errorf(err.Error())
	}
	checkPairs(t, tokenPairs, sOrder)
}

func TestQueryPairsTokenABFilter(t *testing.T) {
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
	sOrder := sampleOrder(t)
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *sOrder
	if err := dbOrder.Save(tx, dbModule.StatusOpen).Error; err != nil {
		t.Errorf(err.Error())
	}
	tokenPairs, err := dbModule.GetTokenABPairs(tx, sOrder.TakerAssetData.Address(), sOrder.MakerAssetData.Address())
	if err != nil {
		t.Errorf(err.Error())
	}
	checkPairs(t, tokenPairs, sOrder)
}

func TestQueryPairsTokenEmptyFilter(t *testing.T) {
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
	sOrder := sampleOrder(t)
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *sOrder
	if err := dbOrder.Save(tx, dbModule.StatusOpen).Error; err != nil {
		t.Errorf(err.Error())
	}
	tokenPairs, err := dbModule.GetTokenAPairs(tx, sOrder.Taker, 0, 10)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(tokenPairs) != 0 {
		t.Errorf("Expected 0 values, got %v", tokenPairs)
		return
	}
}

func TestMarshalPairs(t *testing.T) {
	sOrder := sampleOrder(t)
	pair := &dbModule.Pair{sOrder.MakerAssetData.Address(), sOrder.TakerAssetData.Address()}
	pairJSON, _ := json.Marshal(pair)
	if string(pairJSON) != "{\"tokenA\":{\"address\":\"0x1dad4783cf3fe3085c1426157ab175a6119a04ba\",\"minAmount\":\"1\",\"maxAmount\":\"115792089237316195423570985008687907853269984665640564039457584007913129639935\",\"precision\":5},\"tokenB\":{\"address\":\"0x05d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"minAmount\":\"1\",\"maxAmount\":\"115792089237316195423570985008687907853269984665640564039457584007913129639935\",\"precision\":5}}" {
		t.Errorf("Unexpected response, got '%v'", string(pairJSON))
	}
}
