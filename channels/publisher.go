package channels

import (
	"gopkg.in/redis.v3"
	"log"
)

type Publisher interface {
	Publish(payload string) bool
}

type redisQueuePublisher struct {
	key         string
	redisClient *redis.Client
}

func NewRedisQueuePublisher(key string, client *redis.Client) Publisher {
	return &redisQueuePublisher{key, client}
}

func (publisher *redisQueuePublisher) Publish(payload string) bool {
	if len(payload) == 0 {
		log.Printf("Trying to publish empty message. Skipping")
		return false
	}
	return !redisErrIsNil(publisher.redisClient.LPush(publisher.key, payload))
}

type redisTopicPublisher struct {
	key         string
	redisClient *redis.Client
}

func NewRedisTopicPublisher(key string, client *redis.Client) Publisher {
	return &redisTopicPublisher{key, client}
}

func (publisher *redisTopicPublisher) Publish(payload string) bool {
	return !redisErrIsNil(publisher.redisClient.Publish(publisher.key, payload))
}
