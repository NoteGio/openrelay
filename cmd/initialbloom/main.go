package main

import (
	"github.com/notegio/openrelay/fillbloom"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"os"
	"log"
	"strconv"
)

func main() {
	rpcURL := os.Args[1]
	conn, err := ethclient.Dial(rpcURL)
	if err != nil { log.Panicf(err.Error()); }
	_, err = strconv.Atoi(os.Args[2])
	if err != nil { log.Panicf(err.Error()); }
	endBlock, err := strconv.Atoi(os.Args[3])
	if err != nil { log.Panicf(err.Error()); }
	exchangeAddress := common.HexToAddress(os.Args[4])
	bloom, err := fillbloom.NewFillBloom(os.Args[5])
	if err != nil { log.Panicf(err.Error()); }
	err = bloom.Initialize(conn, int64(endBlock), []common.Address{exchangeAddress})
	if err != nil { log.Panicf(err.Error()); }
	if err = bloom.Save(); err != nil {
		log.Panicf(err.Error());
	}
}
