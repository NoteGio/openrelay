package channels

import (
	"errors"
	"gopkg.in/redis.v3"
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

type URITranslator interface {
	ConsumerFromURI(string) (ConsumerChannel, error)
	PublisherFromURI(string) (Publisher, error)
}

type RedisURITranslator struct {
	redisClient *redis.Client
}

func (rut *RedisURITranslator) ConsumerFromURI(uri string) (ConsumerChannel, error) {
	return ConsumerFromURI(uri, rut.redisClient)
}

func (rut *RedisURITranslator) PublisherFromURI(uri string) (Publisher, error) {
	return PublisherFromURI(uri, rut.redisClient)
}

func NewRedisURITranslator(redisClient *redis.Client) (URITranslator) {
	return &RedisURITranslator{redisClient}
}

type MockURITranslator struct {
	redisClient *redis.Client
}

func (rut *MockURITranslator) ConsumerFromURI(uri string) (ConsumerChannel, error) {
	return ConsumerFromURI(uri, rut.redisClient)
}

func (rut *MockURITranslator) PublisherFromURI(uri string) (Publisher, error) {
	return PublisherFromURI(uri, rut.redisClient)
}

func NewMockURITranslator(redisClient *redis.Client) (URITranslator) {
	return &MockURITranslator{redisClient}
}
