package cmdutils

import (
	"github.com/notegio/openrelay/channels"
	"strings"
	"gopkg.in/redis.v3"
)

// x=>y=>z;q

func ParseChannels(channelString string, redisClient *redis.Client) (channels.ConsumerChannel, channels.MultiPublisher, channels.Publisher, error) {
	publishers := channels.MultiPublisher{}
	altPublisherStringSlice := strings.Split(channelString, ";")
	var altPublisher channels.Publisher
	if len(altPublisherStringSlice) > 1 {
		var err error
		altPublisher, err = channels.PublisherFromURI(altPublisherStringSlice[1], redisClient)
		if err != nil {
			return nil, publishers, nil, err
		}
	}
	channelStringSlice := strings.Split(altPublisherStringSlice[0], "=>")
	sourceChannel, err := channels.ConsumerFromURI(channelStringSlice[0], redisClient)
	if err != nil {
		return nil, publishers, nil, err
	}
	for _, channelString := range channelStringSlice[1:] {
		publisher, err := channels.PublisherFromURI(channelString, redisClient)
		if err != nil {
			return nil, publishers, nil, err
		}
		publishers = append(publishers[:], publisher)
	}
	return sourceChannel, publishers, altPublisher, nil
}
