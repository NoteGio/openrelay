package pool

import (
	"context"
	"net/http"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/search"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/gorm"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	urlModule "net/url"
	"regexp"
	"strings"
	"fmt"
)

type Pool struct {
	SearchTerms   string
	Expiration    uint
	Nonce         uint
	ValidateLive  bool
	FeeShare      uint
	ID            []byte
	FilterAddress *types.Address
	conn          bind.ContractBackend
}

func (pool *Pool) SetConn(conn bind.ContractBackend) {
	pool.conn = conn
}

func (pool Pool) Filter(query *gorm.DB) (*gorm.DB, error) {
	queryObject, err := urlModule.ParseQuery(pool.SearchTerms)
	if err != nil {
		return nil, err
	}
	query, errs := search.QueryFilter(query, queryObject)
	if errCount := len(errs); errCount > 0 {
		return nil, fmt.Errorf("Found %v errors in pool query string", errCount)
	}
	return query, nil
}

var poolRegex = regexp.MustCompile("^(/[^/]*)?/v2/")

func PoolDecorator(db *gorm.DB, fn func(http.ResponseWriter, *http.Request, *Pool)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		poolName := "default"
		match := poolRegex.FindStringSubmatch(r.URL.Path)
		if len(match) == 2 {
			poolName = strings.TrimPrefix(match[1], "/")
		}
		pool :=  &Pool{}
		if q := db.Model(&Pool{}).Where("ID = ?", sha3.NewKeccak256().Sum([]byte(poolName))).First(pool); q.Error != nil {
			w.WriteHeader(404)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("{\"code\":100,\"reason\":\"Pool Not Found\"}")))
		} else {
			fn(w, r, pool)
		}
	}
}

func PoolDecoratorConn(db *gorm.DB, rpcURL string, fn func(http.ResponseWriter, *http.Request, *Pool)) (func(http.ResponseWriter, *http.Request), error) {
	conn, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	if _, err = conn.SyncProgress(context.Background()); err != nil {
		// This is just here so that an NewRpcFeeToken can't be instantiated
		// successfully if the RPC server isn't responding properly. What RPC
		// function we call isn't important, but SyncProgress is pretty cheap.
		return nil, err
	}
	wfn := func (w http.ResponseWriter, r *http.Request, pool *Pool) {
		pool.SetConn(conn)
		fn(w, r, pool)
	}
	return PoolDecorator(db, wfn), nil
}
