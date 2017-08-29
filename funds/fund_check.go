package funds

import (
	tokenModule "github.com/notegio/0xrelay/token"
	"github.com/notegio/0xrelay/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"encoding/hex"
	"log"
)

type FundChecker interface {
	GetBalance(tokenAddress, userAddress [20]byte) (*big.Int, error)
	ValidateOrder(order *types.Order) (bool)
}

type rpcFundChecker struct {
	conn bind.ContractBackend
}

func bytesToAddress(data [20]byte) (common.Address) {
	return common.HexToAddress(hex.EncodeToString(data[:]))
}

func (funds *rpcFundChecker)GetBalance(tokenAddrBytes, userAddrBytes [20]byte) (*big.Int, error) {
	token, err := tokenModule.NewToken(bytesToAddress(tokenAddrBytes), funds.conn)
	if err != nil {
		return nil, err
	}
	return token.BalanceOf(nil, bytesToAddress(userAddrBytes))
}

func (funds *rpcFundChecker) checkFunds(tokenAddress, userAddress [20]byte, required [32]byte, respond chan bool) {
	requiredInt := new(big.Int)
	requiredInt.SetBytes(required[:])
	balance, err := funds.GetBalance(tokenAddress, userAddress)
	if err != nil {
		log.Printf("'%v': '%v'", err.Error(), hex.EncodeToString(tokenAddress[:]))
		respond <- false
		return
	}
	respond <- (requiredInt.Cmp(balance) >= 0)
}

// ValidateOrder makes sure that the maker of an order has sufficient funds to
// fill the order and pay makerFees
func (funds *rpcFundChecker)ValidateOrder(order *types.Order) bool {
	// TODO: Look this up from somewhere so it can work on different chains
	feeToken := common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")
	makerChan := make(chan bool)
	feeChan := make(chan bool)
	go funds.checkFunds(order.MakerToken, order.Maker, order.MakerTokenAmount, makerChan)
	go funds.checkFunds(feeToken, order.Maker, order.MakerFee, feeChan)
	return (<-makerChan && <-feeChan)
}

func NewRpcFundChecker(rpcUrl string) (FundChecker, error){
	conn, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}
	return &rpcFundChecker{conn}, nil
}
