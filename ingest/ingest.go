package ingest

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	accountsModule "github.com/notegio/openrelay/accounts"
	affiliatesModule "github.com/notegio/openrelay/affiliates"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
)

type IngestError struct {
	Code             int               `json:"code"`
	Reason           string            `json:"reason"`
	ValidationErrors []ValidationError `json:"validationErrors,omitempty"`
}

type ValidationError struct {
	Field  string `json:"field"`
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

func valInList(val *types.Address, list []types.Address) bool {
	for _, a := range list {
		if bytes.Equal(a[:], val[:]) {
			return true
		}
	}
	return false
}

func returnError(w http.ResponseWriter, errResp IngestError, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	errBytes, err := json.Marshal(errResp)
	if err != nil {
		log.Printf(err.Error())
	}
	w.Write(errBytes)
}

func Handler(publisher channels.Publisher, accounts accountsModule.AccountService, affiliates affiliatesModule.AffiliateService) func(http.ResponseWriter, *http.Request) {
	var contentType string
	ValidExchangeAddresses := []types.Address{}
	// TODO: Look up valid exchanges from Redis dynamically
	addrBytes := &types.Address{}
	addr, _ := hex.DecodeString("12459c951127e0c374ff9105dda097662a027093")
	copy(addrBytes[:], addr)
	ValidExchangeAddresses = append(ValidExchangeAddresses, *addrBytes)
	addr, _ = hex.DecodeString("479cc461fecd078f766ecc58533d6f69580cf3ac")
	copy(addrBytes[:], addr)
	ValidExchangeAddresses = append(ValidExchangeAddresses, *addrBytes)
	addr, _ = hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364")
	copy(addrBytes[:], addr)
	ValidExchangeAddresses = append(ValidExchangeAddresses, *addrBytes)
	addr, _ = hex.DecodeString("b69e673309512a9d726f87304c6984054f87a93b")
	copy(addrBytes[:], addr)
	ValidExchangeAddresses = append(ValidExchangeAddresses, *addrBytes)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Health checks
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "{\"ok\": true}")
			return
		}
		if typeVal, ok := r.Header["Content-Type"]; ok {
			contentType = strings.Split(typeVal[0], ";")[0]
		} else {
			contentType = "unknown"
		}
		order := types.Order{}
		if contentType == "application/octet-stream" {
			var data [441]byte
			length, err := r.Body.Read(data[:])
			if length != 377 {
				returnError(w, IngestError{
					100,
					"Orders should be exactly 377 bytes",
					nil,
				}, 400)
				return
			} else if err != nil && err != io.EOF {
				log.Printf(err.Error())
				returnError(w, IngestError{
					100,
					"Error reading content",
					nil,
				}, 500)
				return
			}
			order.FromBytes(data)
		} else if contentType == "application/json" {
			var data [1024]byte
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
			if err := json.Unmarshal(data[:jsonLength], &order); err != nil {
				log.Printf("%v: '%v'", err.Error(), string(data[:]))
				returnError(w, IngestError{
					101,
					"Malformed JSON",
					nil,
				}, 400)
				return
			}
		} else {
			returnError(w, IngestError{
				100,
				"Unsupported content-type",
				nil,
			}, 415)
			return
		}
		// At this point we've errored out, or we have an Order object
		emptyBytes := [20]byte{}
		if !bytes.Equal(order.Taker[:], emptyBytes[:]) {
			log.Printf("'%v' != '%v'", hex.EncodeToString(order.Taker[:]), hex.EncodeToString(emptyBytes[:]))
			returnError(w, IngestError{
				100,
				"Validation Failed",
				[]ValidationError{ValidationError{
					"taker",
					1002,
					"Taker address must be empty",
				}},
			}, 400)
			return
		}
		if !valInList(order.ExchangeAddress, ValidExchangeAddresses) {
			returnError(w, IngestError{
				100,
				"Validation Failed",
				[]ValidationError{ValidationError{
					"exchangeContractAddress",
					1002,
					"Unknown exchangeContractAddress",
				}},
			}, 400)
			return
		}
		if !order.Signature.Verify(order.Maker) {
			returnError(w, IngestError{
				100,
				"Validation Failed",
				[]ValidationError{ValidationError{
					"ecSignature",
					1005,
					"Signature validation failed",
				}},
			}, 400)
			return
		}
		// Now that we have a complete order, request the account from redis
		// asynchronously since this may have some latency
		makerChan := make(chan accountsModule.Account)
		affiliateChan := make(chan affiliatesModule.Affiliate)
		go func() {
			feeRecipient, err := affiliates.Get(order.FeeRecipient)
			if err != nil {
				affiliateChan <- nil
			} else {
				affiliateChan <- feeRecipient
			}
		}()
		go func() { makerChan <- accounts.Get(order.Maker) }()
		makerFee := new(big.Int)
		takerFee := new(big.Int)
		totalFee := new(big.Int)
		makerFee.SetBytes(order.MakerFee[:])
		takerFee.SetBytes(order.TakerFee[:])
		totalFee.Add(makerFee, takerFee)
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
		account := <-makerChan
		minFee := new(big.Int)
		// A fee recipient's Fee() value is the base fee for that recipient. A
		// maker's Discount() is the discount that recipient gets from the base
		// fee. Thus, the minimum fee required is feeRecipient.Fee() -
		// maker.Discount()
		minFee.Sub(feeRecipient.Fee(), account.Discount())
		if totalFee.Cmp(minFee) < 0 {
			returnError(w, IngestError{
				100,
				"Validation Failed",
				[]ValidationError{ValidationError{
					"makerFee",
					1004,
					"Total fee must be at least: " + minFee.Text(10),
				},
					ValidationError{
						"takerFee",
						1004,
						"Total fee must be at least: " + minFee.Text(10),
					},
				},
			}, 402)
			return
		}
		if account.Blacklisted() {
			w.WriteHeader(202)
			fmt.Fprintf(w, "")
			return
		}
		w.WriteHeader(202)
		fmt.Fprintf(w, "")
		orderBytes := order.Bytes()
		if err := publisher.Publish(string(orderBytes[:])); !err {
			log.Println("Failed to publish '%v'", hex.EncodeToString(order.Hash()))
		}
	}
}
