package channels

import (
  "gopkg.in/redis.v3"
)

type Publisher interface {
  Publish(payload string) bool
}

type queuePublisher struct {
  key string
  redisClient *redis.Client
}

func NewQueuePublisher(key string, client *redis.Client)(Publisher) {
  return &queuePublisher{key, client}
}

func (publisher *queuePublisher)Publish(payload string) bool {
  return !redisErrIsNil(publisher.redisClient.LPush(publisher.key, payload))
}

type topicPublisher struct {
  key string
  redisClient *redis.Client
}

func NewTopicPublisher(key string, client *redis.Client)(Publisher) {
  return &topicPublisher{key, client}
}

func (publisher *topicPublisher)Publish(payload string) bool {
  return !redisErrIsNil(publisher.redisClient.Publish(publisher.key + "::ready", payload))
}
