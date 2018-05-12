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
)

func getTestOrderBytes() [441]byte {
	testOrderBytes, _ := hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364324454186bb728a3ea55750e0618ff1b18ce6cf800000000000000000000000000000000000000001dad4783cf3fe3085c1426157ab175a6119a04ba05d090b51c40b020eab3bfcb6a2dff130df22e9c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000159938ac4000643508ff7019bfb134363a86e98746f6c33262e68daf992b8df064217222b1b37adbc51c87a2f4c8c40c25fab5a73c65d078322f1db5739ee6fd49f18ce44637382de9b4cf7ceaf602f221132c9ddf41b83fb9666839022703da852d4ed88af")
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
	order := sampleOrder()
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
	fmt.Printf("Filled: %#x", dbOrder.TakerTokenAmountFilled)
	fmt.Printf("Cancelled: %#x", dbOrder.TakerTokenAmountCancelled)
	if !bytes.Equal(dbOrder.MakerTokenRemaining[:], dbOrder.MakerTokenAmount[:]) {
		t.Errorf("Unexpected MakerTokenRemaining, expected %#x got : %#x", dbOrder.MakerTokenAmount, dbOrder.MakerTokenRemaining)
	}
	if !bytes.Equal(dbOrder.MakerFeeRemaining[:], dbOrder.MakerFee[:]) {
		t.Errorf("Unexpected MakerTokenRemaining, expected %#x got : %#x", dbOrder.MakerFee, dbOrder.MakerTokenRemaining)
	}
}

func checkPairs(t *testing.T, tokenPairs []dbModule.Pair, sOrder *types.Order) {
	if len(tokenPairs) != 1 {
		t.Errorf("Expected 1 value, got %v", len(tokenPairs))
		return
	}
	if !reflect.DeepEqual(tokenPairs[0].TokenA, sOrder.TakerToken) {
		t.Errorf("Expected %#x, got %#x", sOrder.TakerToken[:], tokenPairs[0].TokenA[:])
	}
	if !reflect.DeepEqual(tokenPairs[0].TokenB, sOrder.MakerToken) {
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
	defer func() {
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
	defer func() {
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
	defer func() {
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
	defer func() {
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

func TestMarshalPairs(t *testing.T) {
	sOrder := sampleOrder()
	pair := &dbModule.Pair{sOrder.MakerToken, sOrder.TakerToken}
	pairJSON, _ := json.Marshal(pair)
	if string(pairJSON) != "{\"tokenA\":{\"address\":\"0x1dad4783cf3fe3085c1426157ab175a6119a04ba\",\"minAmount\":\"1\",\"maxAmount\":\"115792089237316195423570985008687907853269984665640564039457584007913129639935\",\"precision\":5},\"tokenB\":{\"address\":\"0x05d090b51c40b020eab3bfcb6a2dff130df22e9c\",\"minAmount\":\"1\",\"maxAmount\":\"115792089237316195423570985008687907853269984665640564039457584007913129639935\",\"precision\":5}}" {
		t.Errorf("Unexpected response, got '%v'", string(pairJSON))
	}
}
