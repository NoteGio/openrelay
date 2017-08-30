package channels

import (
	"gopkg.in/redis.v3"
	"errors"
	"strings"
)

func ConsumerFromURI(uri string, redisClient *redis.Client) (ConsumerChannel, error) {
	if strings.HasPrefix(uri, "topic://") {
		uriTopic := uri[len("topic://"):]
		return NewTopicConsumerChannel(uriTopic, redisClient), nil
	} else if strings.HasPrefix(uri, "queue://") {
		uriQueue := uri[len("queue://"):]
		return NewQueueConsumerChannel(uriQueue, redisClient), nil
	} else {
		return nil, errors.New("Must specify uri starting with queue:// or topic://")
	}
}

func PublisherFromURI(uri string, redisClient *redis.Client) (Publisher, error) {
	if strings.HasPrefix(uri, "topic://") {
		uriTopic := uri[len("topic://"):]
		return NewRedisTopicPublisher(uriTopic, redisClient), nil
	} else if strings.HasPrefix(uri, "queue://") {
		uriQueue := uri[len("queue://"):]
		return NewRedisQueuePublisher(uriQueue, redisClient), nil
	} else {
		return nil, errors.New("Must specify uri starting with queue:// or topic://")
	}
}
