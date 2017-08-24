package channels

import (
	"gopkg.in/redis.v3"
	"time"
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
		// debug(fmt.Sprintf("rmq queue returned unacked delivery %s %s", result.Val(), queue.readyKey)) // COMMENTOUT
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
	// log.Printf("rmq queue started consuming %s %d %s", queue, prefetchLimit, pollDuration)
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
			// log.Printf("rmq queue stopped consuming %s", queue)
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
