// This module borrows heavily from https://github.com/adjust/rmq
// Stripping it down to a limited subset that we need that can also be made
// API compatible with topics.
//
// The MIT License (MIT)
//
// Copyright (c) 2015 adjust
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package channels

import (
	"gopkg.in/redis.v3"
	"time"
	// "fmt"
)

type queueConsumerChannel struct {
	redisClient      *redis.Client
	channelName      string
	readyKey         string
	unackedKey       string
	rejectedKey      string
	consumingStopped chan bool
	deliveryChan     chan Delivery
}

// NewQueueConsumerChannel returns a ConsumerChannel that uses Redis queues for
// communication. Each message delivered through this ConsumerChannel will be
// delivered to only one consumer, assuming the consumer Acks the message.
func NewQueueConsumerChannel(channelName string, redisClient *redis.Client) ConsumerChannel {
	return &queueConsumerChannel{
		redisClient,
		channelName,
		channelName,
		channelName + "::unacked",
		channelName + "::rejected",
		nil,
		nil,
	}
}

// ReturnAllUnacked moves all unacked deliveries back to the ready
// queue and deletes the unacked key afterwards, returns number of returned
// deliveries
func (queue *queueConsumerChannel) ReturnAllUnacked() int {
	result := queue.redisClient.LLen(queue.unackedKey)
	if redisErrIsNil(result) {
		return 0
	}

	unackedCount := int(result.Val())
	for i := 0; i < unackedCount; i++ {
		if redisErrIsNil(queue.redisClient.RPopLPush(queue.unackedKey, queue.readyKey)) {
			return i
		}
	}

	return unackedCount
}

// PurgeRejected removes all rejected deliveries from the queue and returns the number of purged deliveries
func (queue *queueConsumerChannel) PurgeRejected() int {
	return queue.deleteRedisList(queue.rejectedKey)
}

// return number of deleted list items
// https://www.redisgreen.net/blog/deleting-large-lists
func (queue *queueConsumerChannel) deleteRedisList(key string) int {
	llenResult := queue.redisClient.LLen(key)
	total := int(llenResult.Val())
	if total == 0 {
		return 0 // nothing to do
	}

	// delete elements without blocking
	for todo := total; todo > 0; todo -= purgeBatchSize {
		// minimum of purgeBatchSize and todo
		batchSize := purgeBatchSize
		if batchSize > todo {
			batchSize = todo
		}

		// remove one batch
		queue.redisClient.LTrim(key, 0, int64(-1-batchSize))
	}

	return total
}

func (queue *queueConsumerChannel) AddConsumer(consumer Consumer) bool {
	go func() {
		for delivery := range queue.deliveryChan {
			consumer.Consume(delivery)
		}
	}()
	return true
}

func (queue *queueConsumerChannel) StartConsuming() bool {
	if queue.deliveryChan != nil {
		return false // already consuming
	}

	queue.deliveryChan = make(chan Delivery, prefetchLimit)
	go queue.consume()
	return true
}

func (queue *queueConsumerChannel) consume() {
	for {
		result := queue.redisClient.BRPopLPush(queue.readyKey, queue.unackedKey, time.Second)
		if !redisErrIsNil(result) {
			queue.deliveryChan <- newQueueDelivery(result.Val(), queue.unackedKey, queue.rejectedKey, queue.redisClient)
		}
		if queue.consumingStopped != nil {
			queue.consumingStopped <- true
			return
		}
	}
}

func (queue *queueConsumerChannel) StopConsuming() bool {
	if queue.deliveryChan != nil && queue.consumingStopped == nil {
		queue.consumingStopped = make(chan bool)
		return <-queue.consumingStopped
	}
	return false
}
