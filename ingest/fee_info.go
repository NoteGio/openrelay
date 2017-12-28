package ingest

import (
	"encoding/hex"
	"encoding/json"
	accountsModule "github.com/notegio/openrelay/accounts"
	affiliatesModule "github.com/notegio/openrelay/affiliates"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	"io"
	"log"
	"math/big"
	"net/http"
)

// FeeInputPayload only considers maker and feeRecipient when calculating fees.
// Everything else will be ignored.
type FeeInputPayload struct {
	Maker        string `json:"maker"`
	FeeRecipient string `json:"feeRecipient"`
}

type FeeResponse struct {
	MakerFee       string `json:"makerFee"`
	TakerFee       string `json:"takerFee"`
	FeeRecipient   string `json:"feeRecipient"`
	TakerToSpecify string `json:"takerToSpecify"`
}

func FeeHandler(publisher channels.Publisher, accounts accountsModule.AccountService, affiliates affiliatesModule.AffiliateService, defaultFeeRecipient [20]byte) func(http.ResponseWriter, *http.Request) {
	emptyBytes := &types.Address{}
	return func(w http.ResponseWriter, r *http.Request) {
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
					"Invalid fee recpient",
				}},
			}, 402)
			return
		}
		account := <-makerChan
		minFee := new(big.Int)
		// A fee recipient's Fee() value is the base fee for that recipient. A
		// maker's Discount() is the discount that recipient gets from the base
		// fee. Thus, the minimum fee required is feeRecipient.Fee() -
		// maker.Discount()
		minFee.Sub(feeRecipient.Fee(), account.Discount())
		feeResponse := &FeeResponse{
			minFee.Text(10),
			"0",
			"0x" + hex.EncodeToString(feeRecipientAddress[:]),
			"0x" + hex.EncodeToString(emptyBytes[:]),
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		feeBytes, err := json.Marshal(feeResponse)
		if err != nil {
			log.Printf(err.Error())
		}
		w.Write(feeBytes)
		return
	}
}
