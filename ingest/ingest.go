package ingest

import (
	"encoding/json"
	"fmt"
	"github.com/notegio/0xrelay/types"
	"io"
	"math/big"
	"net/http"
	"strings"
)

// Publisher items have a Publish function, allowing the publication of a
// string to a channel
type Publisher interface {
	Publish(string, string) error
}

type Account interface {
	Blacklisted() bool
	IsFeeRecipient() bool
	MinFee() *big.Int
}

type AccountService interface {
	Get([20]byte) Account
}

func Handler(publisher Publisher, accounts AccountService) func(http.ResponseWriter, *http.Request) {
	var contentType string
	return func(w http.ResponseWriter, r *http.Request) {
		if typeVal, ok := r.Header["Content-Type"]; ok {
			contentType = strings.Split(typeVal[0], ";")[0]
		} else {
			contentType = "unknown"
		}
		order := types.Order{}
		if contentType == "application/octet-stream" {
			var data [377]byte
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
				fmt.Printf(err.Error())
				return
			}
			order.FromBytes(data)
		} else if contentType == "application/json" {
			var data [1024]byte
			_, err := r.Body.Read(data[:])
			if err != nil && err != io.EOF {
				w.WriteHeader(500)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, "{\"error\": \"Error reading content\"}")
				fmt.Printf(err.Error())
				return
			}
			if err := json.Unmarshal(data[:], &order); err != nil {
				w.WriteHeader(400)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, "{\"error\": \"Error parsing JSON content\"}")
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
		makerChan := make(chan Account)
		feeChan := make(chan Account)
		go func() { feeChan <- accounts.Get(order.FeeRecipient) }()
		go func() { makerChan <- accounts.Get(order.Maker) }()
		makerFee := new(big.Int)
		takerFee := new(big.Int)
		totalFee := new(big.Int)
		makerFee.SetBytes(order.MakerFee[:])
		takerFee.SetBytes(order.TakerFee[:])
		totalFee.Add(makerFee, takerFee)
		feeRecipient := <-feeChan
		if !feeRecipient.IsFeeRecipient() {
			w.WriteHeader(402)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "{\"error\": \"Fee Recipient must be an authorized address\"}")
			return
		}
		account := <-makerChan
		minFee := account.MinFee()
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
		if err := publisher.Publish("ingest", string(orderBytes[:])); err != nil {
			fmt.Println(err)
		}
	}
}
