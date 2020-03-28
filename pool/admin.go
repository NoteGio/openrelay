package pool

import (
	"github.com/notegio/openrelay/channels"
	dbModule "github.com/notegio/openrelay/db"
	"strings"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"time"
)

func handleError(w http.ResponseWriter, text string, code int) {
  w.WriteHeader(code)
  w.Write([]byte(fmt.Sprintf(`{"ok": false, "error": {"message": "%v"}}`, text)))
}

type call struct {
	Method string `json:"method"`
	Params []interface{} `json:"params"`
	ID interface{} `json:"id"`
	Expiration int64 `json:"expiration"`
}

type response struct {
	ID interface{} `json:"id"`
	Result interface{} `json:"result"`
}

func handleResponse(w http.ResponseWriter, result, id interface{}, code int) {
	data, _ := json.Marshal(response{ID: id, Result: result})
	w.WriteHeader(code)
	w.Write(data)

}


// /poolName/v3/_admin
// POST
// Body: {"jsonrpc": "2.0", "id": STRING|INT, "expiration": INT, "method": STRING, "params": LIST}
// * id: An identifier for you to correlate requests and responses
// * expiration: Unix timestamp in seconds. Must be in the range (now < expiration <= now + 10)
// * method: The RPC method to invoke
// * params: The parameters taken by the specified RPC method.
// Header: Authorization: HMAC_sha256(key, body)
func PoolAdminHandler(db *gorm.DB, cancellationsPublisher channels.Publisher) func(http.ResponseWriter, *http.Request, *Pool) {
	return func(w http.ResponseWriter, r *http.Request, pool *Pool) {
		if len(pool.Key()) == 0 {
			handleError(w, "invalid signature", 401)
		}
		h := hmac.New(sha256.New, pool.Key())
		data, err := ioutil.ReadAll(r.Body)
    if err != nil {
      handleError(w, "error reading body", 400)
      return
    }
		checksum := h.Sum(data)
		sigHex, ok := r.Header["Authorization"]
		if !ok {
			handleError(w, "not authorized", 401)
			return
		}
		sig, err := hex.DecodeString(strings.TrimPrefix(string(sigHex[0]), "0x"))
		if err != nil {
			handleError(w, err.Error(), 401)
			return
		}
		if !hmac.Equal(checksum, sig) {
			handleError(w, "invalid signature", 401)
			return
		}
		c := &call{}
		if err := json.Unmarshal(data, c); err != nil {
			handleError(w, "error parsing request", 400)
			return
		}
		expiration := time.Unix(c.Expiration, 0)
		if expiration.Before(time.Now()) {
			handleError(w, "request expired", 400)
		}
		if expiration.After(time.Now().Add(time.Minute)) {
			handleError(w, "requests should expire no more than 60 seconds in the future", 400)
		}

		if c.Method == "cancellation" {
			result, err, status := cancellation(db, cancellationsPublisher, pool, c.Params)
			if err != nil {
				handleError(w, err.Error(), status)
				return
			}
			handleResponse(w, result, c.ID, 200)
		} else {
			handleError(w, "unknown rpc method", 400)
		}
	}
}


// RPC Method: cancellation
// Params: List of order hashes to cancel
// Orders will only be cancelled if they were submitted through this pool.
// Return value: Number of orders to be canclled. This may not match the length
// of the list, if some of the items in the list do not correspond to orders
// submitted through the pool
func cancellation(db *gorm.DB, cancellationsPublisher channels.Publisher, pool *Pool, params []interface{}) (interface{}, error, int) {
	orderList := [][]byte{}
	for _, p := range params {
		orderHash, ok := p.(string)
		if !ok { return nil, fmt.Errorf("All parameters must be strings"), 400 }
		data, err := hex.DecodeString(strings.TrimPrefix(orderHash, "0x"))
		if err != nil { return nil, err, 400}
		orderList = append(orderList, data)
	}
	dbOrders := []dbModule.Order{}
	where := []string{}
	for _ = range orderList {
		where = append(where, "order_hash = ?")
	}

	parameters := append([]interface{}{pool.ID})
	for _, order := range orderList {
		parameters = append(parameters, order)
	}
	err := db.Model(&dbModule.Order{}).Where(fmt.Sprintf("pool_id = ? AND (%v)", strings.Join(where, " OR ")), parameters...).Find(&dbOrders).Error
	if err != nil {
		return nil, err, 500
	}
	for _, order := range dbOrders {
		fr := dbModule.FillRecord{
			OrderHash: fmt.Sprintf("%#x", order.Hash()),
			Cancel: true,
		}
		data, err := json.Marshal(fr)
		if err != nil { return nil, err, 500}
		cancellationsPublisher.Publish(string(data))
	}
	return len(dbOrders), nil, 200
}
