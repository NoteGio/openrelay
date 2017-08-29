package funds

import (
	tokenModule "github.com/notegio/0xrelay/token"
	orCommon "github.com/notegio/0xrelay/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type BalanceChecker interface {
	GetBalance(tokenAddress, userAddress [20]byte) (*big.Int, error)
}

type rpcBalanceChecker struct {
	conn bind.ContractBackend
}



func (funds *rpcBalanceChecker)GetBalance(tokenAddrBytes, userAddrBytes [20]byte) (*big.Int, error) {
	token, err := tokenModule.NewToken(orCommon.BytesToAddress(tokenAddrBytes), funds.conn)
	if err != nil {
		return nil, err
	}
	return token.BalanceOf(nil, orCommon.BytesToAddress(userAddrBytes))
}

func NewRpcBalanceChecker(rpcUrl string) (BalanceChecker, error){
	conn, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}
	return &rpcBalanceChecker{conn}, nil
}
