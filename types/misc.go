package types

import (
	"encoding/hex"
	"math/big"
	"database/sql/driver"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"errors"
	"fmt"
	"bytes"
)

type Address [20]byte

func (addr *Address) Value() (driver.Value, error) {
	return addr[:], nil
}

func (addr *Address) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		copy(addr[:], v)
		return nil
	default:
		return errors.New("Address scanner src should be []byte")
	}
}

func (addr *Address) String() string {
	return fmt.Sprintf("%#x", addr[:])
}

func (addr *Address) UnmarshalJSON(data []byte) error {
	length, err := hex.Decode(addr[:], bytes.TrimPrefix(bytes.Trim(data, "\""), []byte("0x")))
	if err != nil { return err }
	if length != 20 {
		return errors.New("Invalid address length")
	}
	return nil
}


func (data *Address) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", data)), nil
}

func (data *Address) ToGethAddress() (common.Address) {
	return common.HexToAddress(hex.EncodeToString(data[:]))
}

type Uint256 [32]byte

func (data *Uint256) Value() (driver.Value, error) {
	return data[:], nil
}

func (data *Uint256) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		copy(data[:], v)
		return nil
	default:
		return errors.New("Uint256 scanner src should be []byte")
	}
}

func (data *Uint256) String() (string) {
	return data.Big().String()
}

func (data *Uint256) Big() (*big.Int) {
	return new(big.Int).SetBytes(data[:])
}

func (data *Uint256) Uint() (uint) {
	return uint(new(big.Int).SetBytes(data[:]).Uint64())
}

func (data *Uint256) UnmarshalJSON(jsonData []byte) error {
	numberBytes := bytes.Trim(jsonData, "\"")
	numberBig, ok := new(big.Int).SetString(string(numberBytes), 10)
	if !ok { return fmt.Errorf("Failed to convert number to integer: %v", string(jsonData)) }
	copy(data[:], abi.U256(numberBig))
	return nil
}

func (data *Uint256) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", data)), nil
}
