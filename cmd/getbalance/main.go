package main

import (
	"os"
	"github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/types"
	"log"
	"encoding/hex"
	"fmt"
)

func hexToBytes(address string) *types.Address {
	slice, err := hex.DecodeString(address)
	if err != nil {
		return &types.Address{}
	}
	output := &types.Address{}
	copy(output[:], slice[:])
	return output
}

func main() {
	rpcURL := os.Args[1]
	tokenAddress := os.Args[2]
	userAddress := os.Args[3]

	fundChecker, err := funds.NewRpcBalanceChecker(rpcURL)
	if err != nil {
		log.Fatalf(err.Error())
	}
	balance, err := fundChecker.GetBalance(hexToBytes(tokenAddress), hexToBytes(userAddress))
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("'%v'\n", balance)
}
