package types

import (
	// "encoding/hex"
	"bytes"
	"fmt"
	"database/sql/driver"
	// "strings"
	// "log"
)

type AssetData []byte

var ERC20ProxyID = [4]byte{244, 114, 97, 176} // 0xf47261b0
var ERC721ProxyID = [4]byte{2, 87, 23, 146} // 0x02571792

func (data AssetData) ProxyId() ([4]byte) {
	result := [4]byte{}
	if len(data) >= 4 {
		copy(result[:], data[0:4])
	}
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
	return data.IsType(ERC20ProxyID) || data.IsType(ERC721ProxyID)
}

func (data AssetData) TokenID() (*Uint256) {
	tokenID := &Uint256{}
	if data.IsType(ERC721ProxyID) {
		copy(tokenID[:], data[36:])
	}
	return tokenID
}

func (data AssetData) MarshalJSON() ([]byte, error) {
	if len(data) == 0 {
		return []byte(`"0x"`), nil
	}
	return []byte(fmt.Sprintf("\"%#x\"", data[:])), nil
}

func (data AssetData) Value() (driver.Value, error) {
	return []byte(data[:]), nil
}
