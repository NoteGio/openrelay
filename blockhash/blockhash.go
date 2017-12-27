package blockhash

import (
	"github.com/notegio/openrelay/channels"
	"gopkg.in/redis.v3"
)

// BlockHash will get the latest block hash from the ethereum blockchain
type BlockHash interface {
	Get() (string)
}

// ChanneledBlockHashConsumer listens to a consumerChannel for block hashes,
// and sends them over provided channel
type ChanneledBlockHashConsumer struct {
	channel chan string
}

// Consume processes blockhashes as they arrive from the provided consumer
// channel
func (rbhc *ChanneledBlockHashConsumer) Consume(delivery channels.Delivery) {
	rbhc.channel <- delivery.Payload()
	delivery.Ack()
}

// ChanneledBlockHash is a BlockHash implementation that gets the latest
// block hash by watching a ConsumerChannel
type ChanneledBlockHash struct {
	channel channels.ConsumerChannel
	sourceChan chan string
	sinkChan chan chan string
	started bool
}

// Start kicks off a go routine to listen for changes to the blockhash
func (rbh *ChanneledBlockHash) Start() {
	rbh.channel.AddConsumer(&ChanneledBlockHashConsumer{rbh.sourceChan})
	rbh.channel.StartConsuming()
	go func() {
		// TODO: Make this a random value
		currentHash := "initializing"
		for {
			select {
			case msg := <-rbh.sourceChan:
				currentHash = msg
			case channel := <-rbh.sinkChan:
				channel <- currentHash
			}
		}
	}()
	rbh.started = true
}

// Get retrieves the blockhash from the monitoring go routine
func (rbh *ChanneledBlockHash) Get() string {
	if !rbh.started {
		rbh.Start()
	}
	channel := make(chan string)
	rbh.sinkChan <- channel
	return <-channel
}

// NewChanneledBlockHash returns a BlockHash given a ConsumerChannel
func NewChanneledBlockHash(channel channels.ConsumerChannel) (BlockHash) {
	return &ChanneledBlockHash{
		channel,
		make(chan string),
		make(chan chan string),
		false,
	}
}

// NewRedisBlockHash constructs a ConsumerChannel from a channelURI and a
// redisClient
func NewRedisBlockHash(channelURI string, redisClient *redis.Client) (BlockHash, error) {
	consumerChannel, err := channels.ConsumerFromURI(channelURI, redisClient)
	if err != nil {
		return nil, err
	}
	return NewChanneledBlockHash(consumerChannel), nil
}
