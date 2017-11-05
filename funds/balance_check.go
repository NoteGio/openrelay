package funds

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	orCommon "github.com/notegio/openrelay/common"
	tokenModule "github.com/notegio/openrelay/token"
	"math/big"
)

type BalanceChecker interface {
	GetBalance(tokenAddress, userAddress [20]byte) (*big.Int, error)
	GetAllowance(tokenAddress, ownerAddress, spenderAddress [20]byte) (*big.Int, error)
}

type rpcBalanceChecker struct {
	conn bind.ContractBackend
}

func (funds *rpcBalanceChecker) GetBalance(tokenAddrBytes, userAddrBytes [20]byte) (*big.Int, error) {
	token, err := tokenModule.NewToken(orCommon.BytesToAddress(tokenAddrBytes), funds.conn)
	if err != nil {
			return nil, err
	}
	return token.BalanceOf(nil, orCommon.BytesToAddress(userAddrBytes))
}

func (funds *rpcBalanceChecker) GetAllowance(tokenAddrBytes, ownerAddress, spenderAddress [20]byte) (*big.Int, error) {
	token, err := tokenModule.NewToken(orCommon.BytesToAddress(tokenAddrBytes), funds.conn)
	if err != nil {
		return nil, err
	}
	return token.Allowance(nil, orCommon.BytesToAddress(ownerAddress), orCommon.BytesToAddress(spenderAddress))
}

func NewRpcBalanceChecker(rpcUrl string) (BalanceChecker, error) {
	conn, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}
	if _, err = conn.HeaderByNumber(context.Background(), nil); err != nil {
		return nil, err
	}
	return &rpcBalanceChecker{conn}, nil
}
