package terms

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/types"
	"net/http"
	"time"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type TermsFormat struct {
	Text    string         `json:"text"`
	ID      uint           `json:"id"`
	Mask    *types.Uint256 `json:"mask"`
	MaskID  uint           `json:"maskId"`
}

type TermsSigPayload struct {
	TermsID   uint              `json:"terms_id"`
	MaskID    uint              `json:"mask_id"`
	Signature *types.Signature  `json:"sig"`
	Address   *types.Address    `json:"address"`
	Timestamp string            `json:"timestamp"`
	Nonce     string            `json:"nonce"`
}

type IngestError struct {
	Code             int               `json:"code"`
	Reason           string            `json:"reason"`
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

func TermsHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	tm := dbModule.NewTermsManager(db)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			queryObject := r.URL.Query()
			lang := queryObject.Get("lang")
			if lang == "" {
				lang = "en"
			}
			terms, err := tm.GetTerms(lang)
			if err != nil {
				returnError(w, IngestError{101, err.Error()}, 404)
				return
			}
			mask, mask_id, err := tm.GetNewHashMask(terms)
			if err != nil {
				returnError(w, IngestError{101, err.Error()}, 500)
				return
			}
			tf := &TermsFormat{
				Text: terms.Text,
				ID: terms.ID,
				Mask: &types.Uint256{},
				MaskID: mask_id,
			}
			copy(tf.Mask[32 - len(mask):], mask[:])
			data, err := json.Marshal(tf)
			if err != nil {
				returnError(w, IngestError{101, err.Error()}, 500)
				return
			}
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		} else if r.Method == "POST" {
			var data [1024]byte
			payload := &TermsSigPayload{}
			jsonLength, err := r.Body.Read(data[:])
			if err != nil && err != io.EOF {
				log.Printf(err.Error())
				returnError(w, IngestError{
					100,
					"Error reading content",
				}, 500)
				return
			}
			if err := json.Unmarshal(data[:jsonLength], &payload); err != nil {
				log.Printf("%v: '%v'", err.Error(), string(data[:]))
				returnError(w, IngestError{
					101,
					"Malformed JSON",
				}, 400)
				return
			}
			hashMask, err := tm.GetHashMaskById(payload.MaskID)
			if err != nil {
				returnError(w, IngestError{
					101,
					"Invalid HashMask ID",
				}, 400)
				return
			}
			seconds, err := strconv.Atoi(payload.Timestamp)
			if err != nil {
				returnError(w, IngestError{
					101,
					"Invalid Timestamp",
				}, 400)
				return
			}
			unixTime := time.Unix(int64(seconds), 0)
			now := time.Now()
			if unixTime.After(now) {
				returnError(w, IngestError{
					101,
					"Timestamp in the future",
				}, 400)
				return
			}
			if unixTime.Add(5 * time.Minute).Before(now) {
				returnError(w, IngestError{
					101,
					"Timestamp stale",
				}, 400)
				return
			}
			nonce, err := hex.DecodeString(strings.TrimPrefix(payload.Nonce, "0x"))
			if err != nil {
				returnError(w, IngestError{
					101,
					"Malformed Nonce",
				}, 400)
				return
			}
			ipAddress := r.Header.Get("X-Forwarded-For")
			if ipAddress == "" {
				// If we're not behind a CDN / Load Balancer, use the IP
				ipAddress = r.RemoteAddr
			}
			if err := tm.SaveSig(payload.TermsID, payload.Signature, payload.Address, payload.Timestamp, ipAddress, nonce, hashMask); err != nil {
				returnError(w, IngestError{
					101,
					err.Error(),
				}, 400)
				return
			}
			log.Printf("Saved Signature from: %v", payload.Address)
			w.WriteHeader(202)
		} else {
			returnError(w, IngestError{100, fmt.Sprintf("Unsupported Method: %v", r.Method)}, 405)
		}
	}
}

func TermsCheckHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	tm := dbModule.NewTermsManager(db)
	orderRegex := regexp.MustCompile(".*/_tos/0x([0-9a-fA-F]+)")
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			pathMatch := orderRegex.FindStringSubmatch(r.URL.Path)
			if len(pathMatch) == 0 {
				returnError(w, IngestError{100, "Malformed address"}, 404)
				return
			}
			hashHex := pathMatch[1]
			hashBytes, err := hex.DecodeString(hashHex)
			if err != nil {
				returnError(w, IngestError{100, err.Error()}, 404)
				return
			}
			address := &types.Address{}
			copy(address[:], hashBytes[:])
			if <-tm.CheckAddress(address) {
				w.WriteHeader(204)
			} else {
				w.WriteHeader(404)
			}
		} else {
			returnError(w, IngestError{100, fmt.Sprintf("Unsupported Method: %v", r.Method)}, 405)
		}
	}
}
