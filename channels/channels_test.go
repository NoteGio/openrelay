package channels_test

import (
	"github.com/notegio/0xrelay/channels"
	"gopkg.in/redis.v3"
	"testing"
	"os"
	"fmt"
)

type testConsumer struct {
	channel chan string
	ack chan bool
	done chan bool
}

func (consumer *testConsumer)Consume(msg channels.Delivery){
	consumer.channel <- msg.Payload()
	if <-consumer.ack {
		msg.Ack()
	} else {
		msg.Reject()
	}
	consumer.done <- true
}

func ChannelSendTest(publisher channels.Publisher, consumerChannel channels.ConsumerChannel, t *testing.T) {
	fmt.Println("ChannelSendTest")
	// fmt.Println(publisher)
	consumer := &testConsumer{make(chan string), make(chan bool), make(chan bool)}
	// fmt.Println("Created consumer")
	consumerChannel.AddConsumer(consumer)
	// fmt.Println("Added consumer")
	consumerChannel.StartConsuming()
	// fmt.Println("Started consumer")
	publisher.Publish("test")
	// fmt.Println("Published message")
	result := <-consumer.channel
	// fmt.Println("Got result")
	if result != "test" {
		t.Errorf("Unexpected value")
	}
	consumer.ack <- true
	_ = <-consumer.done
}

func ReturnUnackedTest(publisher channels.Publisher, consumerChannel channels.ConsumerChannel, t *testing.T) {
	fmt.Println("ReturnUnackedTest")
	// fmt.Println(publisher)
	consumer := &testConsumer{make(chan string), make(chan bool), make(chan bool)}
	// fmt.Println("Created consumer")
	consumerChannel.AddConsumer(consumer)
	// fmt.Println("Added consumer")
	consumerChannel.StartConsuming()
	// fmt.Println("Started consumer")
	publisher.Publish("test")
	// fmt.Println("Published message")
	result := <-consumer.channel
	// fmt.Println("Got result")
	if result != "test" {
		t.Errorf("Unexpected value")
	}
	// fmt.Println("Checked result")
	if unackedCount := consumerChannel.ReturnAllUnacked(); unackedCount != 1 {
		t.Errorf("Expected 1 unacked value, got '%v'", unackedCount)
	}
	consumer.ack <- true
	_ = <-consumer.done
	// fmt.Println("done here")
}

func AckTest(publisher channels.Publisher, consumerChannel channels.ConsumerChannel, t *testing.T) {
	fmt.Println("AckTest")
	// fmt.Println(publisher)
	consumer := &testConsumer{make(chan string), make(chan bool), make(chan bool)}
	// fmt.Println("Created consumer")
	consumerChannel.AddConsumer(consumer)
	// fmt.Println("Added consumer")
	consumerChannel.StartConsuming()
	// fmt.Println("Started consumer")
	publisher.Publish("test")
	// fmt.Println("Published message")
	result := <-consumer.channel
	// fmt.Println("Got result")
	if result != "test" {
		t.Errorf("Unexpected value")
	}
	consumer.ack <- true
	_ = <-consumer.done
	if unackedCount := consumerChannel.ReturnAllUnacked(); unackedCount != 0 {
		t.Errorf("Expected 0 unacked value, got '%v'", unackedCount)
	}
}

func RejectTest(publisher channels.Publisher, consumerChannel channels.ConsumerChannel, t *testing.T) {
	fmt.Println("RejectTest")
	// fmt.Println(publisher)
	consumer := &testConsumer{make(chan string), make(chan bool), make(chan bool)}
	// fmt.Println("Created consumer")
	consumerChannel.AddConsumer(consumer)
	// fmt.Println("Added consumer")
	consumerChannel.StartConsuming()
	// fmt.Println("Started consumer")
	publisher.Publish("test")
	// fmt.Println("Published message")
	result := <-consumer.channel
	// fmt.Println("Got result")
	if result != "test" {
		t.Errorf("Unexpected value")
	}
	consumer.ack <- false
	_ = <-consumer.done
	if unackedCount := consumerChannel.ReturnAllUnacked(); unackedCount != 0 {
		t.Errorf("Expected 0 unacked value, got '%v'", unackedCount)
	}
	if purgedCount := consumerChannel.PurgeRejected(); purgedCount != 1 {
		t.Errorf("Expected 1 purged value, got '%v'", purgedCount)
	}
}

func TestMockChannelSend(t *testing.T) {
	publisher, consumerChannel := channels.MockChannel()
	ChannelSendTest(publisher, consumerChannel, t)
}
func TestMockReturnUnacked(t *testing.T) {
	publisher, consumerChannel := channels.MockChannel()
	ReturnUnackedTest(publisher, consumerChannel, t)
}
func TestMockAck(t *testing.T) {
	publisher, consumerChannel := channels.MockChannel()
	AckTest(publisher, consumerChannel, t)
}
func TestMockReject(t *testing.T) {
	publisher, consumerChannel := channels.MockChannel()
	RejectTest(publisher, consumerChannel, t)
}

func redisCleanup(redisClient *redis.Client, consumerChannel channels.ConsumerChannel) {
	key := "test_queue::unacked"
	for int(redisClient.LLen(key).Val()) > 0 {
		redisClient.RPop(key)
	}
	consumerChannel.StopConsuming()
}

func TestRedisQueueChannelSend(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		t.Errorf("Please set the REDIS_URL environment variable")
		return
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	publisher := channels.NewQueuePublisher("test_queue", redisClient)
	consumerChannel := channels.NewQueueConsumerChannel("test_queue", redisClient)
	defer redisCleanup(redisClient, consumerChannel)
	ChannelSendTest(publisher, consumerChannel, t)
}
func TestRedisQueueReturnUnacked(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		t.Errorf("Please set the REDIS_URL environment variable")
		return
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	publisher := channels.NewQueuePublisher("test_queue", redisClient)
	consumerChannel := channels.NewQueueConsumerChannel("test_queue", redisClient)
	defer redisCleanup(redisClient, consumerChannel)
	ReturnUnackedTest(publisher, consumerChannel, t)
}
func TestRedisQueueAck(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		t.Errorf("Please set the REDIS_URL environment variable")
		return
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	publisher := channels.NewQueuePublisher("test_queue", redisClient)
	consumerChannel := channels.NewQueueConsumerChannel("test_queue", redisClient)
	defer redisCleanup(redisClient, consumerChannel)
	AckTest(publisher, consumerChannel, t)
}
func TestRedisQueueReject(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		t.Errorf("Please set the REDIS_URL environment variable")
		return
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	publisher := channels.NewQueuePublisher("test_queue", redisClient)
	consumerChannel := channels.NewQueueConsumerChannel("test_queue", redisClient)
	defer redisCleanup(redisClient, consumerChannel)
	RejectTest(publisher, consumerChannel, t)
}

func TestRedisTopicChannelSend(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		t.Errorf("Please set the REDIS_URL environment variable")
		return
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	publisher := channels.NewTopicPublisher("test_topic", redisClient)
	consumerChannel := channels.NewTopicConsumerChannel("test_topic", redisClient)
	defer consumerChannel.StopConsuming()
	ChannelSendTest(publisher, consumerChannel, t)
}
func TestRedisTopicAck(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		t.Errorf("Please set the REDIS_URL environment variable")
		return
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	publisher := channels.NewTopicPublisher("test_topic", redisClient)
	consumerChannel := channels.NewTopicConsumerChannel("test_topic", redisClient)
	defer consumerChannel.StopConsuming()
	AckTest(publisher, consumerChannel, t)
}
