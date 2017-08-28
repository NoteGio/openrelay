package channels
//TODO: Test channels!
import (
	"gopkg.in/redis.v3"
	// "time"
	// "fmt"
)

type topicConsumerChannel struct {
	redisClient      *redis.Client
	pubsub           *redis.PubSub
	channelName      string
	consumers        []Consumer
	consumingStopped bool
}

// NewTopicConsumerChannel returns a ConsumerChannel that uses Redis PubSub for
// communication. Each message delivered through this consumer channel will be
// delivered once to each consumer. Note, however, that network issues that
// prevent delivery of a message may lead to messages going completely
// undelivered. Consumers may Ack or Reject the messages, but this is a no-op.
func NewTopicConsumerChannel(channelName string, redisClient *redis.Client) (ConsumerChannel) {
	return &topicConsumerChannel{
		redisClient,
		nil,
		channelName,
		[]Consumer{},
		false,
	}
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
