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
)

func Handler(publisher channels.Publisher, accounts accountsModule.AccountService, affiliates affiliatesModule.AffiliateService) func(http.ResponseWriter, *http.Request) {
	var contentType string
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
				fmt.Fprintf(w, "{\"error\": \"Orders should be exactly 377 bytes\"}")
				return
			} else if err != nil && err != io.EOF {
				w.WriteHeader(500)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, "{\"error\": \"Error reading content\"}")
				log.Printf(err.Error())
				return
			}
			order.FromBytes(data)
		} else if contentType == "application/json" {
			var data [1024]byte
			jsonLength, err := r.Body.Read(data[:])
			if err != nil && err != io.EOF {
				w.WriteHeader(500)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, "{\"error\": \"Error reading content\"}")
				log.Printf(err.Error())
				return
			}
			if err := json.Unmarshal(data[:jsonLength], &order); err != nil {
				w.WriteHeader(400)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, "{\"error\": \"Error parsing JSON content\"}")
				log.Printf("%v: '%v'", err.Error(), string(data[:]))
				return
			}
		} else {
			w.WriteHeader(415)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "{\"error\": \"Unsupported content-type\"}")
			return
		}
		// At this point we've errored out, or we have an Order object
		if !order.Signature.Verify(order.Maker) {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "{\"error\": \"Invalid order signature\"}")
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
			fmt.Fprintf(w, "{\"error\": \"Fee Recipient must be an authorized address\"}")
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
			fmt.Fprintf(w, "{\"error\": \"makerFee + takerFee must be at least %v\"}", minFee.Text(10))
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
