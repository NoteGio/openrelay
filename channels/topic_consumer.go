package channels
//TODO: Test channels!
import (
	"gopkg.in/redis.v3"
	// "time"
)

type topicConsumerChannel struct {
	redisClient      *redis.Client
	pubsub           *redis.PubSub
	channelName      string
	consumers        []Consumer
	consumingStopped bool
}

// ReturnAllUnacked is just here for API Compatibility with topics. It does
// nothing
func (topic *topicConsumerChannel) ReturnAllUnacked() int {
	return 0
}

// PurgeRejected is just here for API Compatibility with topics. It does
// nothing
func (topic *topicConsumerChannel) PurgeRejected() int {
	return 0
}

func (topic *topicConsumerChannel) AddConsumer(consumer Consumer) bool {
	topic.consumers = append(topic.consumers, consumer)
	return true
}

func (topic *topicConsumerChannel) StartConsuming() bool {
	// log.Printf("rmq topic started consuming %s %d %s", topic, prefetchLimit, pollDuration)
	if topic.pubsub != nil {
		// Already consuming
		return false
	}
	if pubsub, err := topic.redisClient.Subscribe(topic.channelName); err == nil {
		topic.pubsub = pubsub
		go topic.consume()
		return true
	}
	return false
}

func (topic *topicConsumerChannel) consume() {
	for {
		msg, err := topic.pubsub.ReceiveMessage()
		if err == nil {
			for _, consumer := range topic.consumers {
				go consumer.Consume(newTopicDelivery(msg.Payload, topic.redisClient))
			}
		} else if err.Error() == "redis: client is closed" {
			return
		}
	}
}

func (topic *topicConsumerChannel) StopConsuming() bool {
	if topic.pubsub != nil && !topic.consumingStopped {
		topic.pubsub.Close()
		topic.consumingStopped = true
	}
	return false
}
