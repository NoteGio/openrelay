package common

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/notegio/openrelay/aws"
	"github.com/notegio/openrelay/types"
	"io/ioutil"
	"os"
	"strings"
)

func BytesToAddress(data [20]byte) common.Address {
	return common.HexToAddress(hex.EncodeToString(data[:]))
}

func ToGethAddress(data *types.Address) common.Address {
	return common.HexToAddress(hex.EncodeToString(data[:]))
}

func BytesToOrAddress(data [20]byte) *types.Address {
	addr := &types.Address{}
	copy(addr[:], data[:])
	return addr
}

func HexToBytes(hexString string) ([20]byte, error) {
	slice, err := hex.DecodeString(strings.TrimPrefix(hexString, "0x"))
	result := [20]byte{}
	if err != nil {
		return result, err
	}
	copy(result[:], slice[:])
	return result, nil
}

// GetSecret retrieves a secret from various supported secret stores
func GetSecret(uri string) string {
	if strings.HasPrefix(uri, "file://") {
		secretBytes, err := ioutil.ReadFile(strings.TrimPrefix(uri, "file://"))
		if err == nil {
			return string(secretBytes)
		}
	}
	if strings.HasPrefix(uri, "env://") {
		if secret := os.Getenv(strings.TrimPrefix(uri, "env://")); secret != "" {
			return secret
		}
	}
	if strings.HasPrefix(uri, "ssm-param://") {
		if secret := aws.GetParameter(strings.TrimPrefix(uri, "ssm-param://")); secret != "" {
			return secret
		}
	}
	return uri
}
