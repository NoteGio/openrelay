package types

import (
	"encoding/json"
	// "encoding/hex"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"strconv"
	"errors"
	"database/sql/driver"
	"fmt"
	"math/big"
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

type Uint256 [32]byte

func (data *Uint256) Value() (driver.Value, error) {
	return data[:], nil
}

func (addr *Uint256) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		copy(addr[:], v)
		return nil
	default:
		return errors.New("Uint256 scanner src should be []byte")
	}
}

// Order represents an 0x order object
type Order struct {
	Maker                     *Address `gorm:"index"`
	Taker                     *Address `gorm:"index"`
	MakerToken                *Address `gorm:"index"`
	TakerToken                *Address `gorm:"index"`
	FeeRecipient              *Address `gorm:"index"`
	ExchangeAddress           *Address `gorm:"index"`
	MakerTokenAmount          *Uint256
	TakerTokenAmount          *Uint256
	MakerFee                  *Uint256
	TakerFee                  *Uint256
	ExpirationTimestampInSec  *Uint256 `gorm:"index"`
	Salt                      *Uint256
	Signature                 *Signature `gorm:"type:bytea"`
	TakerTokenAmountFilled    *Uint256
	TakerTokenAmountCancelled *Uint256
}

func (order *Order) Initialize() {
	order.ExchangeAddress = &Address{}
	order.Maker = &Address{}
	order.Taker = &Address{}
	order.MakerToken = &Address{}
	order.TakerToken = &Address{}
	order.FeeRecipient = &Address{}
	order.MakerTokenAmount = &Uint256{}
	order.TakerTokenAmount = &Uint256{}
	order.MakerFee = &Uint256{}
	order.TakerFee = &Uint256{}
	order.ExpirationTimestampInSec = &Uint256{}
	order.Salt = &Uint256{}
	order.TakerTokenAmountFilled = &Uint256{}
	order.TakerTokenAmountCancelled = &Uint256{}
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
	copy(order.MakerToken[:], makerTokenBytes)
	copy(order.TakerToken[:], takerTokenBytes)
	copy(order.FeeRecipient[:], feeRecipientBytes)
	copy(order.ExchangeAddress[:], exchangeAddressBytes)
	copy(order.MakerTokenAmount[:], makerTokenAmountBytes)
	copy(order.TakerTokenAmount[:], takerTokenAmountBytes)
	copy(order.MakerFee[:], makerFeeBytes)
	copy(order.TakerFee[:], takerFeeBytes)
	copy(order.ExpirationTimestampInSec[:], expirationTimestampInSecBytes)
	copy(order.Salt[:], saltBytes)
	order.Signature = &Signature{}
	order.Signature.V = byte(sigVInt)
	copy(order.Signature.S[:], sigSBytes)
	copy(order.Signature.R[:], sigRBytes)
	copy(order.Signature.Hash[:], order.Hash())
	copy(order.TakerTokenAmountFilled[:], takerTokenAmountFilledBytes)
	copy(order.TakerTokenAmountCancelled[:], takerTokenAmountCancelledBytes)
	return nil
}

func (order *Order) Hash() []byte {
	sha := sha3.NewKeccak256()

	sha.Write(order.ExchangeAddress[:])
	sha.Write(order.Maker[:])
	sha.Write(order.Taker[:])
	sha.Write(order.MakerToken[:])
	sha.Write(order.TakerToken[:])
	sha.Write(order.FeeRecipient[:])
	sha.Write(order.MakerTokenAmount[:])
	sha.Write(order.TakerTokenAmount[:])
	sha.Write(order.MakerFee[:])
	sha.Write(order.TakerFee[:])
	sha.Write(order.ExpirationTimestampInSec[:])
	sha.Write(order.Salt[:])
	return sha.Sum(nil)
}

type jsonOrder struct {
	Maker                     string        `json:"maker"`
	Taker                     string        `json:"taker"`
	MakerToken                string        `json:"makerTokenAddress"`
	TakerToken                string        `json:"takerTokenAddress"`
	FeeRecipient              string        `json:"feeRecipient"`
	ExchangeAddress           string        `json:"exchangeContractAddress"`
	MakerTokenAmount          string        `json:"makerTokenAmount"`
	TakerTokenAmount          string        `json:"takerTokenAmount"`
	MakerFee                  string        `json:"makerFee"`
	TakerFee                  string        `json:"takerFee"`
	ExpirationTimestampInSec  string        `json:"expirationUnixTimestampSec"`
	Salt                      string        `json:"salt"`
	Signature                 jsonSignature `json:"ecSignature"`
	TakerTokenAmountFilled    string        `json:"takerTokenAmountFilled"`
	TakerTokenAmountCancelled string        `json:"takerTokenAmountCancelled"`
}

func (order *Order) UnmarshalJSON(b []byte) error {
	jOrder := jsonOrder{}
	if err := json.Unmarshal(b, &jOrder); err != nil {
		return err
	}
	if jOrder.TakerTokenAmountFilled == "" {
		jOrder.TakerTokenAmountFilled = "0"
	}
	if jOrder.TakerTokenAmountCancelled == "" {
		jOrder.TakerTokenAmountCancelled = "0"
	}
	order.fromStrings(
		jOrder.Maker,
		jOrder.Taker,
		jOrder.MakerToken,
		jOrder.TakerToken,
		jOrder.FeeRecipient,
		jOrder.ExchangeAddress,
		jOrder.MakerTokenAmount,
		jOrder.TakerTokenAmount,
		jOrder.MakerFee,
		jOrder.TakerFee,
		jOrder.ExpirationTimestampInSec,
		jOrder.Salt,
		string(jOrder.Signature.V),
		jOrder.Signature.R,
		jOrder.Signature.S,
		jOrder.TakerTokenAmountFilled,
		jOrder.TakerTokenAmountCancelled,
	)

	return nil
}

func (order *Order) MarshalJSON() ([]byte, error) {
	jsonOrder := &jsonOrder{}
	jsonOrder.Maker = fmt.Sprintf("%#x", order.Maker[:])
	jsonOrder.Taker = fmt.Sprintf("%#x", order.Taker[:])
	jsonOrder.MakerToken = fmt.Sprintf("%#x", order.MakerToken[:])
	jsonOrder.TakerToken = fmt.Sprintf("%#x", order.TakerToken[:])
	jsonOrder.FeeRecipient = fmt.Sprintf("%#x", order.FeeRecipient[:])
	jsonOrder.ExchangeAddress = fmt.Sprintf("%#x", order.ExchangeAddress[:])
	jsonOrder.MakerTokenAmount = new(big.Int).SetBytes(order.MakerTokenAmount[:]).String()
	jsonOrder.TakerTokenAmount = new(big.Int).SetBytes(order.TakerTokenAmount[:]).String()
	jsonOrder.MakerFee = new(big.Int).SetBytes(order.MakerFee[:]).String()
	jsonOrder.TakerFee = new(big.Int).SetBytes(order.TakerFee[:]).String()
	jsonOrder.ExpirationTimestampInSec = new(big.Int).SetBytes(order.ExpirationTimestampInSec[:]).String()
	jsonOrder.Salt = new(big.Int).SetBytes(order.Salt[:]).String()
	jsonOrder.Signature = jsonSignature{}
	jsonOrder.Signature.R = fmt.Sprintf("%#x", order.Signature.R[:])
	jsonOrder.Signature.V = json.Number(fmt.Sprintf("%v", order.Signature.V))
	jsonOrder.Signature.S = fmt.Sprintf("%#x", order.Signature.S[:])
	jsonOrder.TakerTokenAmountFilled = new(big.Int).SetBytes(order.TakerTokenAmountFilled[:]).String()
	jsonOrder.TakerTokenAmountCancelled = new(big.Int).SetBytes(order.TakerTokenAmountCancelled[:]).String()
	return json.Marshal(jsonOrder)
}

func (order *Order) Bytes() [441]byte {
	var output [441]byte
	copy(output[0:20], order.ExchangeAddress[:])             // 20
	copy(output[20:40], order.Maker[:])                      // 20
	copy(output[40:60], order.Taker[:])                      // 20
	copy(output[60:80], order.MakerToken[:])                 // 20
	copy(output[80:100], order.TakerToken[:])                // 20
	copy(output[100:120], order.FeeRecipient[:])             // 20
	copy(output[120:152], order.MakerTokenAmount[:])         // 32
	copy(output[152:184], order.TakerTokenAmount[:])         // 32
	copy(output[184:216], order.MakerFee[:])                 // 32
	copy(output[216:248], order.TakerFee[:])                 // 32
	copy(output[248:280], order.ExpirationTimestampInSec[:]) // 32
	copy(output[280:312], order.Salt[:])                     // 32
	output[312] = order.Signature.V
	copy(output[313:345], order.Signature.R[:])
	copy(output[345:377], order.Signature.S[:])
	copy(output[377:409], order.TakerTokenAmountFilled[:])
	copy(output[409:441], order.TakerTokenAmountCancelled[:])
	return output
}

func (order *Order) FromBytes(data [441]byte) {
	order.Initialize()
	copy(order.ExchangeAddress[:], data[0:20])
	copy(order.Maker[:], data[20:40])
	copy(order.Taker[:], data[40:60])
	copy(order.MakerToken[:], data[60:80])
	copy(order.TakerToken[:], data[80:100])
	copy(order.FeeRecipient[:], data[100:120])
	copy(order.MakerTokenAmount[:], data[120:152])
	copy(order.TakerTokenAmount[:], data[152:184])
	copy(order.MakerFee[:], data[184:216])
	copy(order.TakerFee[:], data[216:248])
	copy(order.ExpirationTimestampInSec[:], data[248:280])
	copy(order.Salt[:], data[280:312])
	order.Signature = &Signature{}
	order.Signature.V = data[312]
	copy(order.Signature.R[:], data[313:345])
	copy(order.Signature.S[:], data[345:377])
	copy(order.Signature.Hash[:], order.Hash())
	copy(order.TakerTokenAmountFilled[:], data[377:409])
	copy(order.TakerTokenAmountCancelled[:], data[409:441])
}

func OrderFromBytes(data [441]byte) *Order {
	order := Order{}
	order.FromBytes(data)
	return &order
}
