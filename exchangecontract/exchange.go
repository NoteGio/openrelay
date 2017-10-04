// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package exchangecontract

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ExchangeABI is the input ABI used to generate the binding from.
const ExchangeABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"numerator\",\"type\":\"uint256\"},{\"name\":\"denominator\",\"type\":\"uint256\"},{\"name\":\"target\",\"type\":\"uint256\"}],\"name\":\"isRoundingError\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"filled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"cancelled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderAddresses\",\"type\":\"address[5][]\"},{\"name\":\"orderValues\",\"type\":\"uint256[6][]\"},{\"name\":\"fillTakerTokenAmount\",\"type\":\"uint256\"},{\"name\":\"shouldThrowOnInsufficientBalanceOrAllowance\",\"type\":\"bool\"},{\"name\":\"v\",\"type\":\"uint8[]\"},{\"name\":\"r\",\"type\":\"bytes32[]\"},{\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"fillOrdersUpTo\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderAddresses\",\"type\":\"address[5]\"},{\"name\":\"orderValues\",\"type\":\"uint256[6]\"},{\"name\":\"cancelTakerTokenAmount\",\"type\":\"uint256\"}],\"name\":\"cancelOrder\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ZRX_TOKEN_CONTRACT\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderAddresses\",\"type\":\"address[5][]\"},{\"name\":\"orderValues\",\"type\":\"uint256[6][]\"},{\"name\":\"fillTakerTokenAmounts\",\"type\":\"uint256[]\"},{\"name\":\"v\",\"type\":\"uint8[]\"},{\"name\":\"r\",\"type\":\"bytes32[]\"},{\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"batchFillOrKillOrders\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderAddresses\",\"type\":\"address[5]\"},{\"name\":\"orderValues\",\"type\":\"uint256[6]\"},{\"name\":\"fillTakerTokenAmount\",\"type\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"fillOrKillOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"getUnavailableTakerTokenAmount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"},{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"isValidSignature\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"numerator\",\"type\":\"uint256\"},{\"name\":\"denominator\",\"type\":\"uint256\"},{\"name\":\"target\",\"type\":\"uint256\"}],\"name\":\"getPartialAmount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TOKEN_TRANSFER_PROXY_CONTRACT\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderAddresses\",\"type\":\"address[5][]\"},{\"name\":\"orderValues\",\"type\":\"uint256[6][]\"},{\"name\":\"fillTakerTokenAmounts\",\"type\":\"uint256[]\"},{\"name\":\"shouldThrowOnInsufficientBalanceOrAllowance\",\"type\":\"bool\"},{\"name\":\"v\",\"type\":\"uint8[]\"},{\"name\":\"r\",\"type\":\"bytes32[]\"},{\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"batchFillOrders\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderAddresses\",\"type\":\"address[5][]\"},{\"name\":\"orderValues\",\"type\":\"uint256[6][]\"},{\"name\":\"cancelTakerTokenAmounts\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelOrders\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderAddresses\",\"type\":\"address[5]\"},{\"name\":\"orderValues\",\"type\":\"uint256[6]\"},{\"name\":\"fillTakerTokenAmount\",\"type\":\"uint256\"},{\"name\":\"shouldThrowOnInsufficientBalanceOrAllowance\",\"type\":\"bool\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"fillOrder\",\"outputs\":[{\"name\":\"filledTakerTokenAmount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderAddresses\",\"type\":\"address[5]\"},{\"name\":\"orderValues\",\"type\":\"uint256[6]\"}],\"name\":\"getOrderHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"EXTERNAL_QUERY_GAS_LIMIT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_zrxToken\",\"type\":\"address\"},{\"name\":\"_tokenTransferProxy\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipient\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"makerToken\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"takerToken\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"filledMakerTokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"filledTakerTokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"paidMakerFee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"paidTakerFee\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"tokens\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"LogFill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipient\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"makerToken\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"takerToken\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"cancelledMakerTokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"cancelledTakerTokenAmount\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"tokens\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"LogCancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"errorId\",\"type\":\"uint8\"},{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"LogError\",\"type\":\"event\"}]"

// ExchangeBin is the compiled bytecode used for deploying new contracts.
const ExchangeBin = `0x6060604052341561000f57600080fd5b604051604080610db883398101604052808051919060200180519150505b5b50505b610d78806100406000396000f300606060405236156100f95763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166314df96ee81146100fe578063288cdc911461012e5780632ac1262214610156578063363349be1461017e578063394c21e7146103555780633b30ba59146103ca5780634f15078714610406578063741bcc93146105f95780637e9abb50146106715780638163681e1461069957806398024a8b146106e8578063add1cbc514610716578063b7b2c7d614610752578063baa0181d14610950578063bc61394a14610a83578063cfc4d0ec14610b12578063f06bbf7514610b85578063ffa1ad7414610baf575b600080fd5b341561010957600080fd5b61011a600435602435604435610c3a565b604051901515815260200160405180910390f35b341561013957600080fd5b610144600435610c44565b60405190815260200160405180910390f35b341561016157600080fd5b610144600435610c56565b60405190815260200160405180910390f35b341561018957600080fd5b61014460046024813581810190830135806020818102016040519081016040528181529291906000602085015b828210156101f25760a08083028601906005906040519081016040529190828260a08082843750505091835250506001909101906020016101b6565b505050505091908035906020019082018035906020019080806020026020016040519081016040528181529291906000602085015b828210156102635760c08083028601906006906040519081016040529190828260c0808284375050509183525050600190910190602001610227565b5050505050919080359060200190919080351515906020019091908035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919080359060200190820180359060200190808060200260200160405190810160405280939291908181526020018383602002808284378201915050505050509190803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843750949650610c6895505050505050565b60405190815260200160405180910390f35b341561036057600080fd5b610144600460a481600560a06040519081016040529190828260a080828437820191505050505091908060c001906006806020026040519081016040529190828260c080828437509395505092359250610c3a915050565b60405190815260200160405180910390f35b34156103d557600080fd5b6103dd610c80565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b341561041157600080fd5b6105f760046024813581810190830135806020818102016040519081016040528181529291906000602085015b8282101561047a5760a08083028601906005906040519081016040529190828260a080828437505050918352505060019091019060200161043e565b505050505091908035906020019082018035906020019080806020026020016040519081016040528181529291906000602085015b828210156104eb5760c08083028601906006906040519081016040529190828260c08082843750505091835250506001909101906020016104af565b50505050509190803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091908035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919080359060200190820180359060200190808060200260200160405190810160405280939291908181526020018383602002808284378201915050505050509190803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843750949650610c9c95505050505050565b005b341561060457600080fd5b6105f7600460a481600560a06040519081016040529190828260a080828437820191505050505091908060c001906006806020026040519081016040529190828260c080828437509395505083359360ff602082013516935060408101359250606001359050610c9c565b005b341561067c57600080fd5b610144600435610cae565b60405190815260200160405180910390f35b34156106a457600080fd5b61011a73ffffffffffffffffffffffffffffffffffffffff6004351660243560ff60443516606435608435610cb6565b604051901515815260200160405180910390f35b34156106f357600080fd5b610144600435602435604435610c3a565b60405190815260200160405180910390f35b341561072157600080fd5b6103dd610ccc565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b341561075d57600080fd5b6105f760046024813581810190830135806020818102016040519081016040528181529291906000602085015b828210156107c65760a08083028601906005906040519081016040529190828260a080828437505050918352505060019091019060200161078a565b505050505091908035906020019082018035906020019080806020026020016040519081016040528181529291906000602085015b828210156108375760c08083028601906006906040519081016040529190828260c08082843750505091835250506001909101906020016107fb565b505050505091908035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919080351515906020019091908035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919080359060200190820180359060200190808060200260200160405190810160405280939291908181526020018383602002808284378201915050505050509190803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843750949650610ce895505050505050565b005b341561095b57600080fd5b6105f760046024813581810190830135806020818102016040519081016040528181529291906000602085015b828210156109c45760a08083028601906005906040519081016040529190828260a0808284375050509183525050600190910190602001610988565b505050505091908035906020019082018035906020019080806020026020016040519081016040528181529291906000602085015b82821015610a355760c08083028601906006906040519081016040529190828260c08082843750505091835250506001909101906020016109f9565b50505050509190803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843750949650610cf295505050505050565b005b3415610a8e57600080fd5b610144600460a481600560a06040519081016040529190828260a080828437820191505050505091908060c001906006806020026040519081016040529190828260c080828437509395505083359360208101351515935060ff60408201351692506060810135915060800135610c68565b60405190815260200160405180910390f35b3415610b1d57600080fd5b610144600460a481600560a06040519081016040529190828260a080828437820191505050505091908060c001906006806020026040519081016040529190828260c08082843750939550610d06945050505050565b60405190815260200160405180910390f35b3415610b9057600080fd5b610b98610d0f565b60405161ffff909116815260200160405180910390f35b3415610bba57600080fd5b610bc2610d15565b60405160208082528190810183818151815260200191508051906020019080838360005b83811015610bff5780820151818401525b602001610be6565b50505050905090810190601f168015610c2c5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60005b9392505050565b60026020526000908152604090205481565b60036020526000908152604090205481565b60005b979650505050505050565b60005b9392505050565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b5b505050505050565b5b505050505050565b60005b919050565b60005b95945050505050565b60005b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff1681565b5b50505050505050565b5b505050565b60005b979650505050505050565b60005b92915050565b61138781565b60408051908101604052600581527f312e302e300000000000000000000000000000000000000000000000000000006020820152815600a165627a7a72305820cefd703378f89befd84d037094e3e15b8a1c369307c3e72860c9191cde98a9f90029`

// DeployExchange deploys a new Ethereum contract, binding an instance of Exchange to it.
func DeployExchange(auth *bind.TransactOpts, backend bind.ContractBackend, _zrxToken common.Address, _tokenTransferProxy common.Address) (common.Address, *types.Transaction, *Exchange, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ExchangeBin), backend, _zrxToken, _tokenTransferProxy)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Exchange{ExchangeCaller: ExchangeCaller{contract: contract}, ExchangeTransactor: ExchangeTransactor{contract: contract}}, nil
}

// Exchange is an auto generated Go binding around an Ethereum contract.
type Exchange struct {
	ExchangeCaller     // Read-only binding to the contract
	ExchangeTransactor // Write-only binding to the contract
}

// ExchangeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExchangeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExchangeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExchangeSession struct {
	Contract     *Exchange         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ExchangeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExchangeCallerSession struct {
	Contract *ExchangeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ExchangeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExchangeTransactorSession struct {
	Contract     *ExchangeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ExchangeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExchangeRaw struct {
	Contract *Exchange // Generic contract binding to access the raw methods on
}

// ExchangeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExchangeCallerRaw struct {
	Contract *ExchangeCaller // Generic read-only contract binding to access the raw methods on
}

// ExchangeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExchangeTransactorRaw struct {
	Contract *ExchangeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExchange creates a new instance of Exchange, bound to a specific deployed contract.
func NewExchange(address common.Address, backend bind.ContractBackend) (*Exchange, error) {
	contract, err := bindExchange(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Exchange{ExchangeCaller: ExchangeCaller{contract: contract}, ExchangeTransactor: ExchangeTransactor{contract: contract}}, nil
}

// NewExchangeCaller creates a new read-only instance of Exchange, bound to a specific deployed contract.
func NewExchangeCaller(address common.Address, caller bind.ContractCaller) (*ExchangeCaller, error) {
	contract, err := bindExchange(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeCaller{contract: contract}, nil
}

// NewExchangeTransactor creates a new write-only instance of Exchange, bound to a specific deployed contract.
func NewExchangeTransactor(address common.Address, transactor bind.ContractTransactor) (*ExchangeTransactor, error) {
	contract, err := bindExchange(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &ExchangeTransactor{contract: contract}, nil
}

// bindExchange binds a generic wrapper to an already deployed contract.
func bindExchange(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Exchange *ExchangeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Exchange.Contract.ExchangeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Exchange *ExchangeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchange.Contract.ExchangeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Exchange *ExchangeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Exchange.Contract.ExchangeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Exchange *ExchangeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Exchange.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Exchange *ExchangeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchange.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Exchange *ExchangeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Exchange.Contract.contract.Transact(opts, method, params...)
}

// EXTERNAL_QUERY_GAS_LIMIT is a free data retrieval call binding the contract method 0xf06bbf75.
//
// Solidity: function EXTERNAL_QUERY_GAS_LIMIT() constant returns(uint16)
func (_Exchange *ExchangeCaller) EXTERNAL_QUERY_GAS_LIMIT(opts *bind.CallOpts) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "EXTERNAL_QUERY_GAS_LIMIT")
	return *ret0, err
}

// EXTERNAL_QUERY_GAS_LIMIT is a free data retrieval call binding the contract method 0xf06bbf75.
//
// Solidity: function EXTERNAL_QUERY_GAS_LIMIT() constant returns(uint16)
func (_Exchange *ExchangeSession) EXTERNAL_QUERY_GAS_LIMIT() (uint16, error) {
	return _Exchange.Contract.EXTERNAL_QUERY_GAS_LIMIT(&_Exchange.CallOpts)
}

// EXTERNAL_QUERY_GAS_LIMIT is a free data retrieval call binding the contract method 0xf06bbf75.
//
// Solidity: function EXTERNAL_QUERY_GAS_LIMIT() constant returns(uint16)
func (_Exchange *ExchangeCallerSession) EXTERNAL_QUERY_GAS_LIMIT() (uint16, error) {
	return _Exchange.Contract.EXTERNAL_QUERY_GAS_LIMIT(&_Exchange.CallOpts)
}

// TOKEN_TRANSFER_PROXY_CONTRACT is a free data retrieval call binding the contract method 0xadd1cbc5.
//
// Solidity: function TOKEN_TRANSFER_PROXY_CONTRACT() constant returns(address)
func (_Exchange *ExchangeCaller) TOKEN_TRANSFER_PROXY_CONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "TOKEN_TRANSFER_PROXY_CONTRACT")
	return *ret0, err
}

// TOKEN_TRANSFER_PROXY_CONTRACT is a free data retrieval call binding the contract method 0xadd1cbc5.
//
// Solidity: function TOKEN_TRANSFER_PROXY_CONTRACT() constant returns(address)
func (_Exchange *ExchangeSession) TOKEN_TRANSFER_PROXY_CONTRACT() (common.Address, error) {
	return _Exchange.Contract.TOKEN_TRANSFER_PROXY_CONTRACT(&_Exchange.CallOpts)
}

// TOKEN_TRANSFER_PROXY_CONTRACT is a free data retrieval call binding the contract method 0xadd1cbc5.
//
// Solidity: function TOKEN_TRANSFER_PROXY_CONTRACT() constant returns(address)
func (_Exchange *ExchangeCallerSession) TOKEN_TRANSFER_PROXY_CONTRACT() (common.Address, error) {
	return _Exchange.Contract.TOKEN_TRANSFER_PROXY_CONTRACT(&_Exchange.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Exchange *ExchangeCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Exchange *ExchangeSession) VERSION() (string, error) {
	return _Exchange.Contract.VERSION(&_Exchange.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Exchange *ExchangeCallerSession) VERSION() (string, error) {
	return _Exchange.Contract.VERSION(&_Exchange.CallOpts)
}

// ZRX_TOKEN_CONTRACT is a free data retrieval call binding the contract method 0x3b30ba59.
//
// Solidity: function ZRX_TOKEN_CONTRACT() constant returns(address)
func (_Exchange *ExchangeCaller) ZRX_TOKEN_CONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "ZRX_TOKEN_CONTRACT")
	return *ret0, err
}

// ZRX_TOKEN_CONTRACT is a free data retrieval call binding the contract method 0x3b30ba59.
//
// Solidity: function ZRX_TOKEN_CONTRACT() constant returns(address)
func (_Exchange *ExchangeSession) ZRX_TOKEN_CONTRACT() (common.Address, error) {
	return _Exchange.Contract.ZRX_TOKEN_CONTRACT(&_Exchange.CallOpts)
}

// ZRX_TOKEN_CONTRACT is a free data retrieval call binding the contract method 0x3b30ba59.
//
// Solidity: function ZRX_TOKEN_CONTRACT() constant returns(address)
func (_Exchange *ExchangeCallerSession) ZRX_TOKEN_CONTRACT() (common.Address, error) {
	return _Exchange.Contract.ZRX_TOKEN_CONTRACT(&_Exchange.CallOpts)
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled( bytes32) constant returns(uint256)
func (_Exchange *ExchangeCaller) Cancelled(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "cancelled", arg0)
	return *ret0, err
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled( bytes32) constant returns(uint256)
func (_Exchange *ExchangeSession) Cancelled(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.Cancelled(&_Exchange.CallOpts, arg0)
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled( bytes32) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) Cancelled(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.Cancelled(&_Exchange.CallOpts, arg0)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled( bytes32) constant returns(uint256)
func (_Exchange *ExchangeCaller) Filled(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "filled", arg0)
	return *ret0, err
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled( bytes32) constant returns(uint256)
func (_Exchange *ExchangeSession) Filled(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.Filled(&_Exchange.CallOpts, arg0)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled( bytes32) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) Filled(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.Filled(&_Exchange.CallOpts, arg0)
}

// GetOrderHash is a free data retrieval call binding the contract method 0xcfc4d0ec.
//
// Solidity: function getOrderHash(orderAddresses address[5], orderValues uint256[6]) constant returns(bytes32)
func (_Exchange *ExchangeCaller) GetOrderHash(opts *bind.CallOpts, orderAddresses [5]common.Address, orderValues [6]*big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "getOrderHash", orderAddresses, orderValues)
	return *ret0, err
}

// GetOrderHash is a free data retrieval call binding the contract method 0xcfc4d0ec.
//
// Solidity: function getOrderHash(orderAddresses address[5], orderValues uint256[6]) constant returns(bytes32)
func (_Exchange *ExchangeSession) GetOrderHash(orderAddresses [5]common.Address, orderValues [6]*big.Int) ([32]byte, error) {
	return _Exchange.Contract.GetOrderHash(&_Exchange.CallOpts, orderAddresses, orderValues)
}

// GetOrderHash is a free data retrieval call binding the contract method 0xcfc4d0ec.
//
// Solidity: function getOrderHash(orderAddresses address[5], orderValues uint256[6]) constant returns(bytes32)
func (_Exchange *ExchangeCallerSession) GetOrderHash(orderAddresses [5]common.Address, orderValues [6]*big.Int) ([32]byte, error) {
	return _Exchange.Contract.GetOrderHash(&_Exchange.CallOpts, orderAddresses, orderValues)
}

// GetPartialAmount is a free data retrieval call binding the contract method 0x98024a8b.
//
// Solidity: function getPartialAmount(numerator uint256, denominator uint256, target uint256) constant returns(uint256)
func (_Exchange *ExchangeCaller) GetPartialAmount(opts *bind.CallOpts, numerator *big.Int, denominator *big.Int, target *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "getPartialAmount", numerator, denominator, target)
	return *ret0, err
}

// GetPartialAmount is a free data retrieval call binding the contract method 0x98024a8b.
//
// Solidity: function getPartialAmount(numerator uint256, denominator uint256, target uint256) constant returns(uint256)
func (_Exchange *ExchangeSession) GetPartialAmount(numerator *big.Int, denominator *big.Int, target *big.Int) (*big.Int, error) {
	return _Exchange.Contract.GetPartialAmount(&_Exchange.CallOpts, numerator, denominator, target)
}

// GetPartialAmount is a free data retrieval call binding the contract method 0x98024a8b.
//
// Solidity: function getPartialAmount(numerator uint256, denominator uint256, target uint256) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) GetPartialAmount(numerator *big.Int, denominator *big.Int, target *big.Int) (*big.Int, error) {
	return _Exchange.Contract.GetPartialAmount(&_Exchange.CallOpts, numerator, denominator, target)
}

// GetUnavailableTakerTokenAmount is a free data retrieval call binding the contract method 0x7e9abb50.
//
// Solidity: function getUnavailableTakerTokenAmount(orderHash bytes32) constant returns(uint256)
func (_Exchange *ExchangeCaller) GetUnavailableTakerTokenAmount(opts *bind.CallOpts, orderHash [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "getUnavailableTakerTokenAmount", orderHash)
	return *ret0, err
}

// GetUnavailableTakerTokenAmount is a free data retrieval call binding the contract method 0x7e9abb50.
//
// Solidity: function getUnavailableTakerTokenAmount(orderHash bytes32) constant returns(uint256)
func (_Exchange *ExchangeSession) GetUnavailableTakerTokenAmount(orderHash [32]byte) (*big.Int, error) {
	return _Exchange.Contract.GetUnavailableTakerTokenAmount(&_Exchange.CallOpts, orderHash)
}

// GetUnavailableTakerTokenAmount is a free data retrieval call binding the contract method 0x7e9abb50.
//
// Solidity: function getUnavailableTakerTokenAmount(orderHash bytes32) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) GetUnavailableTakerTokenAmount(orderHash [32]byte) (*big.Int, error) {
	return _Exchange.Contract.GetUnavailableTakerTokenAmount(&_Exchange.CallOpts, orderHash)
}

// IsRoundingError is a free data retrieval call binding the contract method 0x14df96ee.
//
// Solidity: function isRoundingError(numerator uint256, denominator uint256, target uint256) constant returns(bool)
func (_Exchange *ExchangeCaller) IsRoundingError(opts *bind.CallOpts, numerator *big.Int, denominator *big.Int, target *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "isRoundingError", numerator, denominator, target)
	return *ret0, err
}

// IsRoundingError is a free data retrieval call binding the contract method 0x14df96ee.
//
// Solidity: function isRoundingError(numerator uint256, denominator uint256, target uint256) constant returns(bool)
func (_Exchange *ExchangeSession) IsRoundingError(numerator *big.Int, denominator *big.Int, target *big.Int) (bool, error) {
	return _Exchange.Contract.IsRoundingError(&_Exchange.CallOpts, numerator, denominator, target)
}

// IsRoundingError is a free data retrieval call binding the contract method 0x14df96ee.
//
// Solidity: function isRoundingError(numerator uint256, denominator uint256, target uint256) constant returns(bool)
func (_Exchange *ExchangeCallerSession) IsRoundingError(numerator *big.Int, denominator *big.Int, target *big.Int) (bool, error) {
	return _Exchange.Contract.IsRoundingError(&_Exchange.CallOpts, numerator, denominator, target)
}

// IsValidSignature is a free data retrieval call binding the contract method 0x8163681e.
//
// Solidity: function isValidSignature(signer address, hash bytes32, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_Exchange *ExchangeCaller) IsValidSignature(opts *bind.CallOpts, signer common.Address, hash [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "isValidSignature", signer, hash, v, r, s)
	return *ret0, err
}

// IsValidSignature is a free data retrieval call binding the contract method 0x8163681e.
//
// Solidity: function isValidSignature(signer address, hash bytes32, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_Exchange *ExchangeSession) IsValidSignature(signer common.Address, hash [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	return _Exchange.Contract.IsValidSignature(&_Exchange.CallOpts, signer, hash, v, r, s)
}

// IsValidSignature is a free data retrieval call binding the contract method 0x8163681e.
//
// Solidity: function isValidSignature(signer address, hash bytes32, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_Exchange *ExchangeCallerSession) IsValidSignature(signer common.Address, hash [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	return _Exchange.Contract.IsValidSignature(&_Exchange.CallOpts, signer, hash, v, r, s)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xbaa0181d.
//
// Solidity: function batchCancelOrders(orderAddresses address[5][], orderValues uint256[6][], cancelTakerTokenAmounts uint256[]) returns()
func (_Exchange *ExchangeTransactor) BatchCancelOrders(opts *bind.TransactOpts, orderAddresses [5]common.Address, orderValues [6]*big.Int, cancelTakerTokenAmounts []*big.Int) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchCancelOrders", orderAddresses, orderValues, cancelTakerTokenAmounts)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xbaa0181d.
//
// Solidity: function batchCancelOrders(orderAddresses address[5][], orderValues uint256[6][], cancelTakerTokenAmounts uint256[]) returns()
func (_Exchange *ExchangeSession) BatchCancelOrders(orderAddresses [5]common.Address, orderValues [6]*big.Int, cancelTakerTokenAmounts []*big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.BatchCancelOrders(&_Exchange.TransactOpts, orderAddresses, orderValues, cancelTakerTokenAmounts)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xbaa0181d.
//
// Solidity: function batchCancelOrders(orderAddresses address[5][], orderValues uint256[6][], cancelTakerTokenAmounts uint256[]) returns()
func (_Exchange *ExchangeTransactorSession) BatchCancelOrders(orderAddresses [5]common.Address, orderValues [6]*big.Int, cancelTakerTokenAmounts []*big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.BatchCancelOrders(&_Exchange.TransactOpts, orderAddresses, orderValues, cancelTakerTokenAmounts)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0x4f150787.
//
// Solidity: function batchFillOrKillOrders(orderAddresses address[5][], orderValues uint256[6][], fillTakerTokenAmounts uint256[], v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeTransactor) BatchFillOrKillOrders(opts *bind.TransactOpts, orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmounts []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchFillOrKillOrders", orderAddresses, orderValues, fillTakerTokenAmounts, v, r, s)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0x4f150787.
//
// Solidity: function batchFillOrKillOrders(orderAddresses address[5][], orderValues uint256[6][], fillTakerTokenAmounts uint256[], v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeSession) BatchFillOrKillOrders(orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmounts []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrKillOrders(&_Exchange.TransactOpts, orderAddresses, orderValues, fillTakerTokenAmounts, v, r, s)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0x4f150787.
//
// Solidity: function batchFillOrKillOrders(orderAddresses address[5][], orderValues uint256[6][], fillTakerTokenAmounts uint256[], v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeTransactorSession) BatchFillOrKillOrders(orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmounts []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrKillOrders(&_Exchange.TransactOpts, orderAddresses, orderValues, fillTakerTokenAmounts, v, r, s)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0xb7b2c7d6.
//
// Solidity: function batchFillOrders(orderAddresses address[5][], orderValues uint256[6][], fillTakerTokenAmounts uint256[], shouldThrowOnInsufficientBalanceOrAllowance bool, v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeTransactor) BatchFillOrders(opts *bind.TransactOpts, orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmounts []*big.Int, shouldThrowOnInsufficientBalanceOrAllowance bool, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchFillOrders", orderAddresses, orderValues, fillTakerTokenAmounts, shouldThrowOnInsufficientBalanceOrAllowance, v, r, s)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0xb7b2c7d6.
//
// Solidity: function batchFillOrders(orderAddresses address[5][], orderValues uint256[6][], fillTakerTokenAmounts uint256[], shouldThrowOnInsufficientBalanceOrAllowance bool, v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeSession) BatchFillOrders(orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmounts []*big.Int, shouldThrowOnInsufficientBalanceOrAllowance bool, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrders(&_Exchange.TransactOpts, orderAddresses, orderValues, fillTakerTokenAmounts, shouldThrowOnInsufficientBalanceOrAllowance, v, r, s)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0xb7b2c7d6.
//
// Solidity: function batchFillOrders(orderAddresses address[5][], orderValues uint256[6][], fillTakerTokenAmounts uint256[], shouldThrowOnInsufficientBalanceOrAllowance bool, v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeTransactorSession) BatchFillOrders(orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmounts []*big.Int, shouldThrowOnInsufficientBalanceOrAllowance bool, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrders(&_Exchange.TransactOpts, orderAddresses, orderValues, fillTakerTokenAmounts, shouldThrowOnInsufficientBalanceOrAllowance, v, r, s)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x394c21e7.
//
// Solidity: function cancelOrder(orderAddresses address[5], orderValues uint256[6], cancelTakerTokenAmount uint256) returns(uint256)
func (_Exchange *ExchangeTransactor) CancelOrder(opts *bind.TransactOpts, orderAddresses [5]common.Address, orderValues [6]*big.Int, cancelTakerTokenAmount *big.Int) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "cancelOrder", orderAddresses, orderValues, cancelTakerTokenAmount)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x394c21e7.
//
// Solidity: function cancelOrder(orderAddresses address[5], orderValues uint256[6], cancelTakerTokenAmount uint256) returns(uint256)
func (_Exchange *ExchangeSession) CancelOrder(orderAddresses [5]common.Address, orderValues [6]*big.Int, cancelTakerTokenAmount *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrder(&_Exchange.TransactOpts, orderAddresses, orderValues, cancelTakerTokenAmount)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x394c21e7.
//
// Solidity: function cancelOrder(orderAddresses address[5], orderValues uint256[6], cancelTakerTokenAmount uint256) returns(uint256)
func (_Exchange *ExchangeTransactorSession) CancelOrder(orderAddresses [5]common.Address, orderValues [6]*big.Int, cancelTakerTokenAmount *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrder(&_Exchange.TransactOpts, orderAddresses, orderValues, cancelTakerTokenAmount)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0x741bcc93.
//
// Solidity: function fillOrKillOrder(orderAddresses address[5], orderValues uint256[6], fillTakerTokenAmount uint256, v uint8, r bytes32, s bytes32) returns()
func (_Exchange *ExchangeTransactor) FillOrKillOrder(opts *bind.TransactOpts, orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmount *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "fillOrKillOrder", orderAddresses, orderValues, fillTakerTokenAmount, v, r, s)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0x741bcc93.
//
// Solidity: function fillOrKillOrder(orderAddresses address[5], orderValues uint256[6], fillTakerTokenAmount uint256, v uint8, r bytes32, s bytes32) returns()
func (_Exchange *ExchangeSession) FillOrKillOrder(orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmount *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrKillOrder(&_Exchange.TransactOpts, orderAddresses, orderValues, fillTakerTokenAmount, v, r, s)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0x741bcc93.
//
// Solidity: function fillOrKillOrder(orderAddresses address[5], orderValues uint256[6], fillTakerTokenAmount uint256, v uint8, r bytes32, s bytes32) returns()
func (_Exchange *ExchangeTransactorSession) FillOrKillOrder(orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmount *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrKillOrder(&_Exchange.TransactOpts, orderAddresses, orderValues, fillTakerTokenAmount, v, r, s)
}

// FillOrder is a paid mutator transaction binding the contract method 0xbc61394a.
//
// Solidity: function fillOrder(orderAddresses address[5], orderValues uint256[6], fillTakerTokenAmount uint256, shouldThrowOnInsufficientBalanceOrAllowance bool, v uint8, r bytes32, s bytes32) returns(filledTakerTokenAmount uint256)
func (_Exchange *ExchangeTransactor) FillOrder(opts *bind.TransactOpts, orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmount *big.Int, shouldThrowOnInsufficientBalanceOrAllowance bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "fillOrder", orderAddresses, orderValues, fillTakerTokenAmount, shouldThrowOnInsufficientBalanceOrAllowance, v, r, s)
}

// FillOrder is a paid mutator transaction binding the contract method 0xbc61394a.
//
// Solidity: function fillOrder(orderAddresses address[5], orderValues uint256[6], fillTakerTokenAmount uint256, shouldThrowOnInsufficientBalanceOrAllowance bool, v uint8, r bytes32, s bytes32) returns(filledTakerTokenAmount uint256)
func (_Exchange *ExchangeSession) FillOrder(orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmount *big.Int, shouldThrowOnInsufficientBalanceOrAllowance bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrder(&_Exchange.TransactOpts, orderAddresses, orderValues, fillTakerTokenAmount, shouldThrowOnInsufficientBalanceOrAllowance, v, r, s)
}

// FillOrder is a paid mutator transaction binding the contract method 0xbc61394a.
//
// Solidity: function fillOrder(orderAddresses address[5], orderValues uint256[6], fillTakerTokenAmount uint256, shouldThrowOnInsufficientBalanceOrAllowance bool, v uint8, r bytes32, s bytes32) returns(filledTakerTokenAmount uint256)
func (_Exchange *ExchangeTransactorSession) FillOrder(orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmount *big.Int, shouldThrowOnInsufficientBalanceOrAllowance bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrder(&_Exchange.TransactOpts, orderAddresses, orderValues, fillTakerTokenAmount, shouldThrowOnInsufficientBalanceOrAllowance, v, r, s)
}

// FillOrdersUpTo is a paid mutator transaction binding the contract method 0x363349be.
//
// Solidity: function fillOrdersUpTo(orderAddresses address[5][], orderValues uint256[6][], fillTakerTokenAmount uint256, shouldThrowOnInsufficientBalanceOrAllowance bool, v uint8[], r bytes32[], s bytes32[]) returns(uint256)
func (_Exchange *ExchangeTransactor) FillOrdersUpTo(opts *bind.TransactOpts, orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmount *big.Int, shouldThrowOnInsufficientBalanceOrAllowance bool, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "fillOrdersUpTo", orderAddresses, orderValues, fillTakerTokenAmount, shouldThrowOnInsufficientBalanceOrAllowance, v, r, s)
}

// FillOrdersUpTo is a paid mutator transaction binding the contract method 0x363349be.
//
// Solidity: function fillOrdersUpTo(orderAddresses address[5][], orderValues uint256[6][], fillTakerTokenAmount uint256, shouldThrowOnInsufficientBalanceOrAllowance bool, v uint8[], r bytes32[], s bytes32[]) returns(uint256)
func (_Exchange *ExchangeSession) FillOrdersUpTo(orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmount *big.Int, shouldThrowOnInsufficientBalanceOrAllowance bool, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrdersUpTo(&_Exchange.TransactOpts, orderAddresses, orderValues, fillTakerTokenAmount, shouldThrowOnInsufficientBalanceOrAllowance, v, r, s)
}

// FillOrdersUpTo is a paid mutator transaction binding the contract method 0x363349be.
//
// Solidity: function fillOrdersUpTo(orderAddresses address[5][], orderValues uint256[6][], fillTakerTokenAmount uint256, shouldThrowOnInsufficientBalanceOrAllowance bool, v uint8[], r bytes32[], s bytes32[]) returns(uint256)
func (_Exchange *ExchangeTransactorSession) FillOrdersUpTo(orderAddresses [5]common.Address, orderValues [6]*big.Int, fillTakerTokenAmount *big.Int, shouldThrowOnInsufficientBalanceOrAllowance bool, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrdersUpTo(&_Exchange.TransactOpts, orderAddresses, orderValues, fillTakerTokenAmount, shouldThrowOnInsufficientBalanceOrAllowance, v, r, s)
}
