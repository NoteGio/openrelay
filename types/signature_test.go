package types_test

import (
	// "database/sql"
	// "database/sql/driver"
	// "encoding/hex"
	// "encoding/json"
	"github.com/notegio/openrelay/types"
	// "io/ioutil"
	// "reflect"
	"testing"
	// "bytes"
	// "log"
)

func TestVerifyEthSig(t *testing.T) {
	order := &types.Order{}
	order.Initialize()
	exchangeAddressBytes, err := types.HexStringToBytes("b69e673309512a9d726f87304c6984054f87a93b")
	if err != nil { t.Errorf(err.Error()) }
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	signature := types.Signature{}
	// Signature generated with ganache CLI given order hash
	sigBytes, err := types.HexStringToBytes("006bcc503876436ae6ebddecc16f95fdc74945ba85aa7debabdfa4a708a80b0272520d4f331a50396583db9a06bce884abc82219bfe180ef0093b0534786c996c2")
	if err != nil { t.Errorf(err.Error()) }
	signature = append(signature, sigBytes...)
	// Default account for ganache CLI
	signerAddressBytes, err := types.HexStringToBytes("627306090abab3a6e1400e9345bc60c78a8bef57")
	if err != nil { t.Errorf(err.Error()) }
	signerAddress := &types.Address{}
	copy(signerAddress[:], signerAddressBytes[:])
	signature = append(signature, types.SigTypeEthSign)
	if !signature.Verify(signerAddress, order.Hash()) {
		t.Fatalf("Signature invalid: %#x", signature[:])
	}
}

func TestVerifyEIP712Sig(t *testing.T) {
	order := &types.Order{}
	order.Initialize()
	exchangeAddressBytes, err := types.HexStringToBytes("b69e673309512a9d726f87304c6984054f87a93b")
	if err != nil { t.Errorf(err.Error()) }
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	signature := types.Signature{}
	// Signature generated with metamask
	sigBytes, err := types.HexStringToBytes("0138b836c68c21074797b0247652ef56bfe75f8e095544b60e19faaac2ee5592cb59f2988c1a37f3acbe6cfd41fecc71130c364a45cf34bb895addb8f0afea47d9")
	if err != nil { t.Errorf(err.Error()) }
	signature = append(signature, sigBytes...)
	// log.Printf("%#x", signature[:])
	signerAddressBytes, err := types.HexStringToBytes("627306090abab3a6e1400e9345bc60c78a8bef57")
	if err != nil { t.Errorf(err.Error()) }
	signerAddress := &types.Address{}
	copy(signerAddress[:], signerAddressBytes[:])
	signature = append(signature, types.SigTypeEIP712)
	if !signature.Verify(signerAddress, order.Hash()) {
		t.Fatalf("Signature invalid: %#x", signature[:])
	}
}
