package common

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/notegio/openrelay/types"
	"strings"
)

func BytesToAddress(data [20]byte) common.Address {
	return common.HexToAddress(hex.EncodeToString(data[:]))
}

func ToGethAddress(data *types.Address) common.Address {
	return common.HexToAddress(hex.EncodeToString(data[:]))
}

func HexToBytes(hexString string) ([20]byte, error) {
	slice, err := hex.DecodeString(strings.TrimPrefix(hexString))
	result := [20]byte{}
	if err != nil {
		return result, err
	}
	copy(result[:], slice[:])
	return result, nil
}
