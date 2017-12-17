package db_test

import (
	"encoding/hex"
	"github.com/jinzhu/gorm"
	dbModule "github.com/notegio/openrelay/db"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/notegio/openrelay/types"
	"testing"
	"fmt"
	"os"
	"reflect"
)

func getTestOrderBytes() [441]byte {
	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000059938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
	var testOrderByteArray [441]byte
	copy(testOrderByteArray[:], testOrderBytes[:])
	return testOrderByteArray
}

func sampleOrder() *types.Order {
	order := &types.Order{}
	order.FromBytes(getTestOrderBytes())
	return order
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

func TestSaveOrder(t *testing.T) {
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
	order := sampleOrder()
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *order
	if err := dbOrder.Save(tx, dbModule.StatusOpen).Error; err != nil {
		t.Errorf(err.Error())
	}
}

func TestQueryOrder(t *testing.T) {
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
	order := sampleOrder()
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
}

func checkPairs(t *testing.T, tokenPairs []dbModule.Pairs, sOrder *types.Order) {
	if len(tokenPairs) != 1 {
		t.Errorf("Expected 1 value, got %v", len(tokenPairs))
		return
	}
	if ! reflect.DeepEqual(tokenPairs[0].TokenA, sOrder.TakerToken) {
		t.Errorf("Expected %#x, got %#x", sOrder.TakerToken[:], tokenPairs[0].TokenA[:])
	}
	if ! reflect.DeepEqual(tokenPairs[0].TokenB, sOrder.MakerToken) {
		t.Errorf("Expected %#x, got %#x", sOrder.MakerToken[:], tokenPairs[0].TokenB[:])
	}
}

func TestQueryPairs(t *testing.T) {
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
	sOrder := sampleOrder()
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
	defer func(){
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	sOrder := sampleOrder()
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *sOrder
	if err := dbOrder.Save(tx, dbModule.StatusOpen).Error; err != nil {
		t.Errorf(err.Error())
	}
	tokenPairs, err := dbModule.GetTokenAPairs(tx, sOrder.TakerToken, 0, 10)
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
	defer func(){
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	sOrder := sampleOrder()
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *sOrder
	if err := dbOrder.Save(tx, dbModule.StatusOpen).Error; err != nil {
		t.Errorf(err.Error())
	}
	tokenPairs, err := dbModule.GetTokenABPairs(tx, sOrder.TakerToken, sOrder.MakerToken)
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
	defer func(){
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	sOrder := sampleOrder()
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
		t.Errorf("Expected 0 values, got %v", len(tokenPairs))
		return
	}
}
