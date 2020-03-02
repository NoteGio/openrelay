package metadata_test

import (
	"bytes"
	"bufio"
	"context"
	"fmt"
	"math/big"
	"github.com/jinzhu/gorm"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/metadata"
	orCommon "github.com/notegio/openrelay/common"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	// "github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/common"
	"testing"
	"os"
	"net/http"
	// "log"
)
type MockContractBackend struct {
	idCounter int
}
func (mock *MockContractBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return []byte{}, nil
}
func (mock *MockContractBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	mock.idCounter += 1
	result := []byte(fmt.Sprintf("http://test/%v", mock.idCounter))
	offset := orCommon.BigToUint256(big.NewInt(32))
	length := orCommon.BigToUint256(big.NewInt(int64(len(result))))
	return append(offset[:], append(length[:], append(result, make([]byte, 32 - (len(result) % 32))...)...)...), nil
}
func (mock *MockContractBackend) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return []byte{}, nil
}
func (mock *MockContractBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return 0, nil
}
func (mock *MockContractBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return new(big.Int), nil
}
func (mock *MockContractBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	return 0, nil
}
func (mock *MockContractBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return nil
}
func (mock *MockContractBackend) FilterLogs(ctx context.Context, filter ethereum.FilterQuery) ([]types.Log, error) {
	return []types.Log{}, nil
}
func (mock *MockContractBackend) SubscribeFilterLogs(ctx context.Context, filter ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return nil, nil
}

type MockHttpClient struct {
	results map[string][]byte
}

func (client *MockHttpClient) Get(url string) (*http.Response, error) {
	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(client.results[url])), nil)
	return resp, err
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

func TestProcessAssetData(t *testing.T) {
	caller := &MockContractBackend{}
	client := &MockHttpClient{make(map[string][]byte)}
	db, err := getDb()
	if err != nil {
		t.Fatalf(err.Error())
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.AssetMetadata{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	if err := tx.AutoMigrate(&dbModule.AssetAttribute{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	client.results["http://test/1"] = []byte(`HTTP/1.1 200 OK
Content-Type: application/json

{"description":"A big cow","external_url":"https://ethercow.openrelay.xyz/cows/1","image":"https://ethercow.openrelay.xyz/cows/1.png","name":"Jeff","attributes":{"base":"cow","level":5,"weight":500,"generation":1}}`)
	metadataConsumer, _ := metadata.NewRawOrderMetadataConsumer(caller, client, tx, 1)
	assetData, _ := orCommon.HexToAssetData("0x02571792000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498")
	metadataConsumer.ProcessAssetData(&assetData)
	count := 0
	assetMetadata := &dbModule.AssetMetadata{}
	tx.Model(&dbModule.AssetMetadata{}).Count(&count)
	if count != 1 {
		t.Errorf("count: expected 1, got %v", count)
	}
	if err := tx.Model(&dbModule.AssetMetadata{}).First(assetMetadata).Error; err != nil {
		t.Errorf(err.Error())
	}
	if tx.Model(&dbModule.AssetAttribute{}).Count(&count); count != 4 {
		t.Errorf("Expected 4 attributes, have %v", count)
	}
	if assetMetadata.Name != "Jeff" {
		t.Errorf("Unexpected name ''%v'", assetMetadata.Name)
	}
	if assetMetadata.Description != "A big cow" {
		t.Errorf("Unexpected Description: '%v'", assetMetadata.Description)
	}
	if assetMetadata.ExternalURL != "https://ethercow.openrelay.xyz/cows/1" {
		t.Errorf("Unexpected ExternalURL: '%v'", assetMetadata.ExternalURL)
	}
	if assetMetadata.Image != "https://ethercow.openrelay.xyz/cows/1.png" {
		t.Errorf("Unexpected Image: '%v'", assetMetadata.Image)
	}
	if assetMetadata.BackgroundColor != "" {
		t.Errorf("Unexpected BackgroundColor: '%v'", assetMetadata.BackgroundColor)
	}
		if !bytes.Equal(assetData[:], assetMetadata.AssetData[:]) {
		t.Errorf("Asset data does not match: %#x != %#x", assetData[:], assetMetadata.AssetData[:])
	}
	if len(assetMetadata.Attributes) != 4 {
		t.Fatalf("Unexpected attribute count %v", len(assetMetadata.Attributes))
	}
	testAttributes := []dbModule.AssetAttribute{
		dbModule.AssetAttribute{Name: "base", Type: "string", Value: "cow", DisplayType: ""},
		dbModule.AssetAttribute{Name: "generation", Type: "number", Value: "1.000000", DisplayType: ""},
		dbModule.AssetAttribute{Name: "level", Type: "number", Value: "5.000000", DisplayType: ""},
		dbModule.AssetAttribute{Name: "weight", Type: "number", Value: "500.000000", DisplayType: ""},
	}
	for i, val := range testAttributes {
		if assetMetadata.Attributes[i].Name != val.Name {
			t.Errorf("Unexpected Name for attribute %v: %v != %v", i, assetMetadata.Attributes[i].Name, val.Name)
		}
		if assetMetadata.Attributes[i].Type != val.Type {
			t.Errorf("Unexpected Type for attribute %v: %v != %v", i, assetMetadata.Attributes[i].Type, val.Type)
		}
		if assetMetadata.Attributes[i].Value != val.Value {
			t.Errorf("Unexpected Value for attribute %v: %v != %v", i, assetMetadata.Attributes[i].Value, val.Value)
		}
		if assetMetadata.Attributes[i].DisplayType != val.DisplayType {
			t.Errorf("Unexpected DisplayType for attribute %v: %v != %v", i, assetMetadata.Attributes[i].DisplayType, val.DisplayType)
		}
		if !bytes.Equal(assetMetadata.Attributes[i].AssetData[:], assetMetadata.AssetData[:]) {
			t.Errorf("Asset data does not match for attribute %v: %#x != %#x", i, assetMetadata.Attributes[i].AssetData[:], assetMetadata.AssetData[:])
		}
	}
}
func TestProcessAssetDataMalformed(t *testing.T) {
	caller := &MockContractBackend{}
	client := &MockHttpClient{make(map[string][]byte)}
	db, err := getDb()
	if err != nil {
		t.Fatalf(err.Error())
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
		db.Close()
	}()
	if err := tx.AutoMigrate(&dbModule.AssetMetadata{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	if err := tx.AutoMigrate(&dbModule.AssetAttribute{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	client.results["http://test/1"] = []byte(`HTTP/1.1 200 OK
Content-Type: application/json

<html></html>`)
	metadataConsumer, _ := metadata.NewRawOrderMetadataConsumer(caller, client, tx, 1)
	assetData, _ := orCommon.HexToAssetData("0x02571792000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498")
	metadataConsumer.ProcessAssetData(&assetData)
	count := 0
	assetMetadata := &dbModule.AssetMetadata{}
	tx.Model(&dbModule.AssetMetadata{}).Count(&count)
	if count != 1 {
		t.Errorf("count: expected 1, got %v", count)
	}
	if err := tx.Model(&dbModule.AssetMetadata{}).First(assetMetadata).Error; err != nil {
		t.Errorf(err.Error())
	}
	if assetMetadata.RawMetadata != "<html></html>" {
		t.Errorf("Expected metadata to be html, got: %v", assetMetadata.RawMetadata)
	}
}
