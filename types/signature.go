package types

import (
	// "github.com/jinzhu/gorm"
	// "database/sql/driver"
	// "encoding/json"
	// "errors"
	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/crypto"
	// "log"
	// "reflect"
)

type Signature []byte

func (sig *Signature) Verify(address *Address, hash []byte) bool {
	return true
	// sigValue, _ := sig.Value()
	// sigBytes := sigValue.([]byte)
	//
	// hashedBytes := append([]byte("\x19Ethereum Signed Message:\n32"), hash[:]...)
	// signedBytes := crypto.Keccak256(hashedBytes)
	// pub, err := crypto.Ecrecover(signedBytes, sigBytes)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return false
	// }
	// recoverAddress := common.BytesToAddress(crypto.Keccak256(pub[1:])[12:])
	// return reflect.DeepEqual(address[:], recoverAddress[:])
}
