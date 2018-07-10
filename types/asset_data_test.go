package types_test

import (
	"encoding/hex"
	"encoding/json"
	"github.com/notegio/openrelay/types"
	"testing"
	"bytes"
)

func getAssetData() (types.AssetData) {
	assetData := make(types.AssetData, 37)
	data, _ := hex.DecodeString("f47261b00000000000000000000000006dfff22588be9b3ef8cf0ad6dc9b84796f9fb45f01")
	copy(assetData[:], data[:])
	return assetData
}

func TestType(t *testing.T) {
	if proxyId := getAssetData().ProxyId(); !bytes.Equal(proxyId[:], types.ERC20ProxyID[:]) {
		t.Errorf("Unexpected ProxyId: %#x", proxyId)
	}
}

func TestAssetDataJSONEncode(t *testing.T) {
	data, err := json.Marshal(getAssetData())
	if err != nil {
		t.Errorf("Error marshalling data: %v", err.Error())
	}
	if string(data) != "\"0xf47261b00000000000000000000000006dfff22588be9b3ef8cf0ad6dc9b84796f9fb45f01\"" {
		t.Errorf("Unexpected JSON data: %v", string(data))
	}
}

// func TestAssetDataJSONDecode(t *testing.T) {
// 	assetData := types.AssetData{}
// 	err := json.Unmarshal(assetData, []byte("\"0xf47261b00000000000000000000000006dfff22588be9b3ef8cf0ad6dc9b84796f9fb45f01\""))
// 	if err != nil {
// 		t.Errorf("Error unmarshalling data: %v", err.Error())
// 	}
// 	if proxyId := assetData.ProxyId(); !bytes.Equal(proxyId[:], types.ERC20ProxyID[:]) {
// 		t.Errorf("Unexpected ProxyId: %#x", proxyId)
// 	}
// }
