package types

import (
	"bytes"
	// "log"
)

type AssetData []byte

var ERC20ProxyID = [4]byte{244, 114, 97, 176}
var ERC721ProxyID = [4]byte{8, 233, 55, 250}

func (data AssetData) ProxyId() ([4]byte) {
	result := [4]byte{}
	copy(result[:], data[0:4])
	return result
}

func (data AssetData) Address() (*Address) {
	address := &Address{}
	if data.IsType(ERC20ProxyID) || data.IsType(ERC721ProxyID) {
		copy(address[:], data[16:36])
	}
	return address
}

func (data AssetData) IsType(proxyId [4]byte) (bool) {
	return bytes.Equal(data[0:4], proxyId[:])
}

func (data AssetData) SupportedType() (bool) {
	return data.IsType(ERC20ProxyID) //|| data.IsType(ERC721ProxyID)
}
