package funds

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	orCommon "github.com/notegio/openrelay/common"
	tokenModule "github.com/notegio/openrelay/token"
	"github.com/notegio/openrelay/types"
	"math/big"
)

type BalanceChecker interface {
	GetBalance(tokenAddress, userAddress *types.Address) (*big.Int, error)
	GetAllowance(tokenAddress, ownerAddress, spenderAddress *types.Address) (*big.Int, error)
}

type rpcBalanceChecker struct {
	conn bind.ContractBackend
}

func (funds *rpcBalanceChecker) GetBalance(tokenAddrBytes, userAddrBytes *types.Address) (*big.Int, error) {
	token, err := tokenModule.NewToken(orCommon.ToGethAddress(tokenAddrBytes), funds.conn)
	if err != nil {
		return nil, err
	}
	return token.BalanceOf(nil, orCommon.ToGethAddress(userAddrBytes))
}

func (funds *rpcBalanceChecker) GetAllowance(tokenAddrBytes, ownerAddress, spenderAddress *types.Address) (*big.Int, error) {
	token, err := tokenModule.NewToken(orCommon.ToGethAddress(tokenAddrBytes), funds.conn)
	if err != nil {
		return nil, err
	}
	return token.Allowance(nil, orCommon.ToGethAddress(ownerAddress), orCommon.ToGethAddress(spenderAddress))
}

func NewRpcBalanceChecker(rpcUrl string) (BalanceChecker, error) {
	conn, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}
	if _, err = conn.SyncProgress(context.Background()); err != nil {
		// This is just here so that an RpcBalanceChecker can't be instantiated
		// successfully if the RPC server isn't responding properly. What RPC
		// function we call isn't important, but SyncProgress is pretty cheap.
		return nil, err
	}
	return &rpcBalanceChecker{conn}, nil
}
