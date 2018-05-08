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
	startBlock, err := strconv.Atoi(os.Args[2])
	if err != nil { log.Panicf(err.Error()); }
	endBlock, err := strconv.Atoi(os.Args[3])
	if err != nil { log.Panicf(err.Error()); }
	exchangeAddress := common.HexToAddress(os.Args[4])
	bloom, err := fillbloom.Initialize(int64(startBlock), int64(endBlock), exchangeAddress, conn)
	if err != nil { log.Panicf(err.Error()); }
	bloom.WriteTo(os.Stdout)
}
