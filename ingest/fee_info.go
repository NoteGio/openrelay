package ingest

import (
	"encoding/json"
	accountsModule "github.com/notegio/openrelay/accounts"
	affiliatesModule "github.com/notegio/openrelay/affiliates"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	poolModule "github.com/notegio/openrelay/pool"
	"io"
	"log"
	"math/big"
	"net/http"
	"fmt"
)

// FeeInputPayload only considers maker and feeRecipient when calculating fees.
// Everything else will be ignored.
type FeeInputPayload struct {
	ChainID      uint `json:"chainId"`
	Maker        string `json:"maker"`
	Exchange     string `json:"exchangeAddress"`
	FeeRecipient string `json:"feeRecipientAddress"`
	Taker        string `json:"taker"`
	Sender       string `json:"senderAddress"`
}

type FeeResponse struct {
	MakerFee       string `json:"makerFee"`
	TakerFee       string `json:"takerFee"`
	FeeRecipient   string `json:"feeRecipientAddress"`
	Sender         string `json:"senderAddress"`
	TakerToSpecify string `json:"takerToSpecify"`
	MakerFeeAssetData string `json:"makerFeeAssetData"`
	TakerFeeAssetData string `json:"takerFeeAssetData"`
}

func FeeHandler(publisher channels.Publisher, accounts accountsModule.AccountService, affiliates affiliatesModule.AffiliateService, defaultFeeRecipient [20]byte, exchangeLookup ExchangeLookup) func(http.ResponseWriter, *http.Request, *poolModule.Pool) {
	emptyBytes := &types.Address{}
	return func(w http.ResponseWriter, r *http.Request, pool *poolModule.Pool) {
		var data [1024]byte
		feeInput := &FeeInputPayload{}

		jsonLength, err := r.Body.Read(data[:])
		if err != nil && err != io.EOF {
			log.Printf(err.Error())
			returnError(w, IngestError{
				100,
				"Error reading content",
				nil,
			}, 500)
			return
		}
		if err := json.Unmarshal(data[:jsonLength], &feeInput); err != nil {
			log.Printf("%v: '%v'", err.Error(), string(data[:]))
			returnError(w, IngestError{
				101,
				"Malformed JSON",
				nil,
			}, 400)
			return
		}
		makerSlice, err := types.HexStringToBytes(feeInput.Maker)
		if err != nil && feeInput.Maker != "" {
			log.Printf("%v: '%v'", err.Error(), string(data[:]))
			returnError(w, IngestError{
				100,
				"Validation failed",
				[]ValidationError{ValidationError{
					"maker",
					1001,
					"Invalid address format",
				},
				},
			}, 400)
			return
		}
		feeRecipientAddressSlice, err := types.HexStringToBytes(feeInput.FeeRecipient)
		if err != nil && feeInput.FeeRecipient != "" {
			log.Printf("%v: '%v'", err.Error(), string(data[:]))
			returnError(w, IngestError{
				100,
				"Validation failed",
				[]ValidationError{ValidationError{
					"feeRecipient",
					1001,
					"Invalid address format",
				},
				},
			}, 400)
			return
		}
		exchangeAddressSlice, err := types.HexStringToBytes(feeInput.Exchange)
		if err != nil && feeInput.FeeRecipient != "" {
			log.Printf("%v: '%v'", err.Error(), string(data[:]))
			returnError(w, IngestError{
				100,
				"Validation failed",
				[]ValidationError{ValidationError{
					"exchangeAddress",
					1001,
					"Invalid address format",
				},
				},
			}, 400)
			return
		}
		exchangeAddress := &types.Address{}
		copy(exchangeAddress[:], exchangeAddressSlice[:])
		networkIDChan := exchangeLookup.ExchangeIsKnown(exchangeAddress)
		makerAddress := &types.Address{}
		copy(makerAddress[:], makerSlice[:])
		feeRecipientAddress := &types.Address{}
		if feeInput.FeeRecipient == "" {
			copy(feeRecipientAddress[:], defaultFeeRecipient[:])
		} else {
			copy(feeRecipientAddress[:], feeRecipientAddressSlice)
		}
		makerChan := make(chan accountsModule.Account)
		affiliateChan := make(chan affiliatesModule.Affiliate)
		go func() {
			feeRecipient, err := affiliates.Get(feeRecipientAddress)
			if err != nil {
				log.Printf("Affiliate error: %v", err.Error())
				affiliateChan <- nil
			} else {
				affiliateChan <- feeRecipient
			}
		}()
		go func() { makerChan <- accounts.Get(makerAddress) }()
		feeRecipient := <-affiliateChan
		if feeRecipient == nil {
			returnError(w, IngestError{
				100,
				"Validation Failed",
				[]ValidationError{ValidationError{
					"feeRecipient",
					1002,
					"Invalid fee recipient",
				}},
			}, 402)
			return
		}
		poolFee, err := pool.Fee()
		if err != nil {
			log.Printf("Pool fee error: %v", err.Error())
			returnError(w, IngestError{
				100,
				"Validation Failed",
				[]ValidationError{ValidationError{
					"pool",
					1002,
					"Pool error",
				}},
			}, 500)
			return
		}
		chainId := feeInput.ChainID
		if chainId == 0 {
			chainId = 1
		}
		feeTokenAssetData, err := pool.FeeAssetData(chainId)
		if err != nil {
			log.Printf("Pool fee token error: %v", err.Error())
			returnError(w, IngestError{
				100,
				"Validation Failed",
				[]ValidationError{ValidationError{
					"pool",
					1002,
					"Pool error",
				}},
			}, 500)
			return
		}
		account := <-makerChan
		minFee := new(big.Int)

		// A fee recipient's Fee() value is the base fee for that recipient. A
		// maker's Discount() is the discount that recipient gets from the base
		// fee. Thus, the minimum fee required is pool.Fee() - maker.Discount()
		minFee.Sub(poolFee, account.Discount())
		takerToSpecify := fmt.Sprintf("%#x", emptyBytes[:])
		networkID := <-networkIDChan
		if networkID == 0 {
			networkID = 1
		}
		var senderToSpecify string
		senderAddress, ok := pool.SenderAddresses[networkID]
		if ok {
			senderToSpecify = fmt.Sprintf("%#x", senderAddress[:])
		} else {
			senderToSpecify = fmt.Sprintf("%#x", emptyBytes[:])

		}
		if feeInput.Taker != "" {
			takerToSpecify = feeInput.Taker
		}
		if feeInput.Sender != "" {
			senderToSpecify = feeInput.Sender
		}
		feeResponse := &FeeResponse{
			minFee.Text(10),
			"0",
			fmt.Sprintf("%#x", feeRecipientAddress[:]),
			senderToSpecify,
			takerToSpecify,
			fmt.Sprintf("%#x", feeTokenAssetData),
			fmt.Sprintf("%#x", feeTokenAssetData),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		feeBytes, err := json.Marshal(feeResponse)
		if err != nil {
			log.Printf(err.Error())
		}
		w.Write(feeBytes)
		return
	}
}
