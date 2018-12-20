package main

import (
	dbModule "github.com/notegio/openrelay/db"
	poolModule "github.com/notegio/openrelay/pool"
	"github.com/notegio/openrelay/common"
	"github.com/notegio/openrelay/types"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"log"
	"os"
	"strconv"
)


func main() {
	if len(os.Args) != 9 {
		log.Fatalf("Usage: poolmgr DB_CONNECTION_STRING DB_PASSWORD POOL_NAME SEARCH_STRING FEE_SHARE SENDER_ADDRESS FILTER_ADDRESS NETWORK_ID")
	}
	db, err := dbModule.GetDB(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err.Error())
	}
	poolHash := sha3.NewKeccak256()
	poolHash.Write([]byte(os.Args[3]))

	senderAddress, err := common.HexToAddress(os.Args[6])
	if err != nil {
		log.Fatalf("Bad senderAddress: %v", err.Error())
	}
	filterAddress, err := common.HexToAddress(os.Args[7])
	if err != nil {
		log.Fatalf("Bad filterAddress: %v", err.Error())
	}

	networkID, err := strconv.Atoi(os.Args[8])
	if err != nil {
		log.Fatalf("Bad network id: %v", err.Error())
	}


	pool := &poolModule.Pool{
		SearchTerms: os.Args[4],
		Expiration: 1744733652,
		Nonce: 0,
		FeeShare: os.Args[5],
		ID: poolHash.Sum(nil),
		SenderAddresses: types.NetworkAddressMap{uint(networkID): senderAddress},
		FilterAddresses: types.NetworkAddressMap{uint(networkID): filterAddress},
	}

	err = db.Debug().Model(&poolModule.Pool{}).Assign(pool).FirstOrCreate(pool).Error
	if err != nil {
		log.Fatalf("Error applying update: %v", err.Error())
	}
	log.Printf("Success!\n")

	// SearchTerms   string
	// Expiration    uint
	// Nonce         uint
	// FeeShare      string
	// ID            []byte
	// SenderAddress *types.Address
	// FilterAddress *types.Address
}
