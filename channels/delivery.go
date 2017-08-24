package channels

import (
  "gopkg.in/redis.v3"
)

type Delivery interface {
	Payload() string
	Ack() bool
	Reject() bool
}

type topicDelivery struct {
  payload string
  redisClient *redis.Client
}

func (delivery *topicDelivery) Payload() string {
  return delivery.payload
}

func (delivery *topicDelivery) Ack() bool {
  // Topics can't actually be Ack'd, but we want the interface to be the same
  return true
}

func (delivery *topicDelivery) Reject() bool {
  // Topics can't actually be Rejected, but we want the interface to be the same
  return true
}

func newTopicDelivery(payload string, client *redis.Client) (*topicDelivery) {
  return &topicDelivery{payload, client}
}

type queueDelivery struct {
  payload string
  unackedKey string
  rejectedKey string
  redisClient *redis.Client
}

func (delivery *queueDelivery) Payload() string {
  return delivery.payload
}

func (delivery *queueDelivery) Ack() bool {
  result := delivery.redisClient.LRem(delivery.unackedKey, 1, delivery.payload)
  if redisErrIsNil(result) {
    return false
  }
  return result.Val() == 1
}

func (delivery *queueDelivery) Reject() bool {
  return delivery.move(delivery.rejectedKey)

}

func (delivery *queueDelivery) move(key string) bool {
	if redisErrIsNil(delivery.redisClient.LPush(key, delivery.payload)) {
		return false
	}

	if redisErrIsNil(delivery.redisClient.LRem(delivery.unackedKey, 1, delivery.payload)) {
		return false
	}

	// debug(fmt.Sprintf("delivery rejected %s", delivery)) // COMMENTOUT
	return true
}

func newQueueDelivery(payload, unackedKey, rejectedKey string, client *redis.Client) (*queueDelivery) {
  return &queueDelivery{payload, unackedKey, rejectedKey, client}
}
