package ingest

import (
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
	"bytes"
)

type IngestError struct {
	Code int			`json:"code"`
	Reason string	`json:"reason"`
	ValidationErrors []ValidationError `json:"validationErrors,omitempty"`
}

type ValidationError struct {
	Field string	`json:"field"`
	Code	int	`json:"code"`
	Reason string	`json:"reason"`
}

func valInList(val [20]byte, list [][20]byte) bool {
	for _, a := range list {
		if bytes.Equal(a[:], val[:]) {
			return true
		}
	}
	return false
}

func Handler(publisher channels.Publisher, accounts accountsModule.AccountService, affiliates affiliatesModule.AffiliateService) func(http.ResponseWriter, *http.Request) {
	var contentType string
	ValidExchangeAddresses := [][20]byte{}
	// TODO: Look up valid exchanges from Redis dynamically
	addrBytes := [20]byte{}
	addr, _ := hex.DecodeString("12459c951127e0c374ff9105dda097662a027093")
	copy(addrBytes[:], addr)
	ValidExchangeAddresses = append(ValidExchangeAddresses, addrBytes)
	addr, _ = hex.DecodeString("479cc461fecd078f766ecc58533d6f69580cf3ac")
	copy(addrBytes[:], addr)
	ValidExchangeAddresses = append(ValidExchangeAddresses, addrBytes)
	addr, _ = hex.DecodeString("90fe2af704b34e0224bf2299c838e04d4dcf1364")
	copy(addrBytes[:], addr)
	ValidExchangeAddresses = append(ValidExchangeAddresses, addrBytes)
	addr, _ = hex.DecodeString("b69e673309512a9d726f87304c6984054f87a93b")
	copy(addrBytes[:], addr)
	ValidExchangeAddresses = append(ValidExchangeAddresses, addrBytes)
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
				w.WriteHeader(400)
				w.Header().Set("Content-Type", "application/json")
				errResp := IngestError{
					100,
					"Orders should be exactly 377 bytes",
					nil,
				}
				errBytes, err := json.Marshal(errResp)
				if err != nil {
					log.Printf(err.Error())
				}
				w.Write(errBytes)
				return
			} else if err != nil && err != io.EOF {
				log.Printf(err.Error())
				w.WriteHeader(500)
				w.Header().Set("Content-Type", "application/json")
				errResp := IngestError{
					100,
					"Error reading content",
					nil,
				}
				errBytes, err := json.Marshal(errResp)
				if err != nil {
					log.Printf(err.Error())
				}
				w.Write(errBytes)
				return
			}
			order.FromBytes(data)
		} else if contentType == "application/json" {
			var data [1024]byte
			jsonLength, err := r.Body.Read(data[:])
			if err != nil && err != io.EOF {
				log.Printf(err.Error())
				w.WriteHeader(500)
				w.Header().Set("Content-Type", "application/json")
				errResp := IngestError{
					100,
					"Error reading content",
					nil,
				}
				errBytes, err := json.Marshal(errResp)
				if err != nil {
					log.Printf(err.Error())
				}
				w.Write(errBytes)
				return
			}
			if err := json.Unmarshal(data[:jsonLength], &order); err != nil {
				w.WriteHeader(400)
				w.Header().Set("Content-Type", "application/json")
				log.Printf("%v: '%v'", err.Error(), string(data[:]))
				errResp := IngestError{
					101,
					"Malformed JSON",
					nil,
				}
				errBytes, err := json.Marshal(errResp)
				if err != nil {
					log.Printf(err.Error())
				}
				w.Write(errBytes)
				return
			}
		} else {
			w.WriteHeader(415)
			w.Header().Set("Content-Type", "application/json")
			errResp := IngestError{
				100,
				"Unsupported content-type",
				nil,
			}
			errBytes, err := json.Marshal(errResp)
			if err != nil {
				log.Printf(err.Error())
			}
			w.Write(errBytes)
			return
		}
		// At this point we've errored out, or we have an Order object
		emptyBytes := [20]byte{}
		if !bytes.Equal(order.Taker[:], emptyBytes[:]) {
			log.Printf("'%v' != '%v'", hex.EncodeToString(order.Taker[:]), hex.EncodeToString(emptyBytes[:]))
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			errResp := IngestError{
				100,
				"Validation Failed",
				[]ValidationError{ValidationError{
						"taker",
						1002,
						"Taker address must be empty",
					}},
			}
			errBytes, err := json.Marshal(errResp)
			if err != nil {
				log.Printf(err.Error())
			}
			w.Write(errBytes)
			return
		}
		if !valInList(order.ExchangeAddress, ValidExchangeAddresses) {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			errResp := IngestError{
				100,
				"Validation Failed",
				[]ValidationError{ValidationError{
						"exchangeContractAddress",
						1002,
						"Unknown exchangeContractAddress",
					}},
			}
			errBytes, err := json.Marshal(errResp)
			if err != nil {
				log.Printf(err.Error())
			}
			w.Write(errBytes)
			return
		}
		if !order.Signature.Verify(order.Maker) {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			errResp := IngestError{
				100,
				"Validation Failed",
				[]ValidationError{ValidationError{
						"ecSignature",
						1005,
						"Signature validation failed",
					}},
			}
			errBytes, err := json.Marshal(errResp)
			if err != nil {
				log.Printf(err.Error())
			}
			w.Write(errBytes)
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
			w.WriteHeader(402)
			w.Header().Set("Content-Type", "application/json")
			// Fee Recipient must be an authorized address
			errResp := IngestError{
				100,
				"Validation Failed",
				[]ValidationError{ValidationError{
						"feeRecipient",
						1002,
						"Invalid fee recpient",
					}},
			}
			errBytes, err := json.Marshal(errResp)
			if err != nil {
				log.Printf(err.Error())
			}
			w.Write(errBytes)
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
			w.WriteHeader(402)
			w.Header().Set("Content-Type", "application/json")
			errResp := IngestError{
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
			}
			errBytes, err := json.Marshal(errResp)
			if err != nil {
				log.Printf(err.Error())
			}
			w.Write(errBytes)
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
