package fillbloom

import (
	"encoding/json"
	"math/big"
	"context"
	"github.com/willf/bloom"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/objectstorage"
	"log"
	"compress/gzip"
	"errors"
	"sync"
)

type FillBloom struct {
	b     *bloom.BloomFilter
	store objectstorage.StoredObject
	Initialized bool
	m *sync.Mutex
}

func (fb *FillBloom) Initialize(lf ethereum.LogFilterer, endBlock int64, exchangeAddresses []common.Address) error {
	itemReader, err := fb.store.Reader()
	if err != nil {
		log.Printf("FillBloom uninitialized: %v - Populating", err.Error())
		var lastBlock int64
		for lastBlock < endBlock {
			query := ethereum.FilterQuery{
				FromBlock: big.NewInt(lastBlock),
				ToBlock: big.NewInt(min(lastBlock + 100000, endBlock)),
				Addresses: exchangeAddresses,
				Topics: [][]common.Hash{
					[]common.Hash{
						common.HexToHash("0x0d0b9391970d9a25552f37d436d2aae2925e2bfe1b2a923754bada030c498cb3"),
						common.HexToHash("0x67d66f160bc93d925d05dae1794c90d2d6d6688b29b84ff069398a9b04587131"),
					},
					nil,
					nil,
				},
			}
			logs, err := lf.FilterLogs(context.Background(), query)
			if err != nil {
				return err
			}
			for _, log := range logs {
				orderHash := log.Data[len(log.Data)-32:]
				fb.Add(orderHash)
			}
			lastBlock = lastBlock + 100000
		}
	} else {
		reader, err := gzip.NewReader(itemReader)
		if err != nil { return err }
		if count, err := fb.b.ReadFrom(reader); err != nil { return err } else {
			log.Printf("Loaded bloom filter with %v bytes", count)
		}
	}
	fb.Initialized = true
	return nil
}

func (fb *FillBloom) Add(data []byte) *FillBloom {
	fb.m.Lock()
	defer fb.m.Unlock()
	fb.b.Add(data)
	return fb
}

func (fb *FillBloom) Test(data []byte) bool {
	fb.m.Lock()
	defer fb.m.Unlock()
	return fb.b.Test(data)
}

func (fb *FillBloom) Save() error {
	fb.m.Lock()
	defer fb.m.Unlock()
	b, err := fb.store.Writer()
	if err != nil { return err }
	w := gzip.NewWriter(b)
	if _, err := fb.b.WriteTo(w); err != nil { return err }
	err = w.Close()
	return err
}

func (fb *FillBloom) Consume(delivery channels.Delivery) {
	fr := &db.FillRecord{}
	if err := json.Unmarshal([]byte(delivery.Payload()), fr); err != nil {
		delivery.Reject()
	}
	orderHash := common.HexToHash(fr.OrderHash)
	fb.Add(orderHash[:])
	delivery.Ack()
}

func min(a, b int64) (int64) {
	if(a < b) {
		return a
	}
	return b
}

func NewFillBloom(storedURI string) (*FillBloom, error) {

	storedObject := objectstorage.GetStoredObject(storedURI)
	if storedObject == nil {
		return nil, errors.New("Invalid stored object URI")
	}

	return &FillBloom{
		bloom.New(419430400, 17),
		storedObject,
		false,
		&sync.Mutex{},
	}, nil
}
