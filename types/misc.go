package types

import (
	"math/big"
	"database/sql/driver"
	"errors"
	"fmt"
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
