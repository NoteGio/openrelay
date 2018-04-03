// The block recorder needs to keep track of the last recorded block so that
// the block monitor can resume in the event that it gets restarted.

package blocks

import (
	"math/big"
	"gopkg.in/redis.v3"
	"errors"
)

type redisBlockRecorder struct {
	redisClient *redis.Client
	key         string
}

func (br *redisBlockRecorder) Record(blockNumber *big.Int) (error) {
	return br.redisClient.Set(br.key, blockNumber.String(), 0).Err()
}

func (br *redisBlockRecorder) Get() (*big.Int, error) {
	result := br.redisClient.Get(br.key)
	if err := result.Err(); err != nil {
		return nil, err
	}
	intResult, err := result.Int64()
	if err != nil {
		return nil, err
	}
	bigInt := big.NewInt(intResult)
	return bigInt, nil
}

func NewRedisBlockRecorder(redisClient *redis.Client, key string) BlockRecorder {
	return &redisBlockRecorder{redisClient, key}
}

type mockBlockRecorder struct {
	blockNumber *big.Int
}

func (br *mockBlockRecorder) Record(blockNumber *big.Int) (error) {
	br.blockNumber = blockNumber
	return nil
}

func (br *mockBlockRecorder) Get() (*big.Int, error) {
	if br.blockNumber != nil {
		return br.blockNumber, nil
	}
	return nil, errors.New("No block set")
}

func NewMockBlockRecorder() BlockRecorder {
	return &mockBlockRecorder{}
}
