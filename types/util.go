package types

import (
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"math/big"
	"strings"
)

func HexStringToBytes(hexString string) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(hexString, "0x"))
}

func intStringToBytes(intString string) ([]byte, error) {
	bigInt := new(big.Int)
	_, success := bigInt.SetString(intString, 10)
	if success {
		return abi.U256(bigInt), nil
	} else {
		return nil, errors.New("Value not a valid integer")
	}
}

func IntStringToUint256(intString string) (*Uint256, error) {
	data, err := intStringToBytes(intString)
	if err != nil {
		return nil, err
	}
	result := &Uint256{}
	copy(result[:], data[:])
	return result, nil
}
