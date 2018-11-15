package types

import (
	"encoding/json"
	// "encoding/hex"

	"fmt"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
	// "github.com/ethereum/go-ethereum/accounts/abi"
	"math/big"
	// "strconv"
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
	TakerFee                  *Uint256
	ExpirationTimestampInSec  *Uint256  `gorm:"index"`
	Salt                      *Uint256
	Signature                 Signature //`gorm:"type:bytea"`
	TakerAssetAmountFilled    *Uint256
	Cancelled                 bool
	PoolID                    []byte    `gorm:"index"`
}

func (order *Order) Initialize() {
	order.ExchangeAddress = &Address{}
	order.Maker = &Address{}
	order.Taker = &Address{}
	order.MakerAssetAddress = &Address{}
	order.TakerAssetAddress = &Address{}
	order.MakerAssetData = make([]byte, 20)
	order.TakerAssetData = make([]byte, 20)
	order.FeeRecipient = &Address{}
	order.SenderAddress = &Address{}
	order.MakerAssetAmount = &Uint256{}
	order.TakerAssetAmount = &Uint256{}
	order.MakerFee = &Uint256{}
	order.TakerFee = &Uint256{}
	order.ExpirationTimestampInSec = &Uint256{}
	order.Salt = &Uint256{}
	order.TakerAssetAmountFilled = &Uint256{}
	order.Cancelled = false
	order.Signature = Signature{}
}

// NewOrder takes string representations of values and converts them into an Order object
func NewOrder(maker, taker, makerToken, takerToken, feeRecipient, exchangeAddress, senderAddress, makerTokenAmount, takerTokenAmount, makerFee, takerFee, expirationTimestampInSec, salt, sig, takerTokenAmountFilled, cancelled string) (*Order, error) {
	order := Order{}
	if err := order.fromStrings(maker, taker, makerToken, takerToken, feeRecipient, exchangeAddress, senderAddress, makerTokenAmount, takerTokenAmount, makerFee, takerFee, expirationTimestampInSec, salt, sig, takerTokenAmountFilled, cancelled); err != nil {
		return nil, err
	}
	return &order, nil
}


func (order *Order) fromStrings(maker, taker, makerToken, takerToken, feeRecipient, exchangeAddress, senderAddress, makerTokenAmount, takerTokenAmount, makerFee, takerFee, expirationTimestampInSec, salt, sig, takerTokenAmountFilled, cancelled string) error {
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
	copy(order.TakerAssetData[:], takerAssetDataBytes)
	copy(order.MakerAssetData[:], makerAssetDataBytes)
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
	order.Cancelled = cancelled == "true"
	return nil
}

func (order *Order) Hash() []byte {

	eip191Header := []byte{25, 1}
	twelveNullBytes := [12]byte{}  // Addresses are 20 bytes, but the hashes expect 32, so we'll just add twelveNullBytes before each address
	domainSchemaSha := sha3.NewKeccak256()
	domainSchemaSha.Write([]byte("EIP712Domain(string name,string version,address verifyingContract)"))
	domainSha := sha3.NewKeccak256()
	nameSha := sha3.NewKeccak256()
	nameSha.Write([]byte("0x Protocol"))
	versionSha := sha3.NewKeccak256()
	versionSha.Write([]byte("2"))
	domainSha.Write(domainSchemaSha.Sum(nil))
	domainSha.Write(nameSha.Sum(nil))
	domainSha.Write(versionSha.Sum(nil))
	domainSha.Write(twelveNullBytes[:])
	domainSha.Write(order.ExchangeAddress[:])

	orderSchemaSha := sha3.NewKeccak256()
	orderSchemaSha.Write([]byte("Order(address makerAddress,address takerAddress,address feeRecipientAddress,address senderAddress,uint256 makerAssetAmount,uint256 takerAssetAmount,uint256 makerFee,uint256 takerFee,uint256 expirationTimeSeconds,uint256 salt,bytes makerAssetData,bytes takerAssetData)"))
	exchangeSha := sha3.NewKeccak256()
	exchangeSha.Write(order.ExchangeAddress[:])
	makerAssetDataSha := sha3.NewKeccak256()
	makerAssetDataSha.Write(order.MakerAssetData[:])
	takerAssetDataSha := sha3.NewKeccak256()
	takerAssetDataSha.Write(order.TakerAssetData[:])
	orderSha := sha3.NewKeccak256()
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

	sha := sha3.NewKeccak256()
	sha.Write(eip191Header)
	sha.Write(domainSha.Sum(nil))
	sha.Write(orderSha.Sum(nil))

	return sha.Sum(nil)
}

type jsonOrder struct {
	Maker                     string  `json:"makerAddress"`
	Taker                     string  `json:"takerAddress"`
	MakerAssetData            string  `json:"makerAssetData"`
	TakerAssetData            string  `json:"takerAssetData"`
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
	)
}

func (order *Order) MarshalJSON() ([]byte, error) {
	jsonOrder := &jsonOrder{}
	jsonOrder.Maker = fmt.Sprintf("%#x", order.Maker[:])
	jsonOrder.Taker = fmt.Sprintf("%#x", order.Taker[:])
	jsonOrder.MakerAssetData = fmt.Sprintf("%#x", order.MakerAssetData[:])
	jsonOrder.TakerAssetData = fmt.Sprintf("%#x", order.TakerAssetData[:])
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
	return json.Marshal(jsonOrder)
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
