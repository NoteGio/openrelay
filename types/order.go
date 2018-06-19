package types

import (
	"encoding/json"
	// "encoding/hex"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	// "github.com/ethereum/go-ethereum/accounts/abi"
	"math/big"
	"strconv"
	// "log"
)

type Address [20]byte

func (addr *Address) Value() (driver.Value, error) {
	return addr[:], nil
}

func (addr *Address) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		copy(addr[:], v)
		return nil
	default:
		return errors.New("Address scanner src should be []byte")
	}
}

func (addr *Address) String() string {
	return fmt.Sprintf("%#x", addr[:])
}

type Uint256 [32]byte

func (data *Uint256) Value() (driver.Value, error) {
	return data[:], nil
}

func (data *Uint256) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		copy(data[:], v)
		return nil
	default:
		return errors.New("Uint256 scanner src should be []byte")
	}
}

func (data *Uint256) String() (string) {
	return data.Big().String()
}

func (data *Uint256) Big() (*big.Int) {
	return new(big.Int).SetBytes(data[:])
}


// Order represents an 0x order object
type Order struct {
	Maker                     *Address `gorm:"index"`
	Taker                     *Address `gorm:"index"`
	MakerAsset                *Address `gorm:"index"`
	TakerAsset                *Address `gorm:"index"`
	MakerAssetData            []byte   `gorm:"index"`
	TakerAssetData            []byte   `gorm:"index"`
	FeeRecipient              *Address `gorm:"index"`
	ExchangeAddress           *Address `gorm:"index"`
	SenderAddress             *Address `gorm:"index"`
	MakerAssetAmount          *Uint256
	TakerAssetAmount          *Uint256
	MakerFee                  *Uint256
	TakerFee                  *Uint256
	ExpirationTimestampInSec  *Uint256 `gorm:"index"`
	Salt                      *Uint256
	Signature                 *[]byte //`gorm:"type:bytea"`
	TakerAssetAmountFilled    *Uint256
	TakerAssetAmountCancelled *Uint256
}

func (order *Order) Initialize() {
	order.ExchangeAddress = &Address{}
	order.Maker = &Address{}
	order.Taker = &Address{}
	order.MakerAsset = &Address{}
	order.TakerAsset = &Address{}
	order.MakerAssetData = []byte{}
	order.TakerAssetData = []byte{}
	order.FeeRecipient = &Address{}
	order.MakerAssetAmount = &Uint256{}
	order.TakerAssetAmount = &Uint256{}
	order.MakerFee = &Uint256{}
	order.TakerFee = &Uint256{}
	order.ExpirationTimestampInSec = &Uint256{}
	order.Salt = &Uint256{}
	order.TakerAssetAmountFilled = &Uint256{}
	order.TakerAssetAmountCancelled = &Uint256{}
	order.Signature = &Signature{}
}

// NewOrder takes string representations of values and converts them into an Order object
func NewOrder(maker, taker, makerToken, takerToken, feeRecipient, exchangeAddress, makerTokenAmount, takerTokenAmount, makerFee, takerFee, expirationTimestampInSec, salt, sigV, sigR, sigS, takerTokenAmountFilled, takerTokenAmountCancelled string) (*Order, error) {
	order := Order{}
	if err := order.fromStrings(maker, taker, makerToken, takerToken, feeRecipient, exchangeAddress, makerTokenAmount, takerTokenAmount, makerFee, takerFee, expirationTimestampInSec, salt, sigV, sigR, sigS, takerTokenAmountFilled, takerTokenAmountCancelled); err != nil {
		return nil, err
	}
	return &order, nil
}


jOrder.Maker,
jOrder.Taker,
jOrder.MakerAsset,
jOrder.TakerAsset,
jOrder.FeeRecipient,
jOrder.ExchangeAddress,
jOrder.MakerAssetAmount,
jOrder.TakerAssetAmount,
jOrder.MakerFee,
jOrder.TakerFee,
jOrder.ExpirationTimestampInSec,
jOrder.Salt,
jOrder.Signature,
jOrder.TakerAssetAmountFilled,
jOrder.TakerAssetAmountCancelled,

func (order *Order) fromStrings(maker, taker, makerToken, takerToken, feeRecipient, exchangeAddress, makerTokenAmount, takerTokenAmount, makerFee, takerFee, expirationTimestampInSec, salt, sigV, sigR, sigS, takerTokenAmountFilled, takerTokenAmountCancelled string) error {
	order.Initialize()
	makerBytes, err := HexStringToBytes(maker)
	if err != nil {
		return err
	}
	takerBytes, err := HexStringToBytes(taker)
	if err != nil {
		return err
	}
	makerTokenBytes, err := HexStringToBytes(makerToken)
	if err != nil {
		return err
	}
	takerTokenBytes, err := HexStringToBytes(takerToken)
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
	sigVInt, err := strconv.Atoi(sigV)
	if err != nil {
		return err
	}
	sigRBytes, err := HexStringToBytes(sigR)
	if err != nil {
		return err
	}
	sigSBytes, err := HexStringToBytes(sigS)
	if err != nil {
		return err
	}
	takerTokenAmountFilledBytes, err := intStringToBytes(takerTokenAmountFilled)
	if err != nil {
		return err
	}
	takerTokenAmountCancelledBytes, err := intStringToBytes(takerTokenAmountCancelled)
	if err != nil {
		return err
	}
	copy(order.Maker[:], makerBytes)
	copy(order.Taker[:], takerBytes)
	copy(order.MakerAsset[:], makerTokenBytes)
	copy(order.TakerAsset[:], takerTokenBytes)
	copy(order.FeeRecipient[:], feeRecipientBytes)
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	copy(order.MakerAssetAmount[:], makerTokenAmountBytes)
	copy(order.TakerAssetAmount[:], takerTokenAmountBytes)
	copy(order.MakerFee[:], makerFeeBytes)
	copy(order.TakerFee[:], takerFeeBytes)
	copy(order.ExpirationTimestampInSec[:], expirationTimestampInSecBytes)
	copy(order.Salt[:], saltBytes)
	order.Signature = &Signature{}
	order.Signature.V = byte(sigVInt)
	copy(order.Signature.S[:], sigSBytes)
	copy(order.Signature.R[:], sigRBytes)
	copy(order.Signature.Hash[:], order.Hash())
	copy(order.TakerAssetAmountFilled[:], takerTokenAmountFilledBytes)
	copy(order.TakerAssetAmountCancelled[:], takerTokenAmountCancelledBytes)
	return nil
}

func (order *Order) Hash() []byte {

// const expectedOrderHash = '0x367ad7730eb8b5feab8a9c9f47c6fcba77a2d4df125ee6a59cc26ac955710f7e';
// const fakeExchangeContractAddress = '0xb69e673309512a9d726f87304c6984054f87a93b';
// const order: Order = {
// makerAddress: constants.NULL_ADDRESS,
// takerAddress: constants.NULL_ADDRESS,
// senderAddress: constants.NULL_ADDRESS,
// feeRecipientAddress: constants.NULL_ADDRESS,
// makerAssetData: constants.NULL_ADDRESS,
// takerAssetData: constants.NULL_ADDRESS,
// exchangeAddress: fakeExchangeContractAddress,
// salt: new BigNumber(0),
// makerFee: new BigNumber(0),
// takerFee: new BigNumber(0),
// makerAssetAmount: new BigNumber(0),
// takerAssetAmount: new BigNumber(0),
// expirationTimeSeconds: new BigNumber(0),
// }

	sha := sha3.NewKeccak256()
	sha.Write(order.ExchangeAddress[:])
	sha.Write(order.Maker[:])
	sha.Write(order.Taker[:])
	sha.Write(order.MakerAsset[:])
	sha.Write(order.TakerAsset[:])
	sha.Write(order.FeeRecipient[:])
	sha.Write(order.MakerAssetAmount[:])
	sha.Write(order.TakerAssetAmount[:])
	sha.Write(order.MakerFee[:])
	sha.Write(order.TakerFee[:])
	sha.Write(order.ExpirationTimestampInSec[:])
	sha.Write(order.Salt[:])
	return sha.Sum(nil)
}

// "makerAddress"
// "takerAddress"
// "feeRecipientAddress"
// "senderAddress"
// "makerAssetAmount"
// "takerAssetAmount"
// "makerFee"
// "takerFee"
// "expirationTimeSeconds"
// "salt"
// "makerAssetData"
// "takerAssetData"
// "exchangeAddress"
// "signature"

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
	TakerAssetAmountCancelled string  `json:"-"`
}

func (order *Order) UnmarshalJSON(b []byte) error {
	jOrder := jsonOrder{}
	if err := json.Unmarshal(b, &jOrder); err != nil {
		return err
	}
	if jOrder.TakerAssetAmountFilled == "" {
		jOrder.TakerAssetAmountFilled = "0"
	}
	if jOrder.TakerAssetAmountCancelled == "" {
		jOrder.TakerAssetAmountCancelled = "0"
	}
	return order.fromStrings(
		jOrder.Maker,
		jOrder.Taker,
		jOrder.MakerAsset,
		jOrder.TakerAsset,
		jOrder.FeeRecipient,
		jOrder.ExchangeAddress,
		jOrder.MakerAssetAmount,
		jOrder.TakerAssetAmount,
		jOrder.MakerFee,
		jOrder.TakerFee,
		jOrder.ExpirationTimestampInSec,
		jOrder.Salt,
		jOrder.Signature,
		jOrder.TakerAssetAmountFilled,
		jOrder.TakerAssetAmountCancelled,
	)
}

func (order *Order) MarshalJSON() ([]byte, error) {
	jsonOrder := &jsonOrder{}
	jsonOrder.Maker = fmt.Sprintf("%#x", order.Maker[:])
	jsonOrder.Taker = fmt.Sprintf("%#x", order.Taker[:])
	jsonOrder.MakerAsset = fmt.Sprintf("%#x", order.MakerAsset[:])
	jsonOrder.TakerAsset = fmt.Sprintf("%#x", order.TakerAsset[:])
	jsonOrder.FeeRecipient = fmt.Sprintf("%#x", order.FeeRecipient[:])
	jsonOrder.ExchangeAddress = fmt.Sprintf("%#x", order.ExchangeAddress[:])
	jsonOrder.MakerAssetAmount = new(big.Int).SetBytes(order.MakerAssetAmount[:]).String()
	jsonOrder.TakerAssetAmount = new(big.Int).SetBytes(order.TakerAssetAmount[:]).String()
	jsonOrder.MakerFee = new(big.Int).SetBytes(order.MakerFee[:]).String()
	jsonOrder.TakerFee = new(big.Int).SetBytes(order.TakerFee[:]).String()
	jsonOrder.ExpirationTimestampInSec = new(big.Int).SetBytes(order.ExpirationTimestampInSec[:]).String()
	jsonOrder.Salt = new(big.Int).SetBytes(order.Salt[:]).String()
	jsonOrder.Signature = jsonSignature{}
	jsonOrder.Signature.R = fmt.Sprintf("%#x", order.Signature.R[:])
	jsonOrder.Signature.V = json.Number(fmt.Sprintf("%v", order.Signature.V))
	jsonOrder.Signature.S = fmt.Sprintf("%#x", order.Signature.S[:])
	jsonOrder.TakerAssetAmountFilled = new(big.Int).SetBytes(order.TakerAssetAmountFilled[:]).String()
	jsonOrder.TakerAssetAmountCancelled = new(big.Int).SetBytes(order.TakerAssetAmountCancelled[:]).String()
	return json.Marshal(jsonOrder)
}

func (order *Order) Bytes() [441]byte {
	var output [441]byte
	copy(output[0:20], order.ExchangeAddress[:])             // 20
	copy(output[20:40], order.Maker[:])                      // 20
	copy(output[40:60], order.Taker[:])                      // 20
	copy(output[60:80], order.MakerAsset[:])                 // 20
	copy(output[80:100], order.TakerAsset[:])                // 20
	copy(output[100:120], order.FeeRecipient[:])             // 20
	copy(output[120:152], order.MakerAssetAmount[:])         // 32
	copy(output[152:184], order.TakerAssetAmount[:])         // 32
	copy(output[184:216], order.MakerFee[:])                 // 32
	copy(output[216:248], order.TakerFee[:])                 // 32
	copy(output[248:280], order.ExpirationTimestampInSec[:]) // 32
	copy(output[280:312], order.Salt[:])                     // 32
	output[312] = order.Signature.V
	copy(output[313:345], order.Signature.R[:])
	copy(output[345:377], order.Signature.S[:])
	copy(output[377:409], order.TakerAssetAmountFilled[:])
	copy(output[409:441], order.TakerAssetAmountCancelled[:])
	return output
}

func (order *Order) FromBytes(data [441]byte) {
	order.Initialize()
	copy(order.ExchangeAddress[:], data[0:20])
	copy(order.Maker[:], data[20:40])
	copy(order.Taker[:], data[40:60])
	copy(order.MakerAsset[:], data[60:80])
	copy(order.TakerAsset[:], data[80:100])
	copy(order.FeeRecipient[:], data[100:120])
	copy(order.MakerAssetAmount[:], data[120:152])
	copy(order.TakerAssetAmount[:], data[152:184])
	copy(order.MakerFee[:], data[184:216])
	copy(order.TakerFee[:], data[216:248])
	copy(order.ExpirationTimestampInSec[:], data[248:280])
	copy(order.Salt[:], data[280:312])
	order.Signature = &Signature{}
	order.Signature.V = data[312]
	copy(order.Signature.R[:], data[313:345])
	copy(order.Signature.S[:], data[345:377])
	copy(order.Signature.Hash[:], order.Hash())
	copy(order.TakerAssetAmountFilled[:], data[377:409])
	copy(order.TakerAssetAmountCancelled[:], data[409:441])
}

func OrderFromBytes(data [441]byte) *Order {
	order := Order{}
	order.FromBytes(data)
	return &order
}
