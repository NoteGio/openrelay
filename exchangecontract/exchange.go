// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package exchangecontract

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// LibFillResultsBatchMatchedFillResults is an auto generated low-level Go binding around an user-defined struct.
type LibFillResultsBatchMatchedFillResults struct {
	Left                    []LibFillResultsFillResults
	Right                   []LibFillResultsFillResults
	ProfitInLeftMakerAsset  *big.Int
	ProfitInRightMakerAsset *big.Int
}

// LibFillResultsFillResults is an auto generated low-level Go binding around an user-defined struct.
type LibFillResultsFillResults struct {
	MakerAssetFilledAmount *big.Int
	TakerAssetFilledAmount *big.Int
	MakerFeePaid           *big.Int
	TakerFeePaid           *big.Int
	ProtocolFeePaid        *big.Int
}

// LibFillResultsMatchedFillResults is an auto generated low-level Go binding around an user-defined struct.
type LibFillResultsMatchedFillResults struct {
	Left                    LibFillResultsFillResults
	Right                   LibFillResultsFillResults
	ProfitInLeftMakerAsset  *big.Int
	ProfitInRightMakerAsset *big.Int
}

// LibOrderOrder is an auto generated low-level Go binding around an user-defined struct.
type LibOrderOrder struct {
	MakerAddress          common.Address
	TakerAddress          common.Address
	FeeRecipientAddress   common.Address
	SenderAddress         common.Address
	MakerAssetAmount      *big.Int
	TakerAssetAmount      *big.Int
	MakerFee              *big.Int
	TakerFee              *big.Int
	ExpirationTimeSeconds *big.Int
	Salt                  *big.Int
	MakerAssetData        []byte
	TakerAssetData        []byte
	MakerFeeAssetData     []byte
	TakerFeeAssetData     []byte
}

// LibOrderOrderInfo is an auto generated low-level Go binding around an user-defined struct.
type LibOrderOrderInfo struct {
	OrderStatus                 uint8
	OrderHash                   [32]byte
	OrderTakerAssetFilledAmount *big.Int
}

// LibZeroExTransactionZeroExTransaction is an auto generated low-level Go binding around an user-defined struct.
type LibZeroExTransactionZeroExTransaction struct {
	Salt                  *big.Int
	ExpirationTimeSeconds *big.Int
	GasPrice              *big.Int
	SignerAddress         common.Address
	Data                  []byte
}

// ExchangecontractABI is the input ABI used to generate the binding from.
const ExchangecontractABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"id\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetProxy\",\"type\":\"address\"}],\"name\":\"AssetProxyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"Cancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"orderSenderAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"orderEpoch\",\"type\":\"uint256\"}],\"name\":\"CancelUpTo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"name\":\"Fill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldProtocolFeeCollector\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"updatedProtocolFeeCollector\",\"type\":\"address\"}],\"name\":\"ProtocolFeeCollectorAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldProtocolFeeMultiplier\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedProtocolFeeMultiplier\",\"type\":\"uint256\"}],\"name\":\"ProtocolFeeMultiplier\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isApproved\",\"type\":\"bool\"}],\"name\":\"SignatureValidatorApproval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"transactionHash\",\"type\":\"bytes32\"}],\"name\":\"TransactionExecution\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"EIP1271_MAGIC_VALUE\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"EIP712_EXCHANGE_DOMAIN_HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedValidators\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"}],\"name\":\"batchCancelOrders\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLibZeroExTransaction.ZeroExTransaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchExecuteTransactions\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"returnData\",\"type\":\"bytes[]\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrKillOrders\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"fillResults\",\"type\":\"tuple[]\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrders\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"fillResults\",\"type\":\"tuple[]\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrdersNoThrow\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"fillResults\",\"type\":\"tuple[]\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"leftOrders\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"rightOrders\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"leftSignatures\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes[]\",\"name\":\"rightSignatures\",\"type\":\"bytes[]\"}],\"name\":\"batchMatchOrders\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"left\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"right\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"profitInLeftMakerAsset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"profitInRightMakerAsset\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.BatchMatchedFillResults\",\"name\":\"batchMatchedFillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"leftOrders\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"rightOrders\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"leftSignatures\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes[]\",\"name\":\"rightSignatures\",\"type\":\"bytes[]\"}],\"name\":\"batchMatchOrdersWithMaximalFill\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"left\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"right\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"profitInLeftMakerAsset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"profitInRightMakerAsset\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.BatchMatchedFillResults\",\"name\":\"batchMatchedFillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"cancelOrder\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"targetOrderEpoch\",\"type\":\"uint256\"}],\"name\":\"cancelOrdersUpTo\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"cancelled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentContextAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"detachProtocolFeeCollector\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLibZeroExTransaction.ZeroExTransaction\",\"name\":\"transaction\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"executeTransaction\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"fillOrKillOrder\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"fillOrder\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"filled\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"}],\"name\":\"getAssetProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"assetProxy\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"getOrderInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"enumLibOrder.OrderStatus\",\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"internalType\":\"structLibOrder.OrderInfo\",\"name\":\"orderInfo\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidHashSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidOrderSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLibZeroExTransaction.ZeroExTransaction\",\"name\":\"transaction\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidTransactionSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketBuyOrdersFillOrKill\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketBuyOrdersNoThrow\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketSellOrdersFillOrKill\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketSellOrdersNoThrow\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"leftOrder\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"rightOrder\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"leftSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rightSignature\",\"type\":\"bytes\"}],\"name\":\"matchOrders\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"left\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"right\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"profitInLeftMakerAsset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"profitInRightMakerAsset\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.MatchedFillResults\",\"name\":\"matchedFillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"leftOrder\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"rightOrder\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"leftSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rightSignature\",\"type\":\"bytes\"}],\"name\":\"matchOrdersWithMaximalFill\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"left\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"right\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"profitInLeftMakerAsset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"profitInRightMakerAsset\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.MatchedFillResults\",\"name\":\"matchedFillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"orderEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"preSign\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"preSigned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"protocolFeeCollector\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"protocolFeeMultiplier\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetProxy\",\"type\":\"address\"}],\"name\":\"registerAssetProxy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"updatedProtocolFeeCollector\",\"type\":\"address\"}],\"name\":\"setProtocolFeeCollectorAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"updatedProtocolFeeMultiplier\",\"type\":\"uint256\"}],\"name\":\"setProtocolFeeMultiplier\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approval\",\"type\":\"bool\"}],\"name\":\"setSignatureValidatorApproval\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"assetData\",\"type\":\"bytes[]\"},{\"internalType\":\"address[]\",\"name\":\"fromAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"toAddresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"simulateDispatchTransferFromCalls\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transactionsExecuted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Exchangecontract is an auto generated Go binding around an Ethereum contract.
type Exchangecontract struct {
	ExchangecontractCaller     // Read-only binding to the contract
	ExchangecontractTransactor // Write-only binding to the contract
	ExchangecontractFilterer   // Log filterer for contract events
}

// ExchangecontractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExchangecontractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangecontractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExchangecontractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangecontractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExchangecontractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangecontractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExchangecontractSession struct {
	Contract     *Exchangecontract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ExchangecontractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExchangecontractCallerSession struct {
	Contract *ExchangecontractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ExchangecontractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExchangecontractTransactorSession struct {
	Contract     *ExchangecontractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ExchangecontractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExchangecontractRaw struct {
	Contract *Exchangecontract // Generic contract binding to access the raw methods on
}

// ExchangecontractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExchangecontractCallerRaw struct {
	Contract *ExchangecontractCaller // Generic read-only contract binding to access the raw methods on
}

// ExchangecontractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExchangecontractTransactorRaw struct {
	Contract *ExchangecontractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExchangecontract creates a new instance of Exchangecontract, bound to a specific deployed contract.
func NewExchangecontract(address common.Address, backend bind.ContractBackend) (*Exchangecontract, error) {
	contract, err := bindExchangecontract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Exchangecontract{ExchangecontractCaller: ExchangecontractCaller{contract: contract}, ExchangecontractTransactor: ExchangecontractTransactor{contract: contract}, ExchangecontractFilterer: ExchangecontractFilterer{contract: contract}}, nil
}

// NewExchangecontractCaller creates a new read-only instance of Exchangecontract, bound to a specific deployed contract.
func NewExchangecontractCaller(address common.Address, caller bind.ContractCaller) (*ExchangecontractCaller, error) {
	contract, err := bindExchangecontract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangecontractCaller{contract: contract}, nil
}

// NewExchangecontractTransactor creates a new write-only instance of Exchangecontract, bound to a specific deployed contract.
func NewExchangecontractTransactor(address common.Address, transactor bind.ContractTransactor) (*ExchangecontractTransactor, error) {
	contract, err := bindExchangecontract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangecontractTransactor{contract: contract}, nil
}

// NewExchangecontractFilterer creates a new log filterer instance of Exchangecontract, bound to a specific deployed contract.
func NewExchangecontractFilterer(address common.Address, filterer bind.ContractFilterer) (*ExchangecontractFilterer, error) {
	contract, err := bindExchangecontract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExchangecontractFilterer{contract: contract}, nil
}

// bindExchangecontract binds a generic wrapper to an already deployed contract.
func bindExchangecontract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangecontractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Exchangecontract *ExchangecontractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Exchangecontract.Contract.ExchangecontractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Exchangecontract *ExchangecontractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchangecontract.Contract.ExchangecontractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Exchangecontract *ExchangecontractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Exchangecontract.Contract.ExchangecontractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Exchangecontract *ExchangecontractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Exchangecontract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Exchangecontract *ExchangecontractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchangecontract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Exchangecontract *ExchangecontractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Exchangecontract.Contract.contract.Transact(opts, method, params...)
}

// EIP1271MAGICVALUE is a free data retrieval call binding the contract method 0xdd885e2d.
//
// Solidity: function EIP1271_MAGIC_VALUE() constant returns(bytes4)
func (_Exchangecontract *ExchangecontractCaller) EIP1271MAGICVALUE(opts *bind.CallOpts) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "EIP1271_MAGIC_VALUE")
	return *ret0, err
}

// EIP1271MAGICVALUE is a free data retrieval call binding the contract method 0xdd885e2d.
//
// Solidity: function EIP1271_MAGIC_VALUE() constant returns(bytes4)
func (_Exchangecontract *ExchangecontractSession) EIP1271MAGICVALUE() ([4]byte, error) {
	return _Exchangecontract.Contract.EIP1271MAGICVALUE(&_Exchangecontract.CallOpts)
}

// EIP1271MAGICVALUE is a free data retrieval call binding the contract method 0xdd885e2d.
//
// Solidity: function EIP1271_MAGIC_VALUE() constant returns(bytes4)
func (_Exchangecontract *ExchangecontractCallerSession) EIP1271MAGICVALUE() ([4]byte, error) {
	return _Exchangecontract.Contract.EIP1271MAGICVALUE(&_Exchangecontract.CallOpts)
}

// EIP712EXCHANGEDOMAINHASH is a free data retrieval call binding the contract method 0xc26cfecd.
//
// Solidity: function EIP712_EXCHANGE_DOMAIN_HASH() constant returns(bytes32)
func (_Exchangecontract *ExchangecontractCaller) EIP712EXCHANGEDOMAINHASH(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "EIP712_EXCHANGE_DOMAIN_HASH")
	return *ret0, err
}

// EIP712EXCHANGEDOMAINHASH is a free data retrieval call binding the contract method 0xc26cfecd.
//
// Solidity: function EIP712_EXCHANGE_DOMAIN_HASH() constant returns(bytes32)
func (_Exchangecontract *ExchangecontractSession) EIP712EXCHANGEDOMAINHASH() ([32]byte, error) {
	return _Exchangecontract.Contract.EIP712EXCHANGEDOMAINHASH(&_Exchangecontract.CallOpts)
}

// EIP712EXCHANGEDOMAINHASH is a free data retrieval call binding the contract method 0xc26cfecd.
//
// Solidity: function EIP712_EXCHANGE_DOMAIN_HASH() constant returns(bytes32)
func (_Exchangecontract *ExchangecontractCallerSession) EIP712EXCHANGEDOMAINHASH() ([32]byte, error) {
	return _Exchangecontract.Contract.EIP712EXCHANGEDOMAINHASH(&_Exchangecontract.CallOpts)
}

// AllowedValidators is a free data retrieval call binding the contract method 0x7b8e3514.
//
// Solidity: function allowedValidators(address , address ) constant returns(bool)
func (_Exchangecontract *ExchangecontractCaller) AllowedValidators(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "allowedValidators", arg0, arg1)
	return *ret0, err
}

// AllowedValidators is a free data retrieval call binding the contract method 0x7b8e3514.
//
// Solidity: function allowedValidators(address , address ) constant returns(bool)
func (_Exchangecontract *ExchangecontractSession) AllowedValidators(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _Exchangecontract.Contract.AllowedValidators(&_Exchangecontract.CallOpts, arg0, arg1)
}

// AllowedValidators is a free data retrieval call binding the contract method 0x7b8e3514.
//
// Solidity: function allowedValidators(address , address ) constant returns(bool)
func (_Exchangecontract *ExchangecontractCallerSession) AllowedValidators(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _Exchangecontract.Contract.AllowedValidators(&_Exchangecontract.CallOpts, arg0, arg1)
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled(bytes32 ) constant returns(bool)
func (_Exchangecontract *ExchangecontractCaller) Cancelled(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "cancelled", arg0)
	return *ret0, err
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled(bytes32 ) constant returns(bool)
func (_Exchangecontract *ExchangecontractSession) Cancelled(arg0 [32]byte) (bool, error) {
	return _Exchangecontract.Contract.Cancelled(&_Exchangecontract.CallOpts, arg0)
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled(bytes32 ) constant returns(bool)
func (_Exchangecontract *ExchangecontractCallerSession) Cancelled(arg0 [32]byte) (bool, error) {
	return _Exchangecontract.Contract.Cancelled(&_Exchangecontract.CallOpts, arg0)
}

// CurrentContextAddress is a free data retrieval call binding the contract method 0xeea086ba.
//
// Solidity: function currentContextAddress() constant returns(address)
func (_Exchangecontract *ExchangecontractCaller) CurrentContextAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "currentContextAddress")
	return *ret0, err
}

// CurrentContextAddress is a free data retrieval call binding the contract method 0xeea086ba.
//
// Solidity: function currentContextAddress() constant returns(address)
func (_Exchangecontract *ExchangecontractSession) CurrentContextAddress() (common.Address, error) {
	return _Exchangecontract.Contract.CurrentContextAddress(&_Exchangecontract.CallOpts)
}

// CurrentContextAddress is a free data retrieval call binding the contract method 0xeea086ba.
//
// Solidity: function currentContextAddress() constant returns(address)
func (_Exchangecontract *ExchangecontractCallerSession) CurrentContextAddress() (common.Address, error) {
	return _Exchangecontract.Contract.CurrentContextAddress(&_Exchangecontract.CallOpts)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled(bytes32 ) constant returns(uint256)
func (_Exchangecontract *ExchangecontractCaller) Filled(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "filled", arg0)
	return *ret0, err
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled(bytes32 ) constant returns(uint256)
func (_Exchangecontract *ExchangecontractSession) Filled(arg0 [32]byte) (*big.Int, error) {
	return _Exchangecontract.Contract.Filled(&_Exchangecontract.CallOpts, arg0)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled(bytes32 ) constant returns(uint256)
func (_Exchangecontract *ExchangecontractCallerSession) Filled(arg0 [32]byte) (*big.Int, error) {
	return _Exchangecontract.Contract.Filled(&_Exchangecontract.CallOpts, arg0)
}

// GetAssetProxy is a free data retrieval call binding the contract method 0x60704108.
//
// Solidity: function getAssetProxy(bytes4 assetProxyId) constant returns(address assetProxy)
func (_Exchangecontract *ExchangecontractCaller) GetAssetProxy(opts *bind.CallOpts, assetProxyId [4]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "getAssetProxy", assetProxyId)
	return *ret0, err
}

// GetAssetProxy is a free data retrieval call binding the contract method 0x60704108.
//
// Solidity: function getAssetProxy(bytes4 assetProxyId) constant returns(address assetProxy)
func (_Exchangecontract *ExchangecontractSession) GetAssetProxy(assetProxyId [4]byte) (common.Address, error) {
	return _Exchangecontract.Contract.GetAssetProxy(&_Exchangecontract.CallOpts, assetProxyId)
}

// GetAssetProxy is a free data retrieval call binding the contract method 0x60704108.
//
// Solidity: function getAssetProxy(bytes4 assetProxyId) constant returns(address assetProxy)
func (_Exchangecontract *ExchangecontractCallerSession) GetAssetProxy(assetProxyId [4]byte) (common.Address, error) {
	return _Exchangecontract.Contract.GetAssetProxy(&_Exchangecontract.CallOpts, assetProxyId)
}

// GetOrderInfo is a free data retrieval call binding the contract method 0x9d3fa4b9.
//
// Solidity: function getOrderInfo(LibOrderOrder order) constant returns(LibOrderOrderInfo orderInfo)
func (_Exchangecontract *ExchangecontractCaller) GetOrderInfo(opts *bind.CallOpts, order LibOrderOrder) (LibOrderOrderInfo, error) {
	var (
		ret0 = new(LibOrderOrderInfo)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "getOrderInfo", order)
	return *ret0, err
}

// GetOrderInfo is a free data retrieval call binding the contract method 0x9d3fa4b9.
//
// Solidity: function getOrderInfo(LibOrderOrder order) constant returns(LibOrderOrderInfo orderInfo)
func (_Exchangecontract *ExchangecontractSession) GetOrderInfo(order LibOrderOrder) (LibOrderOrderInfo, error) {
	return _Exchangecontract.Contract.GetOrderInfo(&_Exchangecontract.CallOpts, order)
}

// GetOrderInfo is a free data retrieval call binding the contract method 0x9d3fa4b9.
//
// Solidity: function getOrderInfo(LibOrderOrder order) constant returns(LibOrderOrderInfo orderInfo)
func (_Exchangecontract *ExchangecontractCallerSession) GetOrderInfo(order LibOrderOrder) (LibOrderOrderInfo, error) {
	return _Exchangecontract.Contract.GetOrderInfo(&_Exchangecontract.CallOpts, order)
}

// IsValidHashSignature is a free data retrieval call binding the contract method 0x8171c407.
//
// Solidity: function isValidHashSignature(bytes32 hash, address signerAddress, bytes signature) constant returns(bool isValid)
func (_Exchangecontract *ExchangecontractCaller) IsValidHashSignature(opts *bind.CallOpts, hash [32]byte, signerAddress common.Address, signature []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "isValidHashSignature", hash, signerAddress, signature)
	return *ret0, err
}

// IsValidHashSignature is a free data retrieval call binding the contract method 0x8171c407.
//
// Solidity: function isValidHashSignature(bytes32 hash, address signerAddress, bytes signature) constant returns(bool isValid)
func (_Exchangecontract *ExchangecontractSession) IsValidHashSignature(hash [32]byte, signerAddress common.Address, signature []byte) (bool, error) {
	return _Exchangecontract.Contract.IsValidHashSignature(&_Exchangecontract.CallOpts, hash, signerAddress, signature)
}

// IsValidHashSignature is a free data retrieval call binding the contract method 0x8171c407.
//
// Solidity: function isValidHashSignature(bytes32 hash, address signerAddress, bytes signature) constant returns(bool isValid)
func (_Exchangecontract *ExchangecontractCallerSession) IsValidHashSignature(hash [32]byte, signerAddress common.Address, signature []byte) (bool, error) {
	return _Exchangecontract.Contract.IsValidHashSignature(&_Exchangecontract.CallOpts, hash, signerAddress, signature)
}

// IsValidOrderSignature is a free data retrieval call binding the contract method 0xa12dcc6f.
//
// Solidity: function isValidOrderSignature(LibOrderOrder order, bytes signature) constant returns(bool isValid)
func (_Exchangecontract *ExchangecontractCaller) IsValidOrderSignature(opts *bind.CallOpts, order LibOrderOrder, signature []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "isValidOrderSignature", order, signature)
	return *ret0, err
}

// IsValidOrderSignature is a free data retrieval call binding the contract method 0xa12dcc6f.
//
// Solidity: function isValidOrderSignature(LibOrderOrder order, bytes signature) constant returns(bool isValid)
func (_Exchangecontract *ExchangecontractSession) IsValidOrderSignature(order LibOrderOrder, signature []byte) (bool, error) {
	return _Exchangecontract.Contract.IsValidOrderSignature(&_Exchangecontract.CallOpts, order, signature)
}

// IsValidOrderSignature is a free data retrieval call binding the contract method 0xa12dcc6f.
//
// Solidity: function isValidOrderSignature(LibOrderOrder order, bytes signature) constant returns(bool isValid)
func (_Exchangecontract *ExchangecontractCallerSession) IsValidOrderSignature(order LibOrderOrder, signature []byte) (bool, error) {
	return _Exchangecontract.Contract.IsValidOrderSignature(&_Exchangecontract.CallOpts, order, signature)
}

// IsValidTransactionSignature is a free data retrieval call binding the contract method 0x8d45cd23.
//
// Solidity: function isValidTransactionSignature(LibZeroExTransactionZeroExTransaction transaction, bytes signature) constant returns(bool isValid)
func (_Exchangecontract *ExchangecontractCaller) IsValidTransactionSignature(opts *bind.CallOpts, transaction LibZeroExTransactionZeroExTransaction, signature []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "isValidTransactionSignature", transaction, signature)
	return *ret0, err
}

// IsValidTransactionSignature is a free data retrieval call binding the contract method 0x8d45cd23.
//
// Solidity: function isValidTransactionSignature(LibZeroExTransactionZeroExTransaction transaction, bytes signature) constant returns(bool isValid)
func (_Exchangecontract *ExchangecontractSession) IsValidTransactionSignature(transaction LibZeroExTransactionZeroExTransaction, signature []byte) (bool, error) {
	return _Exchangecontract.Contract.IsValidTransactionSignature(&_Exchangecontract.CallOpts, transaction, signature)
}

// IsValidTransactionSignature is a free data retrieval call binding the contract method 0x8d45cd23.
//
// Solidity: function isValidTransactionSignature(LibZeroExTransactionZeroExTransaction transaction, bytes signature) constant returns(bool isValid)
func (_Exchangecontract *ExchangecontractCallerSession) IsValidTransactionSignature(transaction LibZeroExTransactionZeroExTransaction, signature []byte) (bool, error) {
	return _Exchangecontract.Contract.IsValidTransactionSignature(&_Exchangecontract.CallOpts, transaction, signature)
}

// OrderEpoch is a free data retrieval call binding the contract method 0xd9bfa73e.
//
// Solidity: function orderEpoch(address , address ) constant returns(uint256)
func (_Exchangecontract *ExchangecontractCaller) OrderEpoch(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "orderEpoch", arg0, arg1)
	return *ret0, err
}

// OrderEpoch is a free data retrieval call binding the contract method 0xd9bfa73e.
//
// Solidity: function orderEpoch(address , address ) constant returns(uint256)
func (_Exchangecontract *ExchangecontractSession) OrderEpoch(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Exchangecontract.Contract.OrderEpoch(&_Exchangecontract.CallOpts, arg0, arg1)
}

// OrderEpoch is a free data retrieval call binding the contract method 0xd9bfa73e.
//
// Solidity: function orderEpoch(address , address ) constant returns(uint256)
func (_Exchangecontract *ExchangecontractCallerSession) OrderEpoch(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Exchangecontract.Contract.OrderEpoch(&_Exchangecontract.CallOpts, arg0, arg1)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Exchangecontract *ExchangecontractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Exchangecontract *ExchangecontractSession) Owner() (common.Address, error) {
	return _Exchangecontract.Contract.Owner(&_Exchangecontract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Exchangecontract *ExchangecontractCallerSession) Owner() (common.Address, error) {
	return _Exchangecontract.Contract.Owner(&_Exchangecontract.CallOpts)
}

// PreSigned is a free data retrieval call binding the contract method 0x82c174d0.
//
// Solidity: function preSigned(bytes32 , address ) constant returns(bool)
func (_Exchangecontract *ExchangecontractCaller) PreSigned(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "preSigned", arg0, arg1)
	return *ret0, err
}

// PreSigned is a free data retrieval call binding the contract method 0x82c174d0.
//
// Solidity: function preSigned(bytes32 , address ) constant returns(bool)
func (_Exchangecontract *ExchangecontractSession) PreSigned(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _Exchangecontract.Contract.PreSigned(&_Exchangecontract.CallOpts, arg0, arg1)
}

// PreSigned is a free data retrieval call binding the contract method 0x82c174d0.
//
// Solidity: function preSigned(bytes32 , address ) constant returns(bool)
func (_Exchangecontract *ExchangecontractCallerSession) PreSigned(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _Exchangecontract.Contract.PreSigned(&_Exchangecontract.CallOpts, arg0, arg1)
}

// ProtocolFeeCollector is a free data retrieval call binding the contract method 0x850a1501.
//
// Solidity: function protocolFeeCollector() constant returns(address)
func (_Exchangecontract *ExchangecontractCaller) ProtocolFeeCollector(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "protocolFeeCollector")
	return *ret0, err
}

// ProtocolFeeCollector is a free data retrieval call binding the contract method 0x850a1501.
//
// Solidity: function protocolFeeCollector() constant returns(address)
func (_Exchangecontract *ExchangecontractSession) ProtocolFeeCollector() (common.Address, error) {
	return _Exchangecontract.Contract.ProtocolFeeCollector(&_Exchangecontract.CallOpts)
}

// ProtocolFeeCollector is a free data retrieval call binding the contract method 0x850a1501.
//
// Solidity: function protocolFeeCollector() constant returns(address)
func (_Exchangecontract *ExchangecontractCallerSession) ProtocolFeeCollector() (common.Address, error) {
	return _Exchangecontract.Contract.ProtocolFeeCollector(&_Exchangecontract.CallOpts)
}

// ProtocolFeeMultiplier is a free data retrieval call binding the contract method 0x1ce4c78b.
//
// Solidity: function protocolFeeMultiplier() constant returns(uint256)
func (_Exchangecontract *ExchangecontractCaller) ProtocolFeeMultiplier(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "protocolFeeMultiplier")
	return *ret0, err
}

// ProtocolFeeMultiplier is a free data retrieval call binding the contract method 0x1ce4c78b.
//
// Solidity: function protocolFeeMultiplier() constant returns(uint256)
func (_Exchangecontract *ExchangecontractSession) ProtocolFeeMultiplier() (*big.Int, error) {
	return _Exchangecontract.Contract.ProtocolFeeMultiplier(&_Exchangecontract.CallOpts)
}

// ProtocolFeeMultiplier is a free data retrieval call binding the contract method 0x1ce4c78b.
//
// Solidity: function protocolFeeMultiplier() constant returns(uint256)
func (_Exchangecontract *ExchangecontractCallerSession) ProtocolFeeMultiplier() (*big.Int, error) {
	return _Exchangecontract.Contract.ProtocolFeeMultiplier(&_Exchangecontract.CallOpts)
}

// TransactionsExecuted is a free data retrieval call binding the contract method 0x0228e168.
//
// Solidity: function transactionsExecuted(bytes32 ) constant returns(bool)
func (_Exchangecontract *ExchangecontractCaller) TransactionsExecuted(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchangecontract.contract.Call(opts, out, "transactionsExecuted", arg0)
	return *ret0, err
}

// TransactionsExecuted is a free data retrieval call binding the contract method 0x0228e168.
//
// Solidity: function transactionsExecuted(bytes32 ) constant returns(bool)
func (_Exchangecontract *ExchangecontractSession) TransactionsExecuted(arg0 [32]byte) (bool, error) {
	return _Exchangecontract.Contract.TransactionsExecuted(&_Exchangecontract.CallOpts, arg0)
}

// TransactionsExecuted is a free data retrieval call binding the contract method 0x0228e168.
//
// Solidity: function transactionsExecuted(bytes32 ) constant returns(bool)
func (_Exchangecontract *ExchangecontractCallerSession) TransactionsExecuted(arg0 [32]byte) (bool, error) {
	return _Exchangecontract.Contract.TransactionsExecuted(&_Exchangecontract.CallOpts, arg0)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xdedfc1f1.
//
// Solidity: function batchCancelOrders([]LibOrderOrder orders) returns()
func (_Exchangecontract *ExchangecontractTransactor) BatchCancelOrders(opts *bind.TransactOpts, orders []LibOrderOrder) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "batchCancelOrders", orders)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xdedfc1f1.
//
// Solidity: function batchCancelOrders([]LibOrderOrder orders) returns()
func (_Exchangecontract *ExchangecontractSession) BatchCancelOrders(orders []LibOrderOrder) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchCancelOrders(&_Exchangecontract.TransactOpts, orders)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xdedfc1f1.
//
// Solidity: function batchCancelOrders([]LibOrderOrder orders) returns()
func (_Exchangecontract *ExchangecontractTransactorSession) BatchCancelOrders(orders []LibOrderOrder) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchCancelOrders(&_Exchangecontract.TransactOpts, orders)
}

// BatchExecuteTransactions is a paid mutator transaction binding the contract method 0xfc74896d.
//
// Solidity: function batchExecuteTransactions([]LibZeroExTransactionZeroExTransaction transactions, bytes[] signatures) returns(bytes[] returnData)
func (_Exchangecontract *ExchangecontractTransactor) BatchExecuteTransactions(opts *bind.TransactOpts, transactions []LibZeroExTransactionZeroExTransaction, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "batchExecuteTransactions", transactions, signatures)
}

// BatchExecuteTransactions is a paid mutator transaction binding the contract method 0xfc74896d.
//
// Solidity: function batchExecuteTransactions([]LibZeroExTransactionZeroExTransaction transactions, bytes[] signatures) returns(bytes[] returnData)
func (_Exchangecontract *ExchangecontractSession) BatchExecuteTransactions(transactions []LibZeroExTransactionZeroExTransaction, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchExecuteTransactions(&_Exchangecontract.TransactOpts, transactions, signatures)
}

// BatchExecuteTransactions is a paid mutator transaction binding the contract method 0xfc74896d.
//
// Solidity: function batchExecuteTransactions([]LibZeroExTransactionZeroExTransaction transactions, bytes[] signatures) returns(bytes[] returnData)
func (_Exchangecontract *ExchangecontractTransactorSession) BatchExecuteTransactions(transactions []LibZeroExTransactionZeroExTransaction, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchExecuteTransactions(&_Exchangecontract.TransactOpts, transactions, signatures)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0xbeee2e14.
//
// Solidity: function batchFillOrKillOrders([]LibOrderOrder orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns([]LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactor) BatchFillOrKillOrders(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "batchFillOrKillOrders", orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0xbeee2e14.
//
// Solidity: function batchFillOrKillOrders([]LibOrderOrder orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns([]LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractSession) BatchFillOrKillOrders(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchFillOrKillOrders(&_Exchangecontract.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0xbeee2e14.
//
// Solidity: function batchFillOrKillOrders([]LibOrderOrder orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns([]LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactorSession) BatchFillOrKillOrders(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchFillOrKillOrders(&_Exchangecontract.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0x9694a402.
//
// Solidity: function batchFillOrders([]LibOrderOrder orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns([]LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactor) BatchFillOrders(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "batchFillOrders", orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0x9694a402.
//
// Solidity: function batchFillOrders([]LibOrderOrder orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns([]LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractSession) BatchFillOrders(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchFillOrders(&_Exchangecontract.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0x9694a402.
//
// Solidity: function batchFillOrders([]LibOrderOrder orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns([]LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactorSession) BatchFillOrders(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchFillOrders(&_Exchangecontract.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrdersNoThrow is a paid mutator transaction binding the contract method 0x8ea8dfe4.
//
// Solidity: function batchFillOrdersNoThrow([]LibOrderOrder orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns([]LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactor) BatchFillOrdersNoThrow(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "batchFillOrdersNoThrow", orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrdersNoThrow is a paid mutator transaction binding the contract method 0x8ea8dfe4.
//
// Solidity: function batchFillOrdersNoThrow([]LibOrderOrder orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns([]LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractSession) BatchFillOrdersNoThrow(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchFillOrdersNoThrow(&_Exchangecontract.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrdersNoThrow is a paid mutator transaction binding the contract method 0x8ea8dfe4.
//
// Solidity: function batchFillOrdersNoThrow([]LibOrderOrder orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns([]LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactorSession) BatchFillOrdersNoThrow(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchFillOrdersNoThrow(&_Exchangecontract.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchMatchOrders is a paid mutator transaction binding the contract method 0x6fcf3e9e.
//
// Solidity: function batchMatchOrders([]LibOrderOrder leftOrders, []LibOrderOrder rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) returns(LibFillResultsBatchMatchedFillResults batchMatchedFillResults)
func (_Exchangecontract *ExchangecontractTransactor) BatchMatchOrders(opts *bind.TransactOpts, leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "batchMatchOrders", leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrders is a paid mutator transaction binding the contract method 0x6fcf3e9e.
//
// Solidity: function batchMatchOrders([]LibOrderOrder leftOrders, []LibOrderOrder rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) returns(LibFillResultsBatchMatchedFillResults batchMatchedFillResults)
func (_Exchangecontract *ExchangecontractSession) BatchMatchOrders(leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchMatchOrders(&_Exchangecontract.TransactOpts, leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrders is a paid mutator transaction binding the contract method 0x6fcf3e9e.
//
// Solidity: function batchMatchOrders([]LibOrderOrder leftOrders, []LibOrderOrder rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) returns(LibFillResultsBatchMatchedFillResults batchMatchedFillResults)
func (_Exchangecontract *ExchangecontractTransactorSession) BatchMatchOrders(leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchMatchOrders(&_Exchangecontract.TransactOpts, leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0x6a1a80fd.
//
// Solidity: function batchMatchOrdersWithMaximalFill([]LibOrderOrder leftOrders, []LibOrderOrder rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) returns(LibFillResultsBatchMatchedFillResults batchMatchedFillResults)
func (_Exchangecontract *ExchangecontractTransactor) BatchMatchOrdersWithMaximalFill(opts *bind.TransactOpts, leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "batchMatchOrdersWithMaximalFill", leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0x6a1a80fd.
//
// Solidity: function batchMatchOrdersWithMaximalFill([]LibOrderOrder leftOrders, []LibOrderOrder rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) returns(LibFillResultsBatchMatchedFillResults batchMatchedFillResults)
func (_Exchangecontract *ExchangecontractSession) BatchMatchOrdersWithMaximalFill(leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchMatchOrdersWithMaximalFill(&_Exchangecontract.TransactOpts, leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0x6a1a80fd.
//
// Solidity: function batchMatchOrdersWithMaximalFill([]LibOrderOrder leftOrders, []LibOrderOrder rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) returns(LibFillResultsBatchMatchedFillResults batchMatchedFillResults)
func (_Exchangecontract *ExchangecontractTransactorSession) BatchMatchOrdersWithMaximalFill(leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.BatchMatchOrdersWithMaximalFill(&_Exchangecontract.TransactOpts, leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x2da62987.
//
// Solidity: function cancelOrder(LibOrderOrder order) returns()
func (_Exchangecontract *ExchangecontractTransactor) CancelOrder(opts *bind.TransactOpts, order LibOrderOrder) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "cancelOrder", order)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x2da62987.
//
// Solidity: function cancelOrder(LibOrderOrder order) returns()
func (_Exchangecontract *ExchangecontractSession) CancelOrder(order LibOrderOrder) (*types.Transaction, error) {
	return _Exchangecontract.Contract.CancelOrder(&_Exchangecontract.TransactOpts, order)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x2da62987.
//
// Solidity: function cancelOrder(LibOrderOrder order) returns()
func (_Exchangecontract *ExchangecontractTransactorSession) CancelOrder(order LibOrderOrder) (*types.Transaction, error) {
	return _Exchangecontract.Contract.CancelOrder(&_Exchangecontract.TransactOpts, order)
}

// CancelOrdersUpTo is a paid mutator transaction binding the contract method 0x4f9559b1.
//
// Solidity: function cancelOrdersUpTo(uint256 targetOrderEpoch) returns()
func (_Exchangecontract *ExchangecontractTransactor) CancelOrdersUpTo(opts *bind.TransactOpts, targetOrderEpoch *big.Int) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "cancelOrdersUpTo", targetOrderEpoch)
}

// CancelOrdersUpTo is a paid mutator transaction binding the contract method 0x4f9559b1.
//
// Solidity: function cancelOrdersUpTo(uint256 targetOrderEpoch) returns()
func (_Exchangecontract *ExchangecontractSession) CancelOrdersUpTo(targetOrderEpoch *big.Int) (*types.Transaction, error) {
	return _Exchangecontract.Contract.CancelOrdersUpTo(&_Exchangecontract.TransactOpts, targetOrderEpoch)
}

// CancelOrdersUpTo is a paid mutator transaction binding the contract method 0x4f9559b1.
//
// Solidity: function cancelOrdersUpTo(uint256 targetOrderEpoch) returns()
func (_Exchangecontract *ExchangecontractTransactorSession) CancelOrdersUpTo(targetOrderEpoch *big.Int) (*types.Transaction, error) {
	return _Exchangecontract.Contract.CancelOrdersUpTo(&_Exchangecontract.TransactOpts, targetOrderEpoch)
}

// DetachProtocolFeeCollector is a paid mutator transaction binding the contract method 0x0efca185.
//
// Solidity: function detachProtocolFeeCollector() returns()
func (_Exchangecontract *ExchangecontractTransactor) DetachProtocolFeeCollector(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "detachProtocolFeeCollector")
}

// DetachProtocolFeeCollector is a paid mutator transaction binding the contract method 0x0efca185.
//
// Solidity: function detachProtocolFeeCollector() returns()
func (_Exchangecontract *ExchangecontractSession) DetachProtocolFeeCollector() (*types.Transaction, error) {
	return _Exchangecontract.Contract.DetachProtocolFeeCollector(&_Exchangecontract.TransactOpts)
}

// DetachProtocolFeeCollector is a paid mutator transaction binding the contract method 0x0efca185.
//
// Solidity: function detachProtocolFeeCollector() returns()
func (_Exchangecontract *ExchangecontractTransactorSession) DetachProtocolFeeCollector() (*types.Transaction, error) {
	return _Exchangecontract.Contract.DetachProtocolFeeCollector(&_Exchangecontract.TransactOpts)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0x2280c910.
//
// Solidity: function executeTransaction(LibZeroExTransactionZeroExTransaction transaction, bytes signature) returns(bytes)
func (_Exchangecontract *ExchangecontractTransactor) ExecuteTransaction(opts *bind.TransactOpts, transaction LibZeroExTransactionZeroExTransaction, signature []byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "executeTransaction", transaction, signature)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0x2280c910.
//
// Solidity: function executeTransaction(LibZeroExTransactionZeroExTransaction transaction, bytes signature) returns(bytes)
func (_Exchangecontract *ExchangecontractSession) ExecuteTransaction(transaction LibZeroExTransactionZeroExTransaction, signature []byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.ExecuteTransaction(&_Exchangecontract.TransactOpts, transaction, signature)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0x2280c910.
//
// Solidity: function executeTransaction(LibZeroExTransactionZeroExTransaction transaction, bytes signature) returns(bytes)
func (_Exchangecontract *ExchangecontractTransactorSession) ExecuteTransaction(transaction LibZeroExTransactionZeroExTransaction, signature []byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.ExecuteTransaction(&_Exchangecontract.TransactOpts, transaction, signature)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0xe14b58c4.
//
// Solidity: function fillOrKillOrder(LibOrderOrder order, uint256 takerAssetFillAmount, bytes signature) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactor) FillOrKillOrder(opts *bind.TransactOpts, order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "fillOrKillOrder", order, takerAssetFillAmount, signature)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0xe14b58c4.
//
// Solidity: function fillOrKillOrder(LibOrderOrder order, uint256 takerAssetFillAmount, bytes signature) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractSession) FillOrKillOrder(order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.FillOrKillOrder(&_Exchangecontract.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0xe14b58c4.
//
// Solidity: function fillOrKillOrder(LibOrderOrder order, uint256 takerAssetFillAmount, bytes signature) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactorSession) FillOrKillOrder(order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.FillOrKillOrder(&_Exchangecontract.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrder is a paid mutator transaction binding the contract method 0x9b44d556.
//
// Solidity: function fillOrder(LibOrderOrder order, uint256 takerAssetFillAmount, bytes signature) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactor) FillOrder(opts *bind.TransactOpts, order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "fillOrder", order, takerAssetFillAmount, signature)
}

// FillOrder is a paid mutator transaction binding the contract method 0x9b44d556.
//
// Solidity: function fillOrder(LibOrderOrder order, uint256 takerAssetFillAmount, bytes signature) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractSession) FillOrder(order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.FillOrder(&_Exchangecontract.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrder is a paid mutator transaction binding the contract method 0x9b44d556.
//
// Solidity: function fillOrder(LibOrderOrder order, uint256 takerAssetFillAmount, bytes signature) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactorSession) FillOrder(order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.FillOrder(&_Exchangecontract.TransactOpts, order, takerAssetFillAmount, signature)
}

// MarketBuyOrdersFillOrKill is a paid mutator transaction binding the contract method 0x8bc8efb3.
//
// Solidity: function marketBuyOrdersFillOrKill([]LibOrderOrder orders, uint256 makerAssetFillAmount, bytes[] signatures) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactor) MarketBuyOrdersFillOrKill(opts *bind.TransactOpts, orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "marketBuyOrdersFillOrKill", orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersFillOrKill is a paid mutator transaction binding the contract method 0x8bc8efb3.
//
// Solidity: function marketBuyOrdersFillOrKill([]LibOrderOrder orders, uint256 makerAssetFillAmount, bytes[] signatures) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractSession) MarketBuyOrdersFillOrKill(orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.MarketBuyOrdersFillOrKill(&_Exchangecontract.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersFillOrKill is a paid mutator transaction binding the contract method 0x8bc8efb3.
//
// Solidity: function marketBuyOrdersFillOrKill([]LibOrderOrder orders, uint256 makerAssetFillAmount, bytes[] signatures) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactorSession) MarketBuyOrdersFillOrKill(orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.MarketBuyOrdersFillOrKill(&_Exchangecontract.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersNoThrow is a paid mutator transaction binding the contract method 0x78d29ac1.
//
// Solidity: function marketBuyOrdersNoThrow([]LibOrderOrder orders, uint256 makerAssetFillAmount, bytes[] signatures) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactor) MarketBuyOrdersNoThrow(opts *bind.TransactOpts, orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "marketBuyOrdersNoThrow", orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersNoThrow is a paid mutator transaction binding the contract method 0x78d29ac1.
//
// Solidity: function marketBuyOrdersNoThrow([]LibOrderOrder orders, uint256 makerAssetFillAmount, bytes[] signatures) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractSession) MarketBuyOrdersNoThrow(orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.MarketBuyOrdersNoThrow(&_Exchangecontract.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersNoThrow is a paid mutator transaction binding the contract method 0x78d29ac1.
//
// Solidity: function marketBuyOrdersNoThrow([]LibOrderOrder orders, uint256 makerAssetFillAmount, bytes[] signatures) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactorSession) MarketBuyOrdersNoThrow(orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.MarketBuyOrdersNoThrow(&_Exchangecontract.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketSellOrdersFillOrKill is a paid mutator transaction binding the contract method 0xa6c3bf33.
//
// Solidity: function marketSellOrdersFillOrKill([]LibOrderOrder orders, uint256 takerAssetFillAmount, bytes[] signatures) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactor) MarketSellOrdersFillOrKill(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "marketSellOrdersFillOrKill", orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersFillOrKill is a paid mutator transaction binding the contract method 0xa6c3bf33.
//
// Solidity: function marketSellOrdersFillOrKill([]LibOrderOrder orders, uint256 takerAssetFillAmount, bytes[] signatures) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractSession) MarketSellOrdersFillOrKill(orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.MarketSellOrdersFillOrKill(&_Exchangecontract.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersFillOrKill is a paid mutator transaction binding the contract method 0xa6c3bf33.
//
// Solidity: function marketSellOrdersFillOrKill([]LibOrderOrder orders, uint256 takerAssetFillAmount, bytes[] signatures) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactorSession) MarketSellOrdersFillOrKill(orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.MarketSellOrdersFillOrKill(&_Exchangecontract.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersNoThrow is a paid mutator transaction binding the contract method 0x369da099.
//
// Solidity: function marketSellOrdersNoThrow([]LibOrderOrder orders, uint256 takerAssetFillAmount, bytes[] signatures) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactor) MarketSellOrdersNoThrow(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "marketSellOrdersNoThrow", orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersNoThrow is a paid mutator transaction binding the contract method 0x369da099.
//
// Solidity: function marketSellOrdersNoThrow([]LibOrderOrder orders, uint256 takerAssetFillAmount, bytes[] signatures) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractSession) MarketSellOrdersNoThrow(orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.MarketSellOrdersNoThrow(&_Exchangecontract.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersNoThrow is a paid mutator transaction binding the contract method 0x369da099.
//
// Solidity: function marketSellOrdersNoThrow([]LibOrderOrder orders, uint256 takerAssetFillAmount, bytes[] signatures) returns(LibFillResultsFillResults fillResults)
func (_Exchangecontract *ExchangecontractTransactorSession) MarketSellOrdersNoThrow(orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.MarketSellOrdersNoThrow(&_Exchangecontract.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MatchOrders is a paid mutator transaction binding the contract method 0x88ec79fb.
//
// Solidity: function matchOrders(LibOrderOrder leftOrder, LibOrderOrder rightOrder, bytes leftSignature, bytes rightSignature) returns(LibFillResultsMatchedFillResults matchedFillResults)
func (_Exchangecontract *ExchangecontractTransactor) MatchOrders(opts *bind.TransactOpts, leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "matchOrders", leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrders is a paid mutator transaction binding the contract method 0x88ec79fb.
//
// Solidity: function matchOrders(LibOrderOrder leftOrder, LibOrderOrder rightOrder, bytes leftSignature, bytes rightSignature) returns(LibFillResultsMatchedFillResults matchedFillResults)
func (_Exchangecontract *ExchangecontractSession) MatchOrders(leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.MatchOrders(&_Exchangecontract.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrders is a paid mutator transaction binding the contract method 0x88ec79fb.
//
// Solidity: function matchOrders(LibOrderOrder leftOrder, LibOrderOrder rightOrder, bytes leftSignature, bytes rightSignature) returns(LibFillResultsMatchedFillResults matchedFillResults)
func (_Exchangecontract *ExchangecontractTransactorSession) MatchOrders(leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.MatchOrders(&_Exchangecontract.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0xb718e292.
//
// Solidity: function matchOrdersWithMaximalFill(LibOrderOrder leftOrder, LibOrderOrder rightOrder, bytes leftSignature, bytes rightSignature) returns(LibFillResultsMatchedFillResults matchedFillResults)
func (_Exchangecontract *ExchangecontractTransactor) MatchOrdersWithMaximalFill(opts *bind.TransactOpts, leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "matchOrdersWithMaximalFill", leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0xb718e292.
//
// Solidity: function matchOrdersWithMaximalFill(LibOrderOrder leftOrder, LibOrderOrder rightOrder, bytes leftSignature, bytes rightSignature) returns(LibFillResultsMatchedFillResults matchedFillResults)
func (_Exchangecontract *ExchangecontractSession) MatchOrdersWithMaximalFill(leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.MatchOrdersWithMaximalFill(&_Exchangecontract.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0xb718e292.
//
// Solidity: function matchOrdersWithMaximalFill(LibOrderOrder leftOrder, LibOrderOrder rightOrder, bytes leftSignature, bytes rightSignature) returns(LibFillResultsMatchedFillResults matchedFillResults)
func (_Exchangecontract *ExchangecontractTransactorSession) MatchOrdersWithMaximalFill(leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.MatchOrdersWithMaximalFill(&_Exchangecontract.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// PreSign is a paid mutator transaction binding the contract method 0x46c02d7a.
//
// Solidity: function preSign(bytes32 hash) returns()
func (_Exchangecontract *ExchangecontractTransactor) PreSign(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "preSign", hash)
}

// PreSign is a paid mutator transaction binding the contract method 0x46c02d7a.
//
// Solidity: function preSign(bytes32 hash) returns()
func (_Exchangecontract *ExchangecontractSession) PreSign(hash [32]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.PreSign(&_Exchangecontract.TransactOpts, hash)
}

// PreSign is a paid mutator transaction binding the contract method 0x46c02d7a.
//
// Solidity: function preSign(bytes32 hash) returns()
func (_Exchangecontract *ExchangecontractTransactorSession) PreSign(hash [32]byte) (*types.Transaction, error) {
	return _Exchangecontract.Contract.PreSign(&_Exchangecontract.TransactOpts, hash)
}

// RegisterAssetProxy is a paid mutator transaction binding the contract method 0xc585bb93.
//
// Solidity: function registerAssetProxy(address assetProxy) returns()
func (_Exchangecontract *ExchangecontractTransactor) RegisterAssetProxy(opts *bind.TransactOpts, assetProxy common.Address) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "registerAssetProxy", assetProxy)
}

// RegisterAssetProxy is a paid mutator transaction binding the contract method 0xc585bb93.
//
// Solidity: function registerAssetProxy(address assetProxy) returns()
func (_Exchangecontract *ExchangecontractSession) RegisterAssetProxy(assetProxy common.Address) (*types.Transaction, error) {
	return _Exchangecontract.Contract.RegisterAssetProxy(&_Exchangecontract.TransactOpts, assetProxy)
}

// RegisterAssetProxy is a paid mutator transaction binding the contract method 0xc585bb93.
//
// Solidity: function registerAssetProxy(address assetProxy) returns()
func (_Exchangecontract *ExchangecontractTransactorSession) RegisterAssetProxy(assetProxy common.Address) (*types.Transaction, error) {
	return _Exchangecontract.Contract.RegisterAssetProxy(&_Exchangecontract.TransactOpts, assetProxy)
}

// SetProtocolFeeCollectorAddress is a paid mutator transaction binding the contract method 0xc0fa16cc.
//
// Solidity: function setProtocolFeeCollectorAddress(address updatedProtocolFeeCollector) returns()
func (_Exchangecontract *ExchangecontractTransactor) SetProtocolFeeCollectorAddress(opts *bind.TransactOpts, updatedProtocolFeeCollector common.Address) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "setProtocolFeeCollectorAddress", updatedProtocolFeeCollector)
}

// SetProtocolFeeCollectorAddress is a paid mutator transaction binding the contract method 0xc0fa16cc.
//
// Solidity: function setProtocolFeeCollectorAddress(address updatedProtocolFeeCollector) returns()
func (_Exchangecontract *ExchangecontractSession) SetProtocolFeeCollectorAddress(updatedProtocolFeeCollector common.Address) (*types.Transaction, error) {
	return _Exchangecontract.Contract.SetProtocolFeeCollectorAddress(&_Exchangecontract.TransactOpts, updatedProtocolFeeCollector)
}

// SetProtocolFeeCollectorAddress is a paid mutator transaction binding the contract method 0xc0fa16cc.
//
// Solidity: function setProtocolFeeCollectorAddress(address updatedProtocolFeeCollector) returns()
func (_Exchangecontract *ExchangecontractTransactorSession) SetProtocolFeeCollectorAddress(updatedProtocolFeeCollector common.Address) (*types.Transaction, error) {
	return _Exchangecontract.Contract.SetProtocolFeeCollectorAddress(&_Exchangecontract.TransactOpts, updatedProtocolFeeCollector)
}

// SetProtocolFeeMultiplier is a paid mutator transaction binding the contract method 0x9331c742.
//
// Solidity: function setProtocolFeeMultiplier(uint256 updatedProtocolFeeMultiplier) returns()
func (_Exchangecontract *ExchangecontractTransactor) SetProtocolFeeMultiplier(opts *bind.TransactOpts, updatedProtocolFeeMultiplier *big.Int) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "setProtocolFeeMultiplier", updatedProtocolFeeMultiplier)
}

// SetProtocolFeeMultiplier is a paid mutator transaction binding the contract method 0x9331c742.
//
// Solidity: function setProtocolFeeMultiplier(uint256 updatedProtocolFeeMultiplier) returns()
func (_Exchangecontract *ExchangecontractSession) SetProtocolFeeMultiplier(updatedProtocolFeeMultiplier *big.Int) (*types.Transaction, error) {
	return _Exchangecontract.Contract.SetProtocolFeeMultiplier(&_Exchangecontract.TransactOpts, updatedProtocolFeeMultiplier)
}

// SetProtocolFeeMultiplier is a paid mutator transaction binding the contract method 0x9331c742.
//
// Solidity: function setProtocolFeeMultiplier(uint256 updatedProtocolFeeMultiplier) returns()
func (_Exchangecontract *ExchangecontractTransactorSession) SetProtocolFeeMultiplier(updatedProtocolFeeMultiplier *big.Int) (*types.Transaction, error) {
	return _Exchangecontract.Contract.SetProtocolFeeMultiplier(&_Exchangecontract.TransactOpts, updatedProtocolFeeMultiplier)
}

// SetSignatureValidatorApproval is a paid mutator transaction binding the contract method 0x77fcce68.
//
// Solidity: function setSignatureValidatorApproval(address validatorAddress, bool approval) returns()
func (_Exchangecontract *ExchangecontractTransactor) SetSignatureValidatorApproval(opts *bind.TransactOpts, validatorAddress common.Address, approval bool) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "setSignatureValidatorApproval", validatorAddress, approval)
}

// SetSignatureValidatorApproval is a paid mutator transaction binding the contract method 0x77fcce68.
//
// Solidity: function setSignatureValidatorApproval(address validatorAddress, bool approval) returns()
func (_Exchangecontract *ExchangecontractSession) SetSignatureValidatorApproval(validatorAddress common.Address, approval bool) (*types.Transaction, error) {
	return _Exchangecontract.Contract.SetSignatureValidatorApproval(&_Exchangecontract.TransactOpts, validatorAddress, approval)
}

// SetSignatureValidatorApproval is a paid mutator transaction binding the contract method 0x77fcce68.
//
// Solidity: function setSignatureValidatorApproval(address validatorAddress, bool approval) returns()
func (_Exchangecontract *ExchangecontractTransactorSession) SetSignatureValidatorApproval(validatorAddress common.Address, approval bool) (*types.Transaction, error) {
	return _Exchangecontract.Contract.SetSignatureValidatorApproval(&_Exchangecontract.TransactOpts, validatorAddress, approval)
}

// SimulateDispatchTransferFromCalls is a paid mutator transaction binding the contract method 0xb04fbddd.
//
// Solidity: function simulateDispatchTransferFromCalls(bytes[] assetData, address[] fromAddresses, address[] toAddresses, uint256[] amounts) returns()
func (_Exchangecontract *ExchangecontractTransactor) SimulateDispatchTransferFromCalls(opts *bind.TransactOpts, assetData [][]byte, fromAddresses []common.Address, toAddresses []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "simulateDispatchTransferFromCalls", assetData, fromAddresses, toAddresses, amounts)
}

// SimulateDispatchTransferFromCalls is a paid mutator transaction binding the contract method 0xb04fbddd.
//
// Solidity: function simulateDispatchTransferFromCalls(bytes[] assetData, address[] fromAddresses, address[] toAddresses, uint256[] amounts) returns()
func (_Exchangecontract *ExchangecontractSession) SimulateDispatchTransferFromCalls(assetData [][]byte, fromAddresses []common.Address, toAddresses []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Exchangecontract.Contract.SimulateDispatchTransferFromCalls(&_Exchangecontract.TransactOpts, assetData, fromAddresses, toAddresses, amounts)
}

// SimulateDispatchTransferFromCalls is a paid mutator transaction binding the contract method 0xb04fbddd.
//
// Solidity: function simulateDispatchTransferFromCalls(bytes[] assetData, address[] fromAddresses, address[] toAddresses, uint256[] amounts) returns()
func (_Exchangecontract *ExchangecontractTransactorSession) SimulateDispatchTransferFromCalls(assetData [][]byte, fromAddresses []common.Address, toAddresses []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Exchangecontract.Contract.SimulateDispatchTransferFromCalls(&_Exchangecontract.TransactOpts, assetData, fromAddresses, toAddresses, amounts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Exchangecontract *ExchangecontractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Exchangecontract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Exchangecontract *ExchangecontractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Exchangecontract.Contract.TransferOwnership(&_Exchangecontract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Exchangecontract *ExchangecontractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Exchangecontract.Contract.TransferOwnership(&_Exchangecontract.TransactOpts, newOwner)
}

// ExchangecontractAssetProxyRegisteredIterator is returned from FilterAssetProxyRegistered and is used to iterate over the raw logs and unpacked data for AssetProxyRegistered events raised by the Exchangecontract contract.
type ExchangecontractAssetProxyRegisteredIterator struct {
	Event *ExchangecontractAssetProxyRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangecontractAssetProxyRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangecontractAssetProxyRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangecontractAssetProxyRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangecontractAssetProxyRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangecontractAssetProxyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangecontractAssetProxyRegistered represents a AssetProxyRegistered event raised by the Exchangecontract contract.
type ExchangecontractAssetProxyRegistered struct {
	Id         [4]byte
	AssetProxy common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAssetProxyRegistered is a free log retrieval operation binding the contract event 0xd2c6b762299c609bdb96520b58a49bfb80186934d4f71a86a367571a15c03194.
//
// Solidity: event AssetProxyRegistered(bytes4 id, address assetProxy)
func (_Exchangecontract *ExchangecontractFilterer) FilterAssetProxyRegistered(opts *bind.FilterOpts) (*ExchangecontractAssetProxyRegisteredIterator, error) {

	logs, sub, err := _Exchangecontract.contract.FilterLogs(opts, "AssetProxyRegistered")
	if err != nil {
		return nil, err
	}
	return &ExchangecontractAssetProxyRegisteredIterator{contract: _Exchangecontract.contract, event: "AssetProxyRegistered", logs: logs, sub: sub}, nil
}

// WatchAssetProxyRegistered is a free log subscription operation binding the contract event 0xd2c6b762299c609bdb96520b58a49bfb80186934d4f71a86a367571a15c03194.
//
// Solidity: event AssetProxyRegistered(bytes4 id, address assetProxy)
func (_Exchangecontract *ExchangecontractFilterer) WatchAssetProxyRegistered(opts *bind.WatchOpts, sink chan<- *ExchangecontractAssetProxyRegistered) (event.Subscription, error) {

	logs, sub, err := _Exchangecontract.contract.WatchLogs(opts, "AssetProxyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangecontractAssetProxyRegistered)
				if err := _Exchangecontract.contract.UnpackLog(event, "AssetProxyRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAssetProxyRegistered is a log parse operation binding the contract event 0xd2c6b762299c609bdb96520b58a49bfb80186934d4f71a86a367571a15c03194.
//
// Solidity: event AssetProxyRegistered(bytes4 id, address assetProxy)
func (_Exchangecontract *ExchangecontractFilterer) ParseAssetProxyRegistered(log types.Log) (*ExchangecontractAssetProxyRegistered, error) {
	event := new(ExchangecontractAssetProxyRegistered)
	if err := _Exchangecontract.contract.UnpackLog(event, "AssetProxyRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangecontractCancelIterator is returned from FilterCancel and is used to iterate over the raw logs and unpacked data for Cancel events raised by the Exchangecontract contract.
type ExchangecontractCancelIterator struct {
	Event *ExchangecontractCancel // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangecontractCancelIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangecontractCancel)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangecontractCancel)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangecontractCancelIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangecontractCancelIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangecontractCancel represents a Cancel event raised by the Exchangecontract contract.
type ExchangecontractCancel struct {
	MakerAddress        common.Address
	FeeRecipientAddress common.Address
	MakerAssetData      []byte
	TakerAssetData      []byte
	SenderAddress       common.Address
	OrderHash           [32]byte
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterCancel is a free log retrieval operation binding the contract event 0x02c310a9a43963ff31a754a4099cc435ed498049687539d72d7818d9b093415c.
//
// Solidity: event Cancel(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, address senderAddress, bytes32 indexed orderHash)
func (_Exchangecontract *ExchangecontractFilterer) FilterCancel(opts *bind.FilterOpts, makerAddress []common.Address, feeRecipientAddress []common.Address, orderHash [][32]byte) (*ExchangecontractCancelIterator, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var feeRecipientAddressRule []interface{}
	for _, feeRecipientAddressItem := range feeRecipientAddress {
		feeRecipientAddressRule = append(feeRecipientAddressRule, feeRecipientAddressItem)
	}

	var orderHashRule []interface{}
	for _, orderHashItem := range orderHash {
		orderHashRule = append(orderHashRule, orderHashItem)
	}

	logs, sub, err := _Exchangecontract.contract.FilterLogs(opts, "Cancel", makerAddressRule, feeRecipientAddressRule, orderHashRule)
	if err != nil {
		return nil, err
	}
	return &ExchangecontractCancelIterator{contract: _Exchangecontract.contract, event: "Cancel", logs: logs, sub: sub}, nil
}

// WatchCancel is a free log subscription operation binding the contract event 0x02c310a9a43963ff31a754a4099cc435ed498049687539d72d7818d9b093415c.
//
// Solidity: event Cancel(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, address senderAddress, bytes32 indexed orderHash)
func (_Exchangecontract *ExchangecontractFilterer) WatchCancel(opts *bind.WatchOpts, sink chan<- *ExchangecontractCancel, makerAddress []common.Address, feeRecipientAddress []common.Address, orderHash [][32]byte) (event.Subscription, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var feeRecipientAddressRule []interface{}
	for _, feeRecipientAddressItem := range feeRecipientAddress {
		feeRecipientAddressRule = append(feeRecipientAddressRule, feeRecipientAddressItem)
	}

	var orderHashRule []interface{}
	for _, orderHashItem := range orderHash {
		orderHashRule = append(orderHashRule, orderHashItem)
	}

	logs, sub, err := _Exchangecontract.contract.WatchLogs(opts, "Cancel", makerAddressRule, feeRecipientAddressRule, orderHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangecontractCancel)
				if err := _Exchangecontract.contract.UnpackLog(event, "Cancel", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCancel is a log parse operation binding the contract event 0x02c310a9a43963ff31a754a4099cc435ed498049687539d72d7818d9b093415c.
//
// Solidity: event Cancel(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, address senderAddress, bytes32 indexed orderHash)
func (_Exchangecontract *ExchangecontractFilterer) ParseCancel(log types.Log) (*ExchangecontractCancel, error) {
	event := new(ExchangecontractCancel)
	if err := _Exchangecontract.contract.UnpackLog(event, "Cancel", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangecontractCancelUpToIterator is returned from FilterCancelUpTo and is used to iterate over the raw logs and unpacked data for CancelUpTo events raised by the Exchangecontract contract.
type ExchangecontractCancelUpToIterator struct {
	Event *ExchangecontractCancelUpTo // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangecontractCancelUpToIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangecontractCancelUpTo)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangecontractCancelUpTo)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangecontractCancelUpToIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangecontractCancelUpToIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangecontractCancelUpTo represents a CancelUpTo event raised by the Exchangecontract contract.
type ExchangecontractCancelUpTo struct {
	MakerAddress       common.Address
	OrderSenderAddress common.Address
	OrderEpoch         *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterCancelUpTo is a free log retrieval operation binding the contract event 0x82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0.
//
// Solidity: event CancelUpTo(address indexed makerAddress, address indexed orderSenderAddress, uint256 orderEpoch)
func (_Exchangecontract *ExchangecontractFilterer) FilterCancelUpTo(opts *bind.FilterOpts, makerAddress []common.Address, orderSenderAddress []common.Address) (*ExchangecontractCancelUpToIterator, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var orderSenderAddressRule []interface{}
	for _, orderSenderAddressItem := range orderSenderAddress {
		orderSenderAddressRule = append(orderSenderAddressRule, orderSenderAddressItem)
	}

	logs, sub, err := _Exchangecontract.contract.FilterLogs(opts, "CancelUpTo", makerAddressRule, orderSenderAddressRule)
	if err != nil {
		return nil, err
	}
	return &ExchangecontractCancelUpToIterator{contract: _Exchangecontract.contract, event: "CancelUpTo", logs: logs, sub: sub}, nil
}

// WatchCancelUpTo is a free log subscription operation binding the contract event 0x82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0.
//
// Solidity: event CancelUpTo(address indexed makerAddress, address indexed orderSenderAddress, uint256 orderEpoch)
func (_Exchangecontract *ExchangecontractFilterer) WatchCancelUpTo(opts *bind.WatchOpts, sink chan<- *ExchangecontractCancelUpTo, makerAddress []common.Address, orderSenderAddress []common.Address) (event.Subscription, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var orderSenderAddressRule []interface{}
	for _, orderSenderAddressItem := range orderSenderAddress {
		orderSenderAddressRule = append(orderSenderAddressRule, orderSenderAddressItem)
	}

	logs, sub, err := _Exchangecontract.contract.WatchLogs(opts, "CancelUpTo", makerAddressRule, orderSenderAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangecontractCancelUpTo)
				if err := _Exchangecontract.contract.UnpackLog(event, "CancelUpTo", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCancelUpTo is a log parse operation binding the contract event 0x82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0.
//
// Solidity: event CancelUpTo(address indexed makerAddress, address indexed orderSenderAddress, uint256 orderEpoch)
func (_Exchangecontract *ExchangecontractFilterer) ParseCancelUpTo(log types.Log) (*ExchangecontractCancelUpTo, error) {
	event := new(ExchangecontractCancelUpTo)
	if err := _Exchangecontract.contract.UnpackLog(event, "CancelUpTo", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangecontractFillIterator is returned from FilterFill and is used to iterate over the raw logs and unpacked data for Fill events raised by the Exchangecontract contract.
type ExchangecontractFillIterator struct {
	Event *ExchangecontractFill // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangecontractFillIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangecontractFill)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangecontractFill)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangecontractFillIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangecontractFillIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangecontractFill represents a Fill event raised by the Exchangecontract contract.
type ExchangecontractFill struct {
	MakerAddress           common.Address
	FeeRecipientAddress    common.Address
	MakerAssetData         []byte
	TakerAssetData         []byte
	MakerFeeAssetData      []byte
	TakerFeeAssetData      []byte
	OrderHash              [32]byte
	TakerAddress           common.Address
	SenderAddress          common.Address
	MakerAssetFilledAmount *big.Int
	TakerAssetFilledAmount *big.Int
	MakerFeePaid           *big.Int
	TakerFeePaid           *big.Int
	ProtocolFeePaid        *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterFill is a free log retrieval operation binding the contract event 0x6869791f0a34781b29882982cc39e882768cf2c96995c2a110c577c53bc932d5.
//
// Solidity: event Fill(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, bytes makerFeeAssetData, bytes takerFeeAssetData, bytes32 indexed orderHash, address takerAddress, address senderAddress, uint256 makerAssetFilledAmount, uint256 takerAssetFilledAmount, uint256 makerFeePaid, uint256 takerFeePaid, uint256 protocolFeePaid)
func (_Exchangecontract *ExchangecontractFilterer) FilterFill(opts *bind.FilterOpts, makerAddress []common.Address, feeRecipientAddress []common.Address, orderHash [][32]byte) (*ExchangecontractFillIterator, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var feeRecipientAddressRule []interface{}
	for _, feeRecipientAddressItem := range feeRecipientAddress {
		feeRecipientAddressRule = append(feeRecipientAddressRule, feeRecipientAddressItem)
	}

	var orderHashRule []interface{}
	for _, orderHashItem := range orderHash {
		orderHashRule = append(orderHashRule, orderHashItem)
	}

	logs, sub, err := _Exchangecontract.contract.FilterLogs(opts, "Fill", makerAddressRule, feeRecipientAddressRule, orderHashRule)
	if err != nil {
		return nil, err
	}
	return &ExchangecontractFillIterator{contract: _Exchangecontract.contract, event: "Fill", logs: logs, sub: sub}, nil
}

// WatchFill is a free log subscription operation binding the contract event 0x6869791f0a34781b29882982cc39e882768cf2c96995c2a110c577c53bc932d5.
//
// Solidity: event Fill(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, bytes makerFeeAssetData, bytes takerFeeAssetData, bytes32 indexed orderHash, address takerAddress, address senderAddress, uint256 makerAssetFilledAmount, uint256 takerAssetFilledAmount, uint256 makerFeePaid, uint256 takerFeePaid, uint256 protocolFeePaid)
func (_Exchangecontract *ExchangecontractFilterer) WatchFill(opts *bind.WatchOpts, sink chan<- *ExchangecontractFill, makerAddress []common.Address, feeRecipientAddress []common.Address, orderHash [][32]byte) (event.Subscription, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var feeRecipientAddressRule []interface{}
	for _, feeRecipientAddressItem := range feeRecipientAddress {
		feeRecipientAddressRule = append(feeRecipientAddressRule, feeRecipientAddressItem)
	}

	var orderHashRule []interface{}
	for _, orderHashItem := range orderHash {
		orderHashRule = append(orderHashRule, orderHashItem)
	}

	logs, sub, err := _Exchangecontract.contract.WatchLogs(opts, "Fill", makerAddressRule, feeRecipientAddressRule, orderHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangecontractFill)
				if err := _Exchangecontract.contract.UnpackLog(event, "Fill", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFill is a log parse operation binding the contract event 0x6869791f0a34781b29882982cc39e882768cf2c96995c2a110c577c53bc932d5.
//
// Solidity: event Fill(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, bytes makerFeeAssetData, bytes takerFeeAssetData, bytes32 indexed orderHash, address takerAddress, address senderAddress, uint256 makerAssetFilledAmount, uint256 takerAssetFilledAmount, uint256 makerFeePaid, uint256 takerFeePaid, uint256 protocolFeePaid)
func (_Exchangecontract *ExchangecontractFilterer) ParseFill(log types.Log) (*ExchangecontractFill, error) {
	event := new(ExchangecontractFill)
	if err := _Exchangecontract.contract.UnpackLog(event, "Fill", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangecontractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Exchangecontract contract.
type ExchangecontractOwnershipTransferredIterator struct {
	Event *ExchangecontractOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangecontractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangecontractOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangecontractOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangecontractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangecontractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangecontractOwnershipTransferred represents a OwnershipTransferred event raised by the Exchangecontract contract.
type ExchangecontractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Exchangecontract *ExchangecontractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ExchangecontractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Exchangecontract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ExchangecontractOwnershipTransferredIterator{contract: _Exchangecontract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Exchangecontract *ExchangecontractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ExchangecontractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Exchangecontract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangecontractOwnershipTransferred)
				if err := _Exchangecontract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Exchangecontract *ExchangecontractFilterer) ParseOwnershipTransferred(log types.Log) (*ExchangecontractOwnershipTransferred, error) {
	event := new(ExchangecontractOwnershipTransferred)
	if err := _Exchangecontract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangecontractProtocolFeeCollectorAddressIterator is returned from FilterProtocolFeeCollectorAddress and is used to iterate over the raw logs and unpacked data for ProtocolFeeCollectorAddress events raised by the Exchangecontract contract.
type ExchangecontractProtocolFeeCollectorAddressIterator struct {
	Event *ExchangecontractProtocolFeeCollectorAddress // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangecontractProtocolFeeCollectorAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangecontractProtocolFeeCollectorAddress)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangecontractProtocolFeeCollectorAddress)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangecontractProtocolFeeCollectorAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangecontractProtocolFeeCollectorAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangecontractProtocolFeeCollectorAddress represents a ProtocolFeeCollectorAddress event raised by the Exchangecontract contract.
type ExchangecontractProtocolFeeCollectorAddress struct {
	OldProtocolFeeCollector     common.Address
	UpdatedProtocolFeeCollector common.Address
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeeCollectorAddress is a free log retrieval operation binding the contract event 0xe1a5430ebec577336427f40f15822f1f36c5e3509ff209d6db9e6c9e6941cb0b.
//
// Solidity: event ProtocolFeeCollectorAddress(address oldProtocolFeeCollector, address updatedProtocolFeeCollector)
func (_Exchangecontract *ExchangecontractFilterer) FilterProtocolFeeCollectorAddress(opts *bind.FilterOpts) (*ExchangecontractProtocolFeeCollectorAddressIterator, error) {

	logs, sub, err := _Exchangecontract.contract.FilterLogs(opts, "ProtocolFeeCollectorAddress")
	if err != nil {
		return nil, err
	}
	return &ExchangecontractProtocolFeeCollectorAddressIterator{contract: _Exchangecontract.contract, event: "ProtocolFeeCollectorAddress", logs: logs, sub: sub}, nil
}

// WatchProtocolFeeCollectorAddress is a free log subscription operation binding the contract event 0xe1a5430ebec577336427f40f15822f1f36c5e3509ff209d6db9e6c9e6941cb0b.
//
// Solidity: event ProtocolFeeCollectorAddress(address oldProtocolFeeCollector, address updatedProtocolFeeCollector)
func (_Exchangecontract *ExchangecontractFilterer) WatchProtocolFeeCollectorAddress(opts *bind.WatchOpts, sink chan<- *ExchangecontractProtocolFeeCollectorAddress) (event.Subscription, error) {

	logs, sub, err := _Exchangecontract.contract.WatchLogs(opts, "ProtocolFeeCollectorAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangecontractProtocolFeeCollectorAddress)
				if err := _Exchangecontract.contract.UnpackLog(event, "ProtocolFeeCollectorAddress", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseProtocolFeeCollectorAddress is a log parse operation binding the contract event 0xe1a5430ebec577336427f40f15822f1f36c5e3509ff209d6db9e6c9e6941cb0b.
//
// Solidity: event ProtocolFeeCollectorAddress(address oldProtocolFeeCollector, address updatedProtocolFeeCollector)
func (_Exchangecontract *ExchangecontractFilterer) ParseProtocolFeeCollectorAddress(log types.Log) (*ExchangecontractProtocolFeeCollectorAddress, error) {
	event := new(ExchangecontractProtocolFeeCollectorAddress)
	if err := _Exchangecontract.contract.UnpackLog(event, "ProtocolFeeCollectorAddress", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangecontractProtocolFeeMultiplierIterator is returned from FilterProtocolFeeMultiplier and is used to iterate over the raw logs and unpacked data for ProtocolFeeMultiplier events raised by the Exchangecontract contract.
type ExchangecontractProtocolFeeMultiplierIterator struct {
	Event *ExchangecontractProtocolFeeMultiplier // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangecontractProtocolFeeMultiplierIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangecontractProtocolFeeMultiplier)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangecontractProtocolFeeMultiplier)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangecontractProtocolFeeMultiplierIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangecontractProtocolFeeMultiplierIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangecontractProtocolFeeMultiplier represents a ProtocolFeeMultiplier event raised by the Exchangecontract contract.
type ExchangecontractProtocolFeeMultiplier struct {
	OldProtocolFeeMultiplier     *big.Int
	UpdatedProtocolFeeMultiplier *big.Int
	Raw                          types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeeMultiplier is a free log retrieval operation binding the contract event 0x3a3e76d7a75e198aef1f53137e4f2a8a2ec74e2e9526db8404d08ccc9f1e621d.
//
// Solidity: event ProtocolFeeMultiplier(uint256 oldProtocolFeeMultiplier, uint256 updatedProtocolFeeMultiplier)
func (_Exchangecontract *ExchangecontractFilterer) FilterProtocolFeeMultiplier(opts *bind.FilterOpts) (*ExchangecontractProtocolFeeMultiplierIterator, error) {

	logs, sub, err := _Exchangecontract.contract.FilterLogs(opts, "ProtocolFeeMultiplier")
	if err != nil {
		return nil, err
	}
	return &ExchangecontractProtocolFeeMultiplierIterator{contract: _Exchangecontract.contract, event: "ProtocolFeeMultiplier", logs: logs, sub: sub}, nil
}

// WatchProtocolFeeMultiplier is a free log subscription operation binding the contract event 0x3a3e76d7a75e198aef1f53137e4f2a8a2ec74e2e9526db8404d08ccc9f1e621d.
//
// Solidity: event ProtocolFeeMultiplier(uint256 oldProtocolFeeMultiplier, uint256 updatedProtocolFeeMultiplier)
func (_Exchangecontract *ExchangecontractFilterer) WatchProtocolFeeMultiplier(opts *bind.WatchOpts, sink chan<- *ExchangecontractProtocolFeeMultiplier) (event.Subscription, error) {

	logs, sub, err := _Exchangecontract.contract.WatchLogs(opts, "ProtocolFeeMultiplier")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangecontractProtocolFeeMultiplier)
				if err := _Exchangecontract.contract.UnpackLog(event, "ProtocolFeeMultiplier", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseProtocolFeeMultiplier is a log parse operation binding the contract event 0x3a3e76d7a75e198aef1f53137e4f2a8a2ec74e2e9526db8404d08ccc9f1e621d.
//
// Solidity: event ProtocolFeeMultiplier(uint256 oldProtocolFeeMultiplier, uint256 updatedProtocolFeeMultiplier)
func (_Exchangecontract *ExchangecontractFilterer) ParseProtocolFeeMultiplier(log types.Log) (*ExchangecontractProtocolFeeMultiplier, error) {
	event := new(ExchangecontractProtocolFeeMultiplier)
	if err := _Exchangecontract.contract.UnpackLog(event, "ProtocolFeeMultiplier", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangecontractSignatureValidatorApprovalIterator is returned from FilterSignatureValidatorApproval and is used to iterate over the raw logs and unpacked data for SignatureValidatorApproval events raised by the Exchangecontract contract.
type ExchangecontractSignatureValidatorApprovalIterator struct {
	Event *ExchangecontractSignatureValidatorApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangecontractSignatureValidatorApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangecontractSignatureValidatorApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangecontractSignatureValidatorApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangecontractSignatureValidatorApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangecontractSignatureValidatorApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangecontractSignatureValidatorApproval represents a SignatureValidatorApproval event raised by the Exchangecontract contract.
type ExchangecontractSignatureValidatorApproval struct {
	SignerAddress    common.Address
	ValidatorAddress common.Address
	IsApproved       bool
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSignatureValidatorApproval is a free log retrieval operation binding the contract event 0xa8656e308026eeabce8f0bc18048433252318ab80ac79da0b3d3d8697dfba891.
//
// Solidity: event SignatureValidatorApproval(address indexed signerAddress, address indexed validatorAddress, bool isApproved)
func (_Exchangecontract *ExchangecontractFilterer) FilterSignatureValidatorApproval(opts *bind.FilterOpts, signerAddress []common.Address, validatorAddress []common.Address) (*ExchangecontractSignatureValidatorApprovalIterator, error) {

	var signerAddressRule []interface{}
	for _, signerAddressItem := range signerAddress {
		signerAddressRule = append(signerAddressRule, signerAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _Exchangecontract.contract.FilterLogs(opts, "SignatureValidatorApproval", signerAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return &ExchangecontractSignatureValidatorApprovalIterator{contract: _Exchangecontract.contract, event: "SignatureValidatorApproval", logs: logs, sub: sub}, nil
}

// WatchSignatureValidatorApproval is a free log subscription operation binding the contract event 0xa8656e308026eeabce8f0bc18048433252318ab80ac79da0b3d3d8697dfba891.
//
// Solidity: event SignatureValidatorApproval(address indexed signerAddress, address indexed validatorAddress, bool isApproved)
func (_Exchangecontract *ExchangecontractFilterer) WatchSignatureValidatorApproval(opts *bind.WatchOpts, sink chan<- *ExchangecontractSignatureValidatorApproval, signerAddress []common.Address, validatorAddress []common.Address) (event.Subscription, error) {

	var signerAddressRule []interface{}
	for _, signerAddressItem := range signerAddress {
		signerAddressRule = append(signerAddressRule, signerAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _Exchangecontract.contract.WatchLogs(opts, "SignatureValidatorApproval", signerAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangecontractSignatureValidatorApproval)
				if err := _Exchangecontract.contract.UnpackLog(event, "SignatureValidatorApproval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSignatureValidatorApproval is a log parse operation binding the contract event 0xa8656e308026eeabce8f0bc18048433252318ab80ac79da0b3d3d8697dfba891.
//
// Solidity: event SignatureValidatorApproval(address indexed signerAddress, address indexed validatorAddress, bool isApproved)
func (_Exchangecontract *ExchangecontractFilterer) ParseSignatureValidatorApproval(log types.Log) (*ExchangecontractSignatureValidatorApproval, error) {
	event := new(ExchangecontractSignatureValidatorApproval)
	if err := _Exchangecontract.contract.UnpackLog(event, "SignatureValidatorApproval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangecontractTransactionExecutionIterator is returned from FilterTransactionExecution and is used to iterate over the raw logs and unpacked data for TransactionExecution events raised by the Exchangecontract contract.
type ExchangecontractTransactionExecutionIterator struct {
	Event *ExchangecontractTransactionExecution // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangecontractTransactionExecutionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangecontractTransactionExecution)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangecontractTransactionExecution)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangecontractTransactionExecutionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangecontractTransactionExecutionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangecontractTransactionExecution represents a TransactionExecution event raised by the Exchangecontract contract.
type ExchangecontractTransactionExecution struct {
	TransactionHash [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterTransactionExecution is a free log retrieval operation binding the contract event 0xa4a7329f1dd821363067e07d359e347b4af9b1efe4b6cccf13240228af3c800d.
//
// Solidity: event TransactionExecution(bytes32 indexed transactionHash)
func (_Exchangecontract *ExchangecontractFilterer) FilterTransactionExecution(opts *bind.FilterOpts, transactionHash [][32]byte) (*ExchangecontractTransactionExecutionIterator, error) {

	var transactionHashRule []interface{}
	for _, transactionHashItem := range transactionHash {
		transactionHashRule = append(transactionHashRule, transactionHashItem)
	}

	logs, sub, err := _Exchangecontract.contract.FilterLogs(opts, "TransactionExecution", transactionHashRule)
	if err != nil {
		return nil, err
	}
	return &ExchangecontractTransactionExecutionIterator{contract: _Exchangecontract.contract, event: "TransactionExecution", logs: logs, sub: sub}, nil
}

// WatchTransactionExecution is a free log subscription operation binding the contract event 0xa4a7329f1dd821363067e07d359e347b4af9b1efe4b6cccf13240228af3c800d.
//
// Solidity: event TransactionExecution(bytes32 indexed transactionHash)
func (_Exchangecontract *ExchangecontractFilterer) WatchTransactionExecution(opts *bind.WatchOpts, sink chan<- *ExchangecontractTransactionExecution, transactionHash [][32]byte) (event.Subscription, error) {

	var transactionHashRule []interface{}
	for _, transactionHashItem := range transactionHash {
		transactionHashRule = append(transactionHashRule, transactionHashItem)
	}

	logs, sub, err := _Exchangecontract.contract.WatchLogs(opts, "TransactionExecution", transactionHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangecontractTransactionExecution)
				if err := _Exchangecontract.contract.UnpackLog(event, "TransactionExecution", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransactionExecution is a log parse operation binding the contract event 0xa4a7329f1dd821363067e07d359e347b4af9b1efe4b6cccf13240228af3c800d.
//
// Solidity: event TransactionExecution(bytes32 indexed transactionHash)
func (_Exchangecontract *ExchangecontractFilterer) ParseTransactionExecution(log types.Log) (*ExchangecontractTransactionExecution, error) {
	event := new(ExchangecontractTransactionExecution)
	if err := _Exchangecontract.contract.UnpackLog(event, "TransactionExecution", log); err != nil {
		return nil, err
	}
	return event, nil
}
