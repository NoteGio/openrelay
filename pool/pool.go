package pool

import (
	"net/http"
	"github.com/notegio/openrelay/config"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/search"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/jinzhu/gorm"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"gopkg.in/redis.v3"
	"math/big"
	urlModule "net/url"
	"regexp"
	"strings"
	"fmt"
)

var feeBaseUnits = big.NewInt(1000000000000000000)

type Pool struct {
	SearchTerms   string
	Expiration    uint
	Nonce         uint
	FeeShare      string
	ID            []byte
	SenderAddress *types.Address
	FilterAddress *types.Address
	conn          bind.ContractBackend
	baseFee       config.BaseFee
}

func (pool *Pool) SetConn(conn bind.ContractBackend) {
	pool.conn = conn
}

func (pool *Pool) SetBaseFee(baseFee config.BaseFee) {
	pool.baseFee = baseFee
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

func (pool Pool) Fee() (*big.Int, error) {
	baseFee, err := pool.baseFee.Get()
	if err != nil {
		return nil, err
	}
	feeShare := new(big.Int)
	if _, ok := feeShare.SetString(pool.FeeShare, 10); !ok {
		// If the fee share is not a valid integer, just return the base fee
		return baseFee, nil
	}
	combined := new(big.Int).Mul(baseFee, feeShare)
	return new(big.Int).Div(combined, feeBaseUnits), nil
}

var poolRegex = regexp.MustCompile("^(/[^/]*)?/v2/")

func PoolDecorator(db *gorm.DB, fn func(http.ResponseWriter, *http.Request, types.Pool)) func(http.ResponseWriter, *http.Request) {
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

func PoolDecoratorBaseFee(db *gorm.DB, redisClient *redis.Client, fn func(http.ResponseWriter, *http.Request, *Pool)) func(http.ResponseWriter, *http.Request) {
	baseFee := config.NewBaseFee(redisClient)
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
			pool.baseFee = baseFee
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
