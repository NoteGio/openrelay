package db_test

import (
	"bytes"
	"encoding/json"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/common"
	"testing"
	// "log"
)

func VerifyAsset(asset *dbModule.AssetMetadata, jsonBytes []byte, rawAttributes string, assetData types.AssetData, t *testing.T) {
	if asset == nil {
		t.Fatalf("AssetMetadata is nil")
	}
	if asset.Name != "Jeff" {
		t.Errorf("Unexpected name '%v'", asset.Name)
	}
	if asset.Description != "A big cow" {
		t.Errorf("Unexpected Description: '%v'", asset.Description)
	}
	if asset.ExternalURL != "https://ethercow.openrelay.xyz/cows/3" {
		t.Errorf("Unexpected ExternalURL: '%v'", asset.ExternalURL)
	}
	if asset.Image != "https://ethercow.openrelay.xyz/cows/3.png" {
		t.Errorf("Unexpected Image: '%v'", asset.Image)
	}
	if asset.BackgroundColor != "" {
		t.Errorf("Unexpected BackgroundColor: '%v'", asset.BackgroundColor)
	}
	if asset.RawMetadata != string(jsonBytes) {
		t.Errorf("Unexpected RawMetadata: '%v' != '%v'", asset.RawMetadata, string(jsonBytes))
	}
	if asset.RawAttributes != rawAttributes {
		t.Errorf("Unexpected RawAttributes: '%v'", asset.RawAttributes)
	}
	if !bytes.Equal(assetData[:], asset.AssetData[:]) {
		t.Errorf("Asset data does not match: %#x != %#x", assetData[:], asset.AssetData[:])
	}
	if len(asset.Attributes) != 4 {
		t.Fatalf("Unexpected attribute count %v", len(asset.Attributes))
	}
}

func VerifyAssetAttributes(asset *dbModule.AssetMetadata, t *testing.T) {
	testAttributes := []dbModule.AssetAttribute{
		dbModule.AssetAttribute{Name: "base", Type: "string", Value: "cow", DisplayType: ""},
		dbModule.AssetAttribute{Name: "generation", Type: "number", Value: "2.000000", DisplayType: ""},
		dbModule.AssetAttribute{Name: "level", Type: "number", Value: "5.000000", DisplayType: ""},
		dbModule.AssetAttribute{Name: "weight", Type: "number", Value: "500.000000", DisplayType: ""},
	}
	for i, val := range testAttributes {
		if asset.Attributes[i].Name != val.Name {
			t.Errorf("Unexpected Name for attribute %v: %v != %v", i, asset.Attributes[i].Name, val.Name)
		}
		if asset.Attributes[i].Type != val.Type {
			t.Errorf("Unexpected Type for attribute %v: %v != %v", i, asset.Attributes[i].Type, val.Type)
		}
		if asset.Attributes[i].Value != val.Value {
			t.Errorf("Unexpected Value for attribute %v: %v != %v", i, asset.Attributes[i].Value, val.Value)
		}
		if asset.Attributes[i].DisplayType != val.DisplayType {
			t.Errorf("Unexpected DisplayType for attribute %v: %v != %v", i, asset.Attributes[i].DisplayType, val.DisplayType)
		}
		if !bytes.Equal(asset.Attributes[i].AssetData[:], asset.AssetData[:]) {
			t.Errorf("Asset data does not match for attribute %v: %#x != %#x", i, asset.Attributes[i].AssetData[:], asset.AssetData[:])
		}
	}
}

func TestParseAssetMetadata(t *testing.T) {
	rawAttributes := `[{"trait_type":"base","value":"cow"},{"trait_type":"level","value":5},{"trait_type":"weight","value":500},{"display_type":"number","trait_type":"generation","value":2}]`
	jsonBytes := []byte(`{
		  "description": "A big cow",
		  "external_url": "https://ethercow.openrelay.xyz/cows/3",
		  "image": "https://ethercow.openrelay.xyz/cows/3.png",
		  "name": "Jeff",
		  "attributes": [
			  {
		      "trait_type": "base",
		      "value": "cow"
		    },
		    {
		      "trait_type": "level",
		      "value": 5
		    },
		    {
		      "trait_type": "weight",
		      "value": 500
		    },
		    {
		      "display_type": "number",
		      "trait_type": "generation",
		      "value": 2
		    }]
		}`)
	asset := &dbModule.AssetMetadata{}
	if err := json.Unmarshal(jsonBytes, asset); err != nil {
		t.Fatalf(err.Error())
	}
	VerifyAsset(asset, jsonBytes, rawAttributes, []byte{}, t)
	testAttributes := []dbModule.AssetAttribute{
		dbModule.AssetAttribute{Name: "base", Type: "string", Value: "cow", DisplayType: ""},
		dbModule.AssetAttribute{Name: "level", Type: "number", Value: "5.000000", DisplayType: ""},
		dbModule.AssetAttribute{Name: "weight", Type: "number", Value: "500.000000", DisplayType: ""},
		dbModule.AssetAttribute{Name: "generation", Type: "number", Value: "2.000000", DisplayType: "number"},
	}
	for i, val := range testAttributes {
		if asset.Attributes[i].Name != val.Name {
			t.Errorf("Unexpected Name for attribute %v: %v", i, asset.Attributes[i].Name)
		}
		if asset.Attributes[i].Type != val.Type {
			t.Errorf("Unexpected Type for attribute %v: %v", i, asset.Attributes[i].Type)
		}
		if asset.Attributes[i].Value != val.Value {
			t.Errorf("Unexpected Value for attribute %v: %v", i, asset.Attributes[i].Value)
		}
		if asset.Attributes[i].DisplayType != val.DisplayType {
			t.Errorf("Unexpected DisplayType for attribute %v: %v", i, asset.Attributes[i].DisplayType)
		}
	}
}
func TestParseAssetMetadataDictAtributes(t *testing.T) {
	rawAttributes := `{"base":"cow","generation":2,"level":5,"weight":500}`
	jsonBytes := []byte(`{
		  "description": "A big cow",
		  "external_url": "https://ethercow.openrelay.xyz/cows/3",
		  "image": "https://ethercow.openrelay.xyz/cows/3.png",
		  "name": "Jeff",
		  "attributes": {
				"base": "cow",
				"level": 5,
				"weight": 500,
				"generation": 2
			}
		}`)
	assetData, _ := common.HexToAssetData("0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498")
	asset := &dbModule.AssetMetadata{}
	if err := json.Unmarshal(jsonBytes, asset); err != nil {
		t.Fatalf(err.Error())
	}
	asset.SetAssetData(assetData)
	VerifyAsset(asset, jsonBytes, rawAttributes, assetData, t)
	VerifyAssetAttributes(asset, t)
}

func TestStoreAssetMetadata(t *testing.T) {
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
	rawAttributes := `{"base":"cow","generation":2,"level":5,"weight":500}`
	jsonBytes := []byte(`{
		  "description": "A big cow",
		  "external_url": "https://ethercow.openrelay.xyz/cows/3",
		  "image": "https://ethercow.openrelay.xyz/cows/3.png",
		  "name": "Jeff",
		  "attributes": {
				"base": "cow",
				"level": 5,
				"weight": 500,
				"generation": 2
			}
		}`)
	assetData, _ := common.HexToAssetData("0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498")
	asset := &dbModule.AssetMetadata{}
	if err := json.Unmarshal(jsonBytes, asset); err != nil {
		t.Fatalf(err.Error())
	}
	asset.SetAssetData(assetData)
	if err := tx.AutoMigrate(&dbModule.AssetMetadata{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	if err := tx.AutoMigrate(&dbModule.AssetAttribute{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	if err := tx.Model(&dbModule.AssetMetadata{}).Create(asset).Error; err != nil {
		t.Fatalf(err.Error())
	}

	dbAsset := dbModule.AssetMetadata{}
	if err := tx.Debug().Model(&dbModule.AssetMetadata{}).Where("asset_data = ?", &assetData).First(&dbAsset).Error; err != nil {
		t.Fatalf(err.Error())
	}
	VerifyAsset(&dbAsset, jsonBytes, rawAttributes, assetData, t)
	VerifyAssetAttributes(&dbAsset, t)
}

func TestPopulateOrderMetadata(t *testing.T) {
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
	rawAttributes := `{"base":"cow","generation":2,"level":5,"weight":500}`
	jsonBytes := []byte(`{
		  "description": "A big cow",
		  "external_url": "https://ethercow.openrelay.xyz/cows/3",
		  "image": "https://ethercow.openrelay.xyz/cows/3.png",
		  "name": "Jeff",
		  "attributes": {
				"base": "cow",
				"level": 5,
				"weight": 500,
				"generation": 2
			}
		}`)
	assetData, _ := common.HexToAssetData("0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498")
	otherAssetData, _ := common.HexToAssetData("0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f499")
	asset := &dbModule.AssetMetadata{}
	if err := json.Unmarshal(jsonBytes, asset); err != nil {
		t.Fatalf(err.Error())
	}
	asset.SetAssetData(assetData)
	if err := tx.AutoMigrate(&dbModule.AssetMetadata{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	if err := tx.AutoMigrate(&dbModule.AssetAttribute{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	if err := tx.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		t.Errorf(err.Error())
	}
	if err := tx.Model(&dbModule.AssetMetadata{}).Create(asset).Error; err != nil {
		t.Fatalf(err.Error())
	}
	asset.Name = "Joe"
	asset.SetAssetData(otherAssetData)
	if err := tx.Model(&dbModule.AssetMetadata{}).Create(asset).Error; err != nil {
		t.Fatalf(err.Error())
	}

	orders := []dbModule.Order{dbModule.Order{}, dbModule.Order{}, dbModule.Order{}}
	orders[0].Initialize()
	orders[1].Initialize()
	orders[2].Initialize()
	orders[0].MakerAssetData = assetData
	orders[1].TakerAssetData = assetData
	orders[2].MakerAssetData = otherAssetData
	dbModule.PopulateAssetMetadata(orders, tx)
	if orders[2].MakerAssetMetadata.Name != "Joe" {
		t.Errorf("Got unexpected metadata for asset 2")
	}
	if count := len(orders[2].MakerAssetMetadata.Attributes); count != 4 {
		t.Errorf("Got unexpected attribute count %v for asset 2", count)
	}

	dbAssetMetadata := []*dbModule.AssetMetadata{orders[0].MakerAssetMetadata, orders[1].TakerAssetMetadata}
	for j, dbAsset := range dbAssetMetadata {
		if dbAsset == nil {
			t.Fatalf("AssetMetadata is nil for iteration %v", j)
		}
		VerifyAsset(dbAsset, jsonBytes, rawAttributes, assetData, t)
		VerifyAssetAttributes(dbAsset, t)
	}
}
