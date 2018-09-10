package types_test

import (
	"github.com/notegio/openrelay/types"
	"testing"
	"log"
)

func TestVerifyEthSig(t *testing.T) {
	order := &types.Order{}
	order.Initialize()
	exchangeAddressBytes, err := types.HexStringToBytes("b69e673309512a9d726f87304c6984054f87a93b")
	if err != nil { t.Errorf(err.Error()) }
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	signature := types.Signature{}
	// Signature generated with ganache CLI given order hash
	sigBytes, err := types.HexStringToBytes("1c41c1957c20f21e35f8b611dd9d1697cfa7b5c591271912534a9f4d7d1eeef01f368fc4818bb472f7450913897cea39e8f07fdacc5284286f10b6cf5eed02d2a4")
	if err != nil { t.Errorf(err.Error()) }
	signature = append(signature, sigBytes...)
	// Default account for ganache CLI
	signerAddressBytes, err := types.HexStringToBytes("627306090abab3a6e1400e9345bc60c78a8bef57")
	if err != nil { t.Errorf(err.Error()) }
	signerAddress := &types.Address{}
	copy(signerAddress[:], signerAddressBytes[:])
	signature = append(signature, types.SigTypeEthSign)
	log.Printf("%#x", order.Hash)
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
	sigBytes, err := types.HexStringToBytes("1b34f5afe988ba953bc074c7c54f6bfa5b9dc1e06e3be5a0b4993e65c1429775445ef3b937e25e16643e621f6b57331b4082b974033785a3045915e90a77077200")
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
