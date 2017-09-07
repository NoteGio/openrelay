package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"strconv"
)

// Order represents an 0x order object
type Order struct {
	Maker                    [20]byte
	Taker                    [20]byte
	MakerToken               [20]byte
	TakerToken               [20]byte
	FeeRecipient             [20]byte
	ExchangeAddress          [20]byte
	MakerTokenAmount         [32]byte
	TakerTokenAmount         [32]byte
	MakerFee                 [32]byte
	TakerFee                 [32]byte
	ExpirationTimestampInSec [32]byte
	Salt                     [32]byte
	Signature                *Signature
}

// NewOrder takes string representations of values and converts them into an Order object
func NewOrder(maker, taker, makerToken, takerToken, feeRecipient, exchangeAddress, makerTokenAmount, takerTokenAmount, makerFee, takerFee, expirationTimestampInSec, salt, sigV, sigR, sigS string) (*Order, error) {
	order := Order{}
	if err := order.fromStrings(maker, taker, makerToken, takerToken, feeRecipient, exchangeAddress, makerTokenAmount, takerTokenAmount, makerFee, takerFee, expirationTimestampInSec, salt, sigV, sigR, sigS); err != nil {
		return nil, err
	}
	return &order, nil
}

func (order *Order) fromStrings(maker, taker, makerToken, takerToken, feeRecipient, exchangeAddress, makerTokenAmount, takerTokenAmount, makerFee, takerFee, expirationTimestampInSec, salt, sigV, sigR, sigS string) error {
	makerBytes, err := hexStringToBytes(maker)
	if err != nil {
		return err
	}
	takerBytes, err := hexStringToBytes(taker)
	if err != nil {
		return err
	}
	makerTokenBytes, err := hexStringToBytes(makerToken)
	if err != nil {
		return err
	}
	takerTokenBytes, err := hexStringToBytes(takerToken)
	if err != nil {
		return err
	}
	feeRecipientBytes, err := hexStringToBytes(feeRecipient)
	if err != nil {
		return err
	}
	exchangeAddressBytes, err := hexStringToBytes(exchangeAddress)
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
	sigRBytes, err := hexStringToBytes(sigR)
	if err != nil {
		return err
	}
	sigSBytes, err := hexStringToBytes(sigS)
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
	Maker                    string        `json:"maker"`
	Taker                    string        `json:"taker"`
	MakerToken               string        `json:"makerTokenAddress"`
	TakerToken               string        `json:"takerTokenAddress"`
	FeeRecipient             string        `json:"feeRecipient"`
	ExchangeAddress          string        `json:"exchangeContractAddress"`
	MakerTokenAmount         string        `json:"makerTokenAmount"`
	TakerTokenAmount         string        `json:"takerTokenAmount"`
	MakerFee                 string        `json:"makerFee"`
	TakerFee                 string        `json:"takerFee"`
	ExpirationTimestampInSec string        `json:"expirationUnixTimestampSec"`
	Salt                     string        `json:"salt"`
	Signature                jsonSignature `json:"ecSignature"`
}

func (order *Order) UnmarshalJSON(b []byte) error {
	jOrder := jsonOrder{}
	if err := json.Unmarshal(b, &jOrder); err != nil {
		return err
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
	)

	return nil
}

func (order *Order) Bytes() [377]byte {
	var output [377]byte
	copy(output[0:20], order.Maker[:])                       // 20
	copy(output[20:40], order.Taker[:])                      // 20
	copy(output[40:60], order.MakerToken[:])                 // 20
	copy(output[60:80], order.TakerToken[:])                 // 20
	copy(output[80:100], order.FeeRecipient[:])              // 20
	copy(output[100:120], order.ExchangeAddress[:])          // 20
	copy(output[120:152], order.MakerTokenAmount[:])         // 32
	copy(output[152:184], order.TakerTokenAmount[:])         // 32
	copy(output[184:216], order.MakerFee[:])                 // 32
	copy(output[216:248], order.TakerFee[:])                 // 32
	copy(output[248:280], order.ExpirationTimestampInSec[:]) // 32
	copy(output[280:312], order.Salt[:])                     // 32
	output[312] = order.Signature.V
	copy(output[313:345], order.Signature.R[:])
	copy(output[345:377], order.Signature.S[:])
	return output
}

func (order *Order) FromBytes(data [377]byte) {
	copy(order.Maker[:], data[0:20])
	copy(order.Taker[:], data[20:40])
	copy(order.MakerToken[:], data[40:60])
	copy(order.TakerToken[:], data[60:80])
	copy(order.FeeRecipient[:], data[80:100])
	copy(order.ExchangeAddress[:], data[100:120])
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
}

func OrderFromBytes(data [377]byte) *Order {
	order := Order{}
	order.FromBytes(data)
	return &order
}
