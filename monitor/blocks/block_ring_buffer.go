package blocks;

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	// "log"
)

// The block monitor will store old blocks in a ring buffer to deal with chain
// reorganiations. When a reorg happens, it will uee the ring buffer to find
// the first block of the new chain and re-emit from that block forward.
type blockRingBuffer struct {
	blocks []*MiniBlock
	size   int
	head   int
}

func newBlockRingBuffer(size int) (*blockRingBuffer) {
	brb := &blockRingBuffer{}
	brb.blocks = make([]*MiniBlock, size)
	brb.size = size
	brb.head = 0
	return brb
}

func (brb *blockRingBuffer) Add(mb *MiniBlock) {
	brb.head = (brb.head + 1) % brb.size
	brb.blocks[brb.head] = mb
}

// Get returns (lastAddedItem - count). The last item added is always
// brb.Get(0), the item before that is brb.Get(1), and so on.
func (brb *blockRingBuffer) Get(count int) (*MiniBlock) {
	index := brb.head - count
	if index < 0 {
		index += brb.size
	}
	return brb.blocks[index]
}

func (brb *blockRingBuffer) HashIndex(hash common.Hash) (int) {
	for i := 0; i < brb.size; i++ {
		block := brb.Get(i)
		if block != nil && bytes.Equal(block.Hash[:], hash[:]) {
			return i
		}
		// if block != nil {
		// 	log.Printf("%#x != %#x", block.Hash[:], hash[:])
		// }
	}
	return -1
}
