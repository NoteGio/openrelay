package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"database/sql/driver"
	"log"
	"reflect"
	"errors"
)

type Signature struct {
	V    byte
	R    [32]byte
	S    [32]byte
	Hash [32]byte
}

type jsonSignature struct {
	V    json.Number `json:"v"`
	R    string      `json:"r"`
	S    string      `json:"s"`
	Hash string
}

func (sig *Signature) Verify(address *Address) bool {
	sigValue, _ := sig.Value()
	sigBytes := sigValue.([]byte)

	hashedBytes := append([]byte("\x19Ethereum Signed Message:\n32"), sig.Hash[:]...)
	signedBytes := crypto.Keccak256(hashedBytes)
	pub, err := crypto.Ecrecover(signedBytes, sigBytes)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	recoverAddress := common.BytesToAddress(crypto.Keccak256(pub[1:])[12:])
	return reflect.DeepEqual(address[:], recoverAddress[:])
}

func (sig *Signature) Value() (driver.Value, error) {
	sigBytes := make([]byte, 65)
	copy(sigBytes[32-len(sig.R):32], sig.R[:])
	copy(sigBytes[64-len(sig.S):64], sig.S[:])
	sigBytes[64] = byte(int(sig.V) - 27)
	return sigBytes, nil
}

func (sig *Signature) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		if len(v) != 65 {
			return errors.New("Signature scanner src should be []byte of length 65")
		}
		copy(sig.R[:], v[0:32])
		copy(sig.S[:], v[32:64])
		sig.V = byte(int(v[64]) + 27)
		return nil
	default:
		return errors.New("Signature scanner src should be []byte of length 65")
	}
}
