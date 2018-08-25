package main

import (
	"os"
	"github.com/notegio/openrelay/funds/balance"
	"github.com/notegio/openrelay/types"
	"log"
	"encoding/hex"
	"fmt"
	"strings"
)

func hexToAddress(address string) *types.Address {
	slice, err := hex.DecodeString(address)
	if err != nil {
		return &types.Address{}
	}
	output := &types.Address{}
	copy(output[:], slice[:])
	return output
}

func hexToAssetData(address string) types.AssetData {
	address = strings.Trim(address, "0x")
	if strings.HasPrefix(address, "f47261b0") {
		address = "f47261b0000000000000000000000000" + address
	}
	slice, err := hex.DecodeString(address)
	if err != nil {
		return types.AssetData{}
	}
	output := make(types.AssetData, len(slice))
	copy(output[:], slice[:])
	return output
}

func main() {
	rpcURL := os.Args[1]
	tokenAddress := os.Args[2]
	userAddress := os.Args[3]

	fundChecker, err := balance.NewRpcRoutingBalanceChecker(rpcURL)
	if err != nil {
		log.Fatalf(err.Error())
	}
	balance, err := fundChecker.GetBalance(hexToAssetData(tokenAddress), hexToAddress(userAddress))
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("'%v'\n", balance)
}
