package balance

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	orCommon "github.com/notegio/openrelay/common"
	tokenModule "github.com/notegio/openrelay/token"
	"github.com/notegio/openrelay/types"
	"math/big"
	"bytes"
	"log"
)

type rpcERC721BalanceChecker struct {
	conn bind.ContractBackend
}

func (funds *rpcERC721BalanceChecker) GetBalance(tokenAsset types.AssetData, userAddrBytes *types.Address) (*big.Int, error) {
	token, err := tokenModule.NewERC721Token(orCommon.ToGethAddress(tokenAsset.Address()), funds.conn)
	if err != nil {
		return nil, err
	}
	if owner, err := token.OwnerOf(nil, tokenAsset.TokenID().Big()); err != nil {
		return nil, err
	} else if !bytes.Equal(userAddrBytes[:], owner[:]) {
		return big.NewInt(0), nil
	}
	return big.NewInt(1), nil
}

func (funds *rpcERC721BalanceChecker) GetAllowance(tokenAsset types.AssetData, ownerAddress, spenderAddress *types.Address) (*big.Int, error) {
	token, err := tokenModule.NewERC721Token(orCommon.ToGethAddress(tokenAsset.Address()), funds.conn)
	if err != nil {
		return nil, err
	}
	if approved, err := token.IsApprovedForAll(nil, orCommon.ToGethAddress(ownerAddress), orCommon.ToGethAddress(spenderAddress)); err != nil {
		if err.Error() != "VM Exception while processing transaction: revert" {
			return nil, err
		} else {
			// Some early ERC721 tokens don't implement this
			log.Printf("Token %#x does not provide isApprovedForAll. Testing getApproved.", tokenAsset.Address())
		}
	} else if approved {
		return big.NewInt(1), nil
	}
	if operator, err := token.GetApproved(nil, tokenAsset.TokenID().Big()); err != nil {
		return nil, err
	} else if bytes.Equal(spenderAddress[:], operator[:]) {
		return big.NewInt(1), nil
	}
	return big.NewInt(0), nil
}


func NewRpcERC721BalanceChecker(conn bind.ContractBackend) (BalanceChecker) {
	return &rpcERC721BalanceChecker{conn}
}
