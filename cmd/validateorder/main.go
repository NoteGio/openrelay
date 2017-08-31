package main

import (
	"os"
	"github.com/notegio/openrelay/funds"
	"github.com/notegio/openrelay/types"
	"log"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	orderFile := os.Args[2]

	fundChecker, err := funds.NewRpcOrderValidator(rpcURL)
	if err != nil {
		log.Fatalf(err.Error())
	}
	newOrder := types.Order{}
	if orderData, err := ioutil.ReadFile(orderFile); err == nil {
		if err := json.Unmarshal(orderData, &newOrder); err != nil {
			log.Fatalf(err.Error())
		}
	}
	fmt.Println(fundChecker.ValidateOrder(&newOrder))
}
