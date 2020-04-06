package blocks

import (
	"context"
	"encoding/json"
	"time"
	"errors"
	"math/big"
	"github.com/ethereum/go-ethereum"
	// "github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/notegio/openrelay/channels"
	"log"
)

// MiniBlock is a subset of the Ethereum block header that has the subset of
// fields we need to monitor for events. The hash is tracked to identify
// specific blocks in the event of a reorg. The block number is tracked to make
// it easy to guess what block should come next (though reorgs can alter this).
// The bloom filter is tracked so that downstream tasks can efficiently
// determine if they need to take action on this block.
type MiniBlock struct {
	Hash   common.Hash  `json:"hash"`
	Number *big.Int     `json:"number"`
	Bloom  types.Bloom `json:"bloom"`
}

// HeaderGetter returns block headers by hash or number. The ethclient provides
// this interface, but for test purposes we want a simpler interface.
type HeaderGetter interface {
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
	HeaderByHash(context.Context, common.Hash) (*types.Header, error)
}

// BlockRecorder keeps track of the last recorded block, primarily so that the
// block monitor can resume where it left off in the event that it restarts.
type BlockRecorder interface {
	Record(*big.Int) (error)
	Get() (*big.Int, error)
}

// BlockMonitor watches a HeaderGetter (probably an ethclient)for new blocks,
// publishing new blocks to a Publisher. In the event of a chain
// reorganization, it will emit any blocks from the new chain, so long as the
// common ancestor is in its block ring buffer. If the HeaderGetter does not
// yet have the next block, the BlockMonitor will poll every queryInterval.
// Finally, a BlockRecorder is used to track the last recorded block number,
// so that the BlockMonitor can resume where it left off in the event of a
// restart.
type BlockMonitor struct {
	brb           *blockRingBuffer
	headerGetter  HeaderGetter
	publisher     channels.Publisher
	queryInterval time.Duration
	blockRecorder BlockRecorder
	quit          chan bool
}

// Process watches for new blocks, publishing each block on the provided
// publisher.
func (bm *BlockMonitor) Process() error {
	blockNumber, err := bm.blockRecorder.Get()
	// If we get an error retrieving the last block, log it, but continue. A nil
	// blockNumber will retrieve the latest headers.
	if err != nil {
		log.Printf("Error getting block number: %v", err.Error())
	}
	header, err := bm.headerGetter.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Printf("Error getting header for block number %v", blockNumber)
		return err
	}
	log.Printf("Starting Block: Number: '%v' - Hash: '%#x' - Pause Time: '%v'", header.Number, header.Hash(), bm.queryInterval)
	// Track the block in the RingBuffer to handle chain re-orgs.
	bm.brb.Add(&MiniBlock{
		header.Hash(),
		header.Number,
		header.Bloom,
	})
	// Only publish the initial block if blocknumber == 0. For later blocks, we
	// should have published the block in an earlier iteration, so we don't need
	// to publish it now.
	if header.Number.Int64() == 0 {
		if err := bm.publish(bm.brb.Get(0)); err != nil {
			log.Printf("Error publishing block")
			return err
		}
	}
	// data, _ := rlp.EncodeToBytes(header)
	// log.Printf("Initial rlp: '%x'", data)
	MAIN_PROCESSING:
	for {
		select {
		case _ = <-bm.quit:
			// Return on the quite signal.
			return nil
		default:
		}
		// Ask the headerGetter for the last known block + 1.
		header, err = bm.headerGetter.HeaderByNumber(context.Background(), new(big.Int).Add(bm.brb.Get(0).Number, big.NewInt(1)))
		if err == ethereum.NotFound {
			// If no block is available, sleep for a bit and try again.
			time.Sleep(bm.queryInterval)
			continue
		} else if err != nil {
			// If we got an unexpected error, return it.
			log.Printf("error getting header for block %v", new(big.Int).Add(bm.brb.Get(0).Number, big.NewInt(1)))
			return err
		}
		// data, _ = rlp.EncodeToBytes(header)
		// log.Printf("Block rlp: '%x'", data)
		// In the event of a chain reorg, the current block's parent won't be
		// present in our block ring buffer. We need to follow the block's parents
		// backwards until we find a recognized ancestor, or until we've exhausted
		// our ring buffer.
		counter := 0
		// TODO: Somehow this is adding block 0, then trying to add block 0 again
		for bm.brb.HashIndex(header.ParentHash) == -1 && counter < bm.brb.size {
			if header.Number.Int64() == 0 {
				break
			}
			log.Printf("Getting parent for header: %v - %#x (Probable chain reorg)", header.Number.Int64(), header.ParentHash)
			counter++
			parentHash := header.ParentHash
			header, err = bm.headerGetter.HeaderByHash(context.Background(), header.ParentHash)
			if err == ethereum.NotFound {
				log.Printf("HeaderGetter is missing parent hash: %#x", header.ParentHash)
				time.Sleep(bm.queryInterval)
				continue MAIN_PROCESSING
			} else if err != nil {
				log.Printf("error getting header for hash %#x", parentHash)
				return err
			}
		}
		if bm.brb.HashIndex(header.Hash()) != -1 {
			log.Fatalf("No parents found, but current block already exists. It's likely that block.Hash() is not being computed properly somewhere.")
		}
		// At this point we either have the next header, we have wound back to the
		// beginning of a reorg, or we've wound back as far as we can given our
		// ring buffer size
		bm.brb.Add(&MiniBlock{
			header.Hash(),
			header.Number,
			header.Bloom,
		})
		log.Printf("Published Block %v - %#x", bm.brb.Get(0).Number, bm.brb.Get(0).Hash)
		if err := bm.publish(bm.brb.Get(0)); err != nil {
			return err
		}
	}
}

// publish sends a JSON marshalled miniblock to the publisher, and records the
// block number in the blockRecorder.
func (bm *BlockMonitor) publish(block *MiniBlock) error {
	data, err := json.Marshal(block)
	if err != nil {
		return err
	}
	result := bm.publisher.Publish(string(data))
	if result {
		return bm.blockRecorder.Record(block.Number)
	} else {
		return errors.New("Failed to publish block")
	}
}

// Stop sends the signal to stop processing.
func (bm *BlockMonitor) Stop() {
	bm.quit <- true
}

// NewBlockMonitor creates an BlockMonitor with the provided HeaderGetter,
// Publisher, and BlockRecorder. When blocks are unavailable, it will sleep for
// `interval` between polling, and it will keep `brbSize` historic block
// headers to deal with chain reorganizations.
func NewBlockMonitor(headerGetter HeaderGetter, publisher channels.Publisher, interval time.Duration, blockRecorder BlockRecorder, brbSize int) (*BlockMonitor) {
	return &BlockMonitor{
		newBlockRingBuffer(brbSize),
		headerGetter,
		publisher,
		interval,
		blockRecorder,
		make(chan bool),
	}
}

// NewRPCBlockMonitor creates a BlockMonitor using an ehtclient to the
// specified rpcURL for a HeaderGetter. Other parameters match NewBlockMonitor.
func NewRPCBlockMonitor(rpcURL string, publisher channels.Publisher, interval time.Duration, blockRecorder BlockRecorder, brbSize int) (*BlockMonitor, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return NewBlockMonitor(client, publisher, interval, blockRecorder, brbSize), nil
}
