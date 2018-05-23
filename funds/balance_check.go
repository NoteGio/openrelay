package funds

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	orCommon "github.com/notegio/openrelay/common"
	tokenModule "github.com/notegio/openrelay/token"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/channels"
	"math/big"
	"sync"
	"fmt"
	"log"
)

type BalanceChecker interface {
	GetBalance(tokenAddress, userAddress *types.Address) (*big.Int, error)
	GetAllowance(tokenAddress, ownerAddress, spenderAddress *types.Address) (*big.Int, error)
	Consume(msg channels.Delivery)
}

type rpcBalanceChecker struct {
	conn bind.ContractBackend
	lookupCache map[string]*big.Int
	cacheMutex *sync.Mutex
}

func (funds *rpcBalanceChecker) GetBalance(tokenAddrBytes, userAddrBytes *types.Address) (*big.Int, error) {
	cacheKey := fmt.Sprintf("b-%#x-%#x", tokenAddrBytes[:], userAddrBytes[:])
	if funds.lookupCache != nil {
		funds.cacheMutex.Lock()
		if balance, ok := funds.lookupCache[cacheKey]; ok {
			funds.cacheMutex.Unlock()
			return balance, nil
		}
		funds.cacheMutex.Unlock()
	}
	token, err := tokenModule.NewToken(orCommon.ToGethAddress(tokenAddrBytes), funds.conn)
	if err != nil {
		return nil, err
	}
	balance, err := token.BalanceOf(nil, orCommon.ToGethAddress(userAddrBytes))
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

func (funds *rpcBalanceChecker) GetAllowance(tokenAddrBytes, ownerAddress, spenderAddress *types.Address) (*big.Int, error) {
	cacheKey := fmt.Sprintf("a-%#x-%#x-%#x", tokenAddrBytes[:], ownerAddress[:], spenderAddress[:])
	if funds.lookupCache != nil {
		funds.cacheMutex.Lock()
		if allowance, ok := funds.lookupCache[cacheKey]; ok {
			funds.cacheMutex.Unlock()
			return allowance, nil
		}
		funds.cacheMutex.Unlock()
	}
	token, err := tokenModule.NewToken(orCommon.ToGethAddress(tokenAddrBytes), funds.conn)
	if err != nil {
		return nil, err
	}
	allowance, err := token.Allowance(nil, orCommon.ToGethAddress(ownerAddress), orCommon.ToGethAddress(spenderAddress))
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

func (funds *rpcBalanceChecker) Consume(msg channels.Delivery) {
	if funds.lookupCache == nil {
		log.Printf("Inititalizing lookup cache")
	}
	defer msg.Ack()
	funds.cacheMutex.Lock()
	defer funds.cacheMutex.Unlock()
	funds.lookupCache = make(map[string]*big.Int)
}

func NewRpcBalanceChecker(rpcURL string) (BalanceChecker, error) {
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
	return &rpcBalanceChecker{conn, nil, &sync.Mutex{}}, nil
}
