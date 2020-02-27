package types_test

import (
	"github.com/notegio/openrelay/types"
	"testing"
	// "log"
)

func TestVerifyEthSig(t *testing.T) {
	order := &types.Order{}
	order.Initialize()
	exchangeAddressBytes, err := types.HexStringToBytes("90fe2af704b34e0224bf2299c838e04d4dcf1364")
	if err != nil { t.Errorf(err.Error()) }
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	signature := types.Signature{}
	sigBytes, err := types.HexStringToBytes("1c64f11dfb26a787964551b7507854a3db9431438a44d922c761982f14d32056e61ce1b466208375a623a68cf3cb2155d6316ad4abb24b31af2761bb2c77737d55")
	if err != nil { t.Errorf(err.Error()) }
	signature = append(signature, sigBytes...)
	signerAddressBytes, err := types.HexStringToBytes("a2123be740a0d348cab7f6d6c44d02c13d4b85a3")
	if err != nil { t.Errorf(err.Error()) }
	copy(order.Maker[:], signerAddressBytes)
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
	exchangeAddressBytes, err := types.HexStringToBytes("90fe2af704b34e0224bf2299c838e04d4dcf1364")
	if err != nil { t.Errorf(err.Error()) }
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	signature := types.Signature{}
	sigBytes, err := types.HexStringToBytes("1c3c5fdf4d451ff39f97834e29f38d093184eeedfbc4c8ffbd304851cc54fdee4035055e0cae0cea0bbbeeebae4a59fafd159194c4608dc8de34c26caa60e3b94c")
	if err != nil { t.Errorf(err.Error()) }
	signature = append(signature, sigBytes...)
	signerAddressBytes, err := types.HexStringToBytes("192e2d8c3d38f2c63917e19f829402ffc4ce0538")
	if err != nil { t.Errorf(err.Error()) }
	copy(order.Maker[:], signerAddressBytes)
	signerAddress := &types.Address{}
	copy(signerAddress[:], signerAddressBytes[:])
	signature = append(signature, types.SigTypeEIP712)
	if !signature.Verify(signerAddress, order.Hash()) {
		t.Fatalf("Signature invalid: %#x", signature[:])
	}
}
