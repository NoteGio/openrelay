package types_test

import (
	"encoding/json"
	"github.com/notegio/openrelay/types"
	"testing"
)

func getUint256() (*types.Uint256) {
	return &types.Uint256{}
}

func getAddress() (*types.Address) {
	return &types.Address{}
}

func TestUint256JSONEncode(t *testing.T) {
	data, err := json.Marshal(getUint256())
	if err != nil {
		t.Errorf("Error marshalling data: %v", err.Error())
	}
	if string(data) != "\"0\"" {
		t.Errorf("Unexpected JSON data: %v", string(data))
	}
}

func TestAddressJSONEncode(t *testing.T) {
	data, err := json.Marshal(getAddress())
	if err != nil {
		t.Errorf("Error marshalling data: %v", err.Error())
	}
	if string(data) != "\"0x0000000000000000000000000000000000000000\"" {
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
