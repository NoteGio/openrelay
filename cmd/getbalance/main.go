package main

import (
	"os"
	"github.com/notegio/openrelay/funds"
	"log"
	"encoding/hex"
	"fmt"
)

func hexToBytes(address string) [20]byte {
	slice, err := hex.DecodeString(address)
	if err != nil {
		return [20]byte{}
	}
	output := [20]byte{}
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
