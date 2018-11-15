package pool

import (
	"net/http"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/search"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	FeeShare      uint
	ID            []byte
	SenderAddress *types.Address
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
		match := poolRegex.FindStringSubmatch(r.URL.Path)
		if len(match) == 2 {
			poolName := strings.TrimPrefix(match[1], "/")
			pool :=  &Pool{}
			if q := db.Model(&Pool{}).Where("ID = ?", sha3.NewKeccak256().Sum([]byte(poolName))).First(pool); q.Error != nil {
				if poolName != "" {
					w.WriteHeader(404)
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(fmt.Sprintf("{\"code\":100,\"reason\":\"Pool Not Found: %v\"}", q.Error.Error())))
					return
				}
				// If no pool was specified and no default pool is in the database,
				// just use an empty pool
			}
			fn(w, r, pool)
		} else {
			// Routing regex shouldn't get here
			w.WriteHeader(404)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("{\"code\":100,\"reason\":\"Not Found\"}")))
			return
		}
	}
}
