package types

import (
	// "github.com/jinzhu/gorm"
	// "database/sql/driver"
	// "encoding/json"
	// "errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	// "reflect"
	"bytes"
)

const (
	SigTypeIllegal = 0
	SigTypeInvalid = 1
	SigTypeEIP712 = 2
	SigTypeEthSign = 3
	SigTypeCaller = 4
	SigTypeWallet = 5
	SigTypeValidator = 6
	SigTypePreSigned = 7
	SigTypeTrezor = 8
	NSigTypes = 9
)

type Signature []byte

func (sig Signature) Type() (byte) {
	return sig[len(sig[:])-1]
}

func (sig Signature) Verify(address *Address, hash []byte) bool {
	switch sigType := sig.Type(); sigType {
	case SigTypeEIP712:
		return sig.verifyEIP712(address, hash)
	case SigTypeEthSign:
		return sig.verifyEthSign(address, hash)
	case SigTypeWallet:
		return sig.verifyWallet(address, hash)
	case SigTypeValidator:
		return sig.verifyValidator(address, hash)
	case SigTypeTrezor:
		return sig.verifyTrezor(address, hash)
	default:
		return false
	}
	return false
}

func (sig Signature) verifyEIP712(address *Address, hash []byte) bool {
	if len(sig[:]) != 66 {
		log.Printf("Invalid length: %v", len(sig[:]))
		return false
	}
	v := sig[0]
	r := sig[1:33]
	s := sig[33:65]
	pub, err := crypto.Ecrecover(hash, append(append(r, s...), v))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	recoverAddress := common.BytesToAddress(crypto.Keccak256(pub[1:])[12:])
	return bytes.Equal(address[:], recoverAddress[:])
}

func (sig Signature) verifyEthSign(address *Address, hash []byte) bool {
	if len(sig[:]) != 66 {
		log.Printf("Invalid length: %v", len(sig[:]))
		return false
	}
	v := sig[0]
	r := sig[1:33]
	s := sig[33:65]
	hashedBytes := append([]byte("\x19Ethereum Signed Message:\n32"), hash[:]...)
	signedBytes := crypto.Keccak256(hashedBytes)
	pub, err := crypto.Ecrecover(signedBytes, append(append(r, s...), v))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	recoverAddress := common.BytesToAddress(crypto.Keccak256(pub[1:])[12:])
	return bytes.Equal(address[:], recoverAddress[:])
}

func (sig Signature) verifyWallet(address *Address, hash []byte) bool {
	// We may never support Wallet verification. Even if we check that an order
	// is verified as we ingest the order, each contract could have custom logic
	// that invalidates orders later, without any way for us to determine that
	// aside from checking every order on a repeating basis. That doesn't scale
	// well.
	return false
}

func (sig Signature) verifyValidator(address *Address, hash []byte) bool {
	// We don't currently support any validators, but we might support a
	// whitelist of validators in the future. To support a given validator, we
	// need to be sure that the validator provides adequate information that we
	// can monitor for events that would invalidate an order in a scalable
	// manner.
	return false
}

func (sig Signature) verifyTrezor(address *Address, hash []byte) bool {
	if len(sig[:]) != 66 {
		log.Printf("Invalid length: %v", len(sig[:]))
		return false
	}
	hashedBytes := append([]byte("\x19Ethereum Signed Message:\n\x20"), hash[:]...)
	signedBytes := crypto.Keccak256(hashedBytes)
	v := sig[0]
	r := sig[1:33]
	s := sig[33:65]
	pub, err := crypto.Ecrecover(signedBytes, append(append(r, s...), v))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	recoverAddress := common.BytesToAddress(crypto.Keccak256(pub[1:])[12:])
	return bytes.Equal(address[:], recoverAddress[:])
}
