package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"reflect"
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

func (sig *Signature) Verify(address [20]byte) bool {
	sigBytes := make([]byte, 65)
	copy(sigBytes[32-len(sig.R):32], sig.R[:])
	copy(sigBytes[64-len(sig.S):64], sig.S[:])
	sigBytes[64] = byte(int(sig.V)-27)
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
