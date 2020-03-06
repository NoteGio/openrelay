package exchangecontract

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/notegio/openrelay/common"
	"os"
	"strconv"
)

func (exc *Exchangecontract) ZRX_ASSET_DATA(opts *bind.CallOpts) ([]byte, error) {
	chainid, err := strconv.Atoi(os.Getenv("CHAIN_ID"))
	if err != nil { return nil, err }
	return common.DefaultFeeAssetData(uint(chainid))
}
