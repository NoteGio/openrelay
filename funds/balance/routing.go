package balance

import (
	"context"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/channels"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"sync"
	"fmt"
	"log"
)

type BalanceChecker interface {
	GetBalance(asset types.AssetData, userAddress *types.Address) (*big.Int, error)
	GetAllowance(asset types.AssetData, ownerAddress, spenderAddress *types.Address) (*big.Int, error)
}

type CachedBalanceChecker interface {
	BalanceChecker
	Consume(msg channels.Delivery)
}

type routingBalanceChecker struct {
	lookupCache map[string]*big.Int
	cacheMutex *sync.Mutex
	assetBalanceCheckers map[string]BalanceChecker
}

func (funds *routingBalanceChecker) GetBalance(tokenAsset types.AssetData, userAddrBytes *types.Address) (*big.Int, error) {
	cacheKey := fmt.Sprintf("b-%#x-%#x", tokenAsset[:], userAddrBytes[:])
	if funds.lookupCache != nil {
		funds.cacheMutex.Lock()
		if balance, ok := funds.lookupCache[cacheKey]; ok {
			funds.cacheMutex.Unlock()

			return balance, nil
		}
		funds.cacheMutex.Unlock()
	}
	balanceChecker, ok := funds.assetBalanceCheckers[fmt.Sprintf("%#x", tokenAsset.ProxyId())]
	if !ok {
		return nil, fmt.Errorf("Could not find balance checker for asset type '%#x'", tokenAsset.ProxyId())
	}
	balance, err := balanceChecker.GetBalance(tokenAsset, userAddrBytes)
	if err != nil {
		return nil, err
	}
	if funds.lookupCache != nil {
		funds.cacheMutex.Lock()
		defer funds.cacheMutex.Unlock()
		funds.lookupCache[cacheKey] = balance
	}
	return balance, nil
}
func (funds *routingBalanceChecker) GetAllowance(tokenAsset types.AssetData, ownerAddress, spenderAddress *types.Address) (*big.Int, error) {
	cacheKey := fmt.Sprintf("a-%#x-%#x-%#x", tokenAsset[:], ownerAddress[:], spenderAddress[:])
	if funds.lookupCache != nil {
		funds.cacheMutex.Lock()
		if balance, ok := funds.lookupCache[cacheKey]; ok {
			funds.cacheMutex.Unlock()

			return balance, nil
		}
		funds.cacheMutex.Unlock()
	}
	balanceChecker, ok := funds.assetBalanceCheckers[fmt.Sprintf("%#x", tokenAsset.ProxyId())]
	if !ok {
		return nil, fmt.Errorf("Could not find balance checker for asset type '%#x'", tokenAsset.ProxyId())
	}
	allowance, err := balanceChecker.GetAllowance(tokenAsset, ownerAddress, spenderAddress)
	if err != nil {
		return nil, err
	}
	if funds.lookupCache != nil {
		funds.cacheMutex.Lock()
		defer funds.cacheMutex.Unlock()
		funds.lookupCache[cacheKey] = allowance
	}
	return allowance, nil
}

func (funds *routingBalanceChecker) Consume(msg channels.Delivery) {
	if funds.lookupCache == nil {
		log.Printf("Inititalizing lookup cache")
	}
	defer msg.Ack()
	funds.cacheMutex.Lock()
	defer funds.cacheMutex.Unlock()
	funds.lookupCache = make(map[string]*big.Int)
}

func NewRpcRoutingBalanceChecker(rpcURL string) (CachedBalanceChecker, error) {
	conn, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	if _, err = conn.SyncProgress(context.Background()); err != nil {
		// This is just here so that an RpcBalanceChecker can't be instantiated
		// successfully if the RPC server isn't responding properly. What RPC
		// function we call isn't important, but SyncProgress is pretty cheap.
		return nil, err
	}
	checkers := make(map[string]BalanceChecker)
	checkers["0xf47261b0"] = NewRpcERC20BalanceChecker(conn)
	checkers["0x02571792"] = NewRpcERC721BalanceChecker(conn)
	return &routingBalanceChecker{nil, &sync.Mutex{}, checkers}, nil
}
