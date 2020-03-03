package exchangecontract

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"os"
)

func (exc *Exchangecontract) ZRX_ASSET_DATA(opts *bind.CallOpts) ([]byte, error) {
	return hex.DecodeString(os.Getenv("DEFAULT_FEE_ASSETDATA"))
}
