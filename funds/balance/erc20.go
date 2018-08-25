package balance

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	orCommon "github.com/notegio/openrelay/common"
	tokenModule "github.com/notegio/openrelay/token"
	"github.com/notegio/openrelay/types"
	"math/big"
)

type rpcERC20BalanceChecker struct {
	conn bind.ContractBackend
}

func (funds *rpcERC20BalanceChecker) GetBalance(tokenAsset types.AssetData, userAddrBytes *types.Address) (*big.Int, error) {
	token, err := tokenModule.NewToken(orCommon.ToGethAddress(tokenAsset.Address()), funds.conn)
	if err != nil {
		return nil, err
	}
	return token.BalanceOf(nil, orCommon.ToGethAddress(userAddrBytes))
}

func (funds *rpcERC20BalanceChecker) GetAllowance(tokenAsset types.AssetData, ownerAddress, spenderAddress *types.Address) (*big.Int, error) {
	token, err := tokenModule.NewToken(orCommon.ToGethAddress(tokenAsset.Address()), funds.conn)
	if err != nil {
		return nil, err
	}
	return token.Allowance(nil, orCommon.ToGethAddress(ownerAddress), orCommon.ToGethAddress(spenderAddress))
}


func NewRpcERC20BalanceChecker(conn bind.ContractBackend) (BalanceChecker) {
	return &rpcERC20BalanceChecker{conn}
}
