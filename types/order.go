package types

import (
	"encoding/json"
	// "encoding/hex"
	"crypto/ecdsa"

	"fmt"
	"golang.org/x/crypto/sha3"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	// "github.com/ethereum/go-ethereum/accounts/abi"
	"math/big"
	"strconv"
	// "log"
)

// Order represents an 0x order object
type Order struct {
	Maker                     *Address  `gorm:"index"`
	Taker                     *Address  `gorm:"index"`
	MakerAssetAddress         *Address  `gorm:"index"`
	TakerAssetAddress         *Address  `gorm:"index"`
	MakerAssetData            AssetData `gorm:"index"`
	TakerAssetData            AssetData `gorm:"index"`
	FeeRecipient              *Address  `gorm:"index"`
	ExchangeAddress           *Address  `gorm:"index"`
	SenderAddress             *Address  `gorm:"index"`
	MakerAssetAmount          *Uint256
	TakerAssetAmount          *Uint256
	MakerFee                  *Uint256
	TakerFee                  *Uint256  `gorm:"index"`
	MakerFeeAssetData         AssetData `gorm:"index"`
	TakerFeeAssetData         AssetData `gorm:"index"`
	ExpirationTimestampInSec  *Uint256  `gorm:"index"`
	Salt                      *Uint256
	Signature                 Signature //`gorm:"type:bytea"`
	TakerAssetAmountFilled    *Uint256
	Cancelled                 bool
	PoolID                    []byte    `gorm:"index"`
	ChainID                   *Uint256  `gorm:"index"`
}

func (order *Order) Initialize() {
	order.ExchangeAddress = &Address{}
	order.Maker = &Address{}
	order.Taker = &Address{}
	order.MakerAssetAddress = &Address{}
	order.TakerAssetAddress = &Address{}
	order.MakerAssetData = []byte{}
	order.TakerAssetData = []byte{}
	order.MakerFeeAssetData = []byte{}
	order.TakerFeeAssetData = []byte{}
	order.FeeRecipient = &Address{}
	order.SenderAddress = &Address{}
	order.MakerAssetAmount = &Uint256{}
	order.TakerAssetAmount = &Uint256{}
	order.MakerFee = &Uint256{}
	order.TakerFee = &Uint256{}
	order.ExpirationTimestampInSec = &Uint256{}
	order.Salt = &Uint256{}
	order.TakerAssetAmountFilled = &Uint256{}
	order.ChainID = &Uint256{}
	order.Cancelled = false
	order.Signature = Signature{}
}

// NewOrder takes string representations of values and converts them into an Order object
func NewOrder(maker, taker, makerToken, takerToken, makerFeeAssetData, takerFeeAssetData, feeRecipient, exchangeAddress, senderAddress, makerTokenAmount, takerTokenAmount, makerFee, takerFee, expirationTimestampInSec, salt, sig, takerTokenAmountFilled, cancelled, chainid string) (*Order, error) {
	order := Order{}
	if err := order.fromStrings(maker, taker, makerToken, takerToken, makerFeeAssetData, takerFeeAssetData, feeRecipient, exchangeAddress, senderAddress, makerTokenAmount, takerTokenAmount, makerFee, takerFee, expirationTimestampInSec, salt, sig, takerTokenAmountFilled, cancelled, chainid); err != nil {
		return nil, err
	}
	return &order, nil
}


func (order *Order) fromStrings(maker, taker, makerToken, takerToken, makerFeeAssetData, takerFeeAssetData, feeRecipient, exchangeAddress, senderAddress, makerTokenAmount, takerTokenAmount, makerFee, takerFee, expirationTimestampInSec, salt, sig, takerTokenAmountFilled, cancelled, chainid string) error {
	order.Initialize()
	makerBytes, err := HexStringToBytes(maker)
	if err != nil {
		return err
	}
	takerBytes, err := HexStringToBytes(taker)
	if err != nil {
		return err
	}
	makerAssetDataBytes, err := HexStringToBytes(makerToken)
	if err != nil {
		return err
	}
	takerAssetDataBytes, err := HexStringToBytes(takerToken)
	if err != nil {
		return err
	}
	makerFeeAssetDataBytes, err := HexStringToBytes(makerFeeAssetData)
	if err != nil {
		return err
	}
	takerFeeAssetDataBytes, err := HexStringToBytes(takerFeeAssetData)
	if err != nil {
		return err
	}
	feeRecipientBytes, err := HexStringToBytes(feeRecipient)
	if err != nil {
		return err
	}
	exchangeAddressBytes, err := HexStringToBytes(exchangeAddress)
	if err != nil {
		return err
	}
	senderAddressBytes, err := HexStringToBytes(senderAddress)
	if err != nil {
		return err
	}
	signatureBytes, err := HexStringToBytes(sig)
	if err != nil {
		return err
	}
	makerTokenAmountBytes, err := intStringToBytes(makerTokenAmount)
	if err != nil {
		return err
	}
	takerTokenAmountBytes, err := intStringToBytes(takerTokenAmount)
	if err != nil {
		return err
	}
	chainIDBytes, err := intStringToBytes(chainid)
	if err != nil {
		return err
	}
	makerFeeBytes, err := intStringToBytes(makerFee)
	if err != nil {
		return err
	}
	takerFeeBytes, err := intStringToBytes(takerFee)
	if err != nil {
		return err
	}
	expirationTimestampInSecBytes, err := intStringToBytes(expirationTimestampInSec)
	if err != nil {
		return err
	}
	saltBytes, err := intStringToBytes(salt)
	if err != nil {
		return err
	}
	takerTokenAmountFilledBytes, err := intStringToBytes(takerTokenAmountFilled)
	if err != nil {
		return err
	}
	copy(order.Maker[:], makerBytes)
	copy(order.Taker[:], takerBytes)
	copy(order.FeeRecipient[:], feeRecipientBytes)
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	copy(order.SenderAddress[:], senderAddressBytes)
	order.TakerAssetData = make(AssetData, len(takerAssetDataBytes))
	order.MakerAssetData = make(AssetData, len(makerAssetDataBytes))
	order.TakerFeeAssetData = make(AssetData, len(takerFeeAssetDataBytes))
	order.MakerFeeAssetData = make(AssetData, len(makerFeeAssetDataBytes))
	copy(order.TakerAssetData[:], takerAssetDataBytes)
	copy(order.MakerAssetData[:], makerAssetDataBytes)
	copy(order.TakerFeeAssetData[:], takerFeeAssetDataBytes)
	copy(order.MakerFeeAssetData[:], makerFeeAssetDataBytes)
	order.TakerAssetAddress = order.TakerAssetData.Address()
	order.MakerAssetAddress = order.MakerAssetData.Address()
	copy(order.MakerAssetAmount[:], makerTokenAmountBytes)
	copy(order.TakerAssetAmount[:], takerTokenAmountBytes)
	copy(order.MakerFee[:], makerFeeBytes)
	copy(order.TakerFee[:], takerFeeBytes)
	copy(order.ExpirationTimestampInSec[:], expirationTimestampInSecBytes)
	copy(order.Salt[:], saltBytes)
	order.Signature = append(order.Signature, signatureBytes[:]...)
	copy(order.Signature[:], signatureBytes)
	copy(order.TakerAssetAmountFilled[:], takerTokenAmountFilledBytes)
	copy(order.ChainID[:], chainIDBytes)
	order.Cancelled = cancelled == "true"
	return nil
}

func (order *Order) Hash() []byte {

	eip191Header := []byte{25, 1}
	twelveNullBytes := [12]byte{}  // Addresses are 20 bytes, but the hashes expect 32, so we'll just add twelveNullBytes before each address
	domainSchemaSha := sha3.NewLegacyKeccak256()
	domainSchemaSha.Write([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))
	domainSha := sha3.NewLegacyKeccak256()
	nameSha := sha3.NewLegacyKeccak256()
	nameSha.Write([]byte("0x Protocol"))
	versionSha := sha3.NewLegacyKeccak256()
	versionSha.Write([]byte("3.0.0"))
	// log.Printf("domain schema sha: %#x", domainSchemaSha.Sum(nil))
	domainSha.Write(domainSchemaSha.Sum(nil))
	domainSha.Write(nameSha.Sum(nil))
	domainSha.Write(versionSha.Sum(nil))
	domainSha.Write(order.ChainID[:])
	domainSha.Write(twelveNullBytes[:])
	domainSha.Write(order.ExchangeAddress[:])
	// log.Printf("domain sha: %#x", domainSha.Sum(nil))

	orderSchemaSha := sha3.NewLegacyKeccak256()
	orderSchemaSha.Write([]byte("Order(address makerAddress,address takerAddress,address feeRecipientAddress,address senderAddress,uint256 makerAssetAmount,uint256 takerAssetAmount,uint256 makerFee,uint256 takerFee,uint256 expirationTimeSeconds,uint256 salt,bytes makerAssetData,bytes takerAssetData,bytes makerFeeAssetData,bytes takerFeeAssetData)"))
	// log.Printf("order schema sha: %#x", orderSchemaSha.Sum(nil))

	exchangeSha := sha3.NewLegacyKeccak256()
	exchangeSha.Write(order.ExchangeAddress[:])
	makerAssetDataSha := sha3.NewLegacyKeccak256()
	makerAssetDataSha.Write(order.MakerAssetData[:])
	takerAssetDataSha := sha3.NewLegacyKeccak256()
	takerAssetDataSha.Write(order.TakerAssetData[:])
	makerFeeAssetDataSha := sha3.NewLegacyKeccak256()
	makerFeeAssetDataSha.Write(order.MakerFeeAssetData[:])
	takerFeeAssetDataSha := sha3.NewLegacyKeccak256()
	takerFeeAssetDataSha.Write(order.TakerFeeAssetData[:])
	orderSha := sha3.NewLegacyKeccak256()
	orderSha.Write(orderSchemaSha.Sum(nil))
	orderSha.Write(twelveNullBytes[:])
	orderSha.Write(order.Maker[:])
	orderSha.Write(twelveNullBytes[:])
	orderSha.Write(order.Taker[:])
	orderSha.Write(twelveNullBytes[:])
	orderSha.Write(order.FeeRecipient[:])
	orderSha.Write(twelveNullBytes[:])
	orderSha.Write(order.SenderAddress[:])
	orderSha.Write(order.MakerAssetAmount[:])
	orderSha.Write(order.TakerAssetAmount[:])
	orderSha.Write(order.MakerFee[:])
	orderSha.Write(order.TakerFee[:])
	orderSha.Write(order.ExpirationTimestampInSec[:])
	orderSha.Write(order.Salt[:])
	orderSha.Write(makerAssetDataSha.Sum(nil))
	orderSha.Write(takerAssetDataSha.Sum(nil))
	orderSha.Write(makerFeeAssetDataSha.Sum(nil))
	orderSha.Write(takerFeeAssetDataSha.Sum(nil))

	sha := sha3.NewLegacyKeccak256()
	sha.Write(eip191Header)
	sha.Write(domainSha.Sum(nil))
	sha.Write(orderSha.Sum(nil))

	return sha.Sum(nil)
}

type jsonOrder struct {
	ChainID                   int64   `json:"chainId"`
	Maker                     string  `json:"makerAddress"`
	Taker                     string  `json:"takerAddress"`
	MakerAssetData            string  `json:"makerAssetData"`
	TakerAssetData            string  `json:"takerAssetData"`
	MakerFeeAssetData         string  `json:"makerFeeAssetData"`
	TakerFeeAssetData         string  `json:"takerFeeAssetData"`
	FeeRecipient              string  `json:"feeRecipientAddress"`
	ExchangeAddress           string  `json:"exchangeAddress"`
	SenderAddress             string  `json:"senderAddress"`
	MakerAssetAmount          string  `json:"makerAssetAmount"`
	TakerAssetAmount          string  `json:"takerAssetAmount"`
	MakerFee                  string  `json:"makerFee"`
	TakerFee                  string  `json:"takerFee"`
	ExpirationTimestampInSec  string  `json:"expirationTimeSeconds"`
	Salt                      string  `json:"salt"`
	Signature                 string  `json:"signature"`
	TakerAssetAmountFilled    string  `json:"-"`
	Cancelled                 string  `json:"-"`
}

func (order *Order) UnmarshalJSON(b []byte) error {
	jOrder := jsonOrder{}
	if err := json.Unmarshal(b, &jOrder); err != nil {
		return err
	}
	if jOrder.TakerAssetAmountFilled == "" {
		jOrder.TakerAssetAmountFilled = "0"
	}
	if jOrder.Cancelled == "" {
		jOrder.Cancelled = "false"
	}
	return order.fromStrings(
		jOrder.Maker,
		jOrder.Taker,
		jOrder.MakerAssetData,
		jOrder.TakerAssetData,
		jOrder.MakerFeeAssetData,
		jOrder.TakerFeeAssetData,
		jOrder.FeeRecipient,
		jOrder.ExchangeAddress,
		jOrder.SenderAddress,
		jOrder.MakerAssetAmount,
		jOrder.TakerAssetAmount,
		jOrder.MakerFee,
		jOrder.TakerFee,
		jOrder.ExpirationTimestampInSec,
		jOrder.Salt,
		jOrder.Signature,
		jOrder.TakerAssetAmountFilled,
		jOrder.Cancelled,
		strconv.Itoa(int(jOrder.ChainID)),
	)
}

func (order *Order) MarshalJSON() ([]byte, error) {
	jsonOrder := &jsonOrder{}
	jsonOrder.Maker = fmt.Sprintf("%#x", order.Maker[:])
	jsonOrder.Taker = fmt.Sprintf("%#x", order.Taker[:])
	jsonOrder.MakerAssetData = fmt.Sprintf("%#x", order.MakerAssetData[:])
	jsonOrder.TakerAssetData = fmt.Sprintf("%#x", order.TakerAssetData[:])
	jsonOrder.MakerFeeAssetData = fmt.Sprintf("%#x", order.MakerFeeAssetData[:])
	jsonOrder.TakerFeeAssetData = fmt.Sprintf("%#x", order.TakerFeeAssetData[:])
	jsonOrder.FeeRecipient = fmt.Sprintf("%#x", order.FeeRecipient[:])
	jsonOrder.ExchangeAddress = fmt.Sprintf("%#x", order.ExchangeAddress[:])
	jsonOrder.SenderAddress = fmt.Sprintf("%#x", order.SenderAddress[:])
	jsonOrder.MakerAssetAmount = new(big.Int).SetBytes(order.MakerAssetAmount[:]).String()
	jsonOrder.TakerAssetAmount = new(big.Int).SetBytes(order.TakerAssetAmount[:]).String()
	jsonOrder.MakerFee = new(big.Int).SetBytes(order.MakerFee[:]).String()
	jsonOrder.TakerFee = new(big.Int).SetBytes(order.TakerFee[:]).String()
	jsonOrder.ExpirationTimestampInSec = new(big.Int).SetBytes(order.ExpirationTimestampInSec[:]).String()
	jsonOrder.Salt = new(big.Int).SetBytes(order.Salt[:]).String()
	jsonOrder.Signature = fmt.Sprintf("%#x", order.Signature)
	jsonOrder.TakerAssetAmountFilled = new(big.Int).SetBytes(order.TakerAssetAmountFilled[:]).String()
	if order.Cancelled {
		jsonOrder.Cancelled = "true"
	} else {
		jsonOrder.Cancelled = "false"
	}
	jsonOrder.ChainID = order.ChainID.Big().Int64()
	return json.Marshal(jsonOrder)
}

func (order *Order) Sign(key *ecdsa.PrivateKey, sigType byte) error {
	address := crypto.PubkeyToAddress(key.PublicKey)
	copy(order.Maker[:], address[:])
	var signedBytes []byte
	switch sigType {
	case SigTypeEthSign:
		hashedBytes := append([]byte("\x19Ethereum Signed Message:\n32"), order.Hash()...)
		signedBytes = crypto.Keccak256(hashedBytes)
	case SigTypeEIP712:
		signedBytes = order.Hash()
	default:
		return fmt.Errorf("Unsupported signature type: %v", sigType)
	}
	sig, _ := crypto.Sign(signedBytes, key)
	order.Signature = make(Signature, 66)
	order.Signature[0] = sig[64] + 27
	copy(order.Signature[1:33], sig[0:32])
	copy(order.Signature[33:65], sig[32:64])
	order.Signature[65] = SigTypeEthSign
	return nil
}

func (order *Order) Bytes() []byte {
	data, _ := rlp.EncodeToBytes(order)
	return data
}

func (order *Order) FromBytes(data []byte) (error) {
	return rlp.DecodeBytes(data, order)
}

func OrderFromBytes(data []byte) (*Order, error)  {
	order := Order{}
	return &order, order.FromBytes(data)
}
