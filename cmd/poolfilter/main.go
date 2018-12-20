package main

import (
	"context"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/cmd/cmdutils"
	"github.com/jinzhu/gorm"
	poolModule "github.com/notegio/openrelay/pool"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/redis.v3"
	"os"
	"os/signal"
	"log"
	"strconv"
	"fmt"
)

type PoolFilter struct {
	db *gorm.DB
	conn bind.ContractCaller
	networkID uint
	poolCache map[string]*poolModule.Pool
}

func (filter *PoolFilter) Filter(delivery channels.Delivery) bool {
	order, err := types.OrderFromBytes([]byte(delivery.Payload()))
	if err != nil {
		log.Printf("Invalid order format: %#x", delivery.Payload())
		return false
	}
	if !order.Signature.Verify(order.Maker, order.Hash()) {
		log.Printf("Invalid order signature");
		return false
	}
	pool, ok := filter.poolCache[fmt.Sprintf("%#x", order.PoolID)]
	if !ok {
		pool = &poolModule.Pool{}
		if err := filter.db.Model(&poolModule.Pool{}).Where("ID = ?", order.PoolID).First(pool).Error; err != nil {
			log.Fatalf("Error getting pool: %#x - Error: %v", order.PoolID, err.Error())
		}
		pool.SetConn(filter.conn)
		filter.poolCache[fmt.Sprintf("%#x", order.PoolID)] = pool
	}
	valid, err := pool.CheckFilter(order, filter.networkID)
	if err != nil {
		delivery.Return()
		log.Fatalf("Error filtering order: %v", err.Error())
	}
	log.Printf("Order %#x is %v valid for target pool %#x", order.Hash(), valid, pool.ID)
	return valid
}

func main() {
	db, err := dbModule.GetDB(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err.Error())
	}
	redisURL := os.Args[3]
	rpcURL := os.Args[4]
	// src := os.Args[3]
	if redisURL == "" {
		log.Fatalf("Please specify redis URL")
	}
	if rpcURL == "" {
		log.Fatalf("Please specify RPC URL")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	conn, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Error connecting to RPC: %v", err.Error())
	}
	var networkID uint64
	if _, err = conn.NetworkID(context.Background()); err != nil {
		log.Fatalf("Error connecting to RPC: %v", err.Error())
	}
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		concurrency = 5
	}
	var channelStrings []string
	for _, arg := range os.Args[5:] {
		channelStrings = append(channelStrings, arg)
	}

	var poolFilter channels.RelayFilter
	poolFilter = &PoolFilter{db, conn, uint(networkID), make(map[string]*poolModule.Pool)}
	var relays []channels.Relay
	for _, channelString := range channelStrings {
		consumerChannel, publisher, _, err := cmdutils.ParseChannels(channelString, redisClient)
		if err != nil { log.Fatalf(err.Error()) }
		relay := channels.NewRelay(consumerChannel, publisher, poolFilter, concurrency)
		relay.Start()
		relays = append(relays, relay)
	}

	log.Printf("Starting poolcheck")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		break
	}
	for _, relay := range relays {
		relay.Stop()
	}
}
