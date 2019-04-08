package subscriptions

import (
	"encoding/json"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/db"
	"log"
)


type SubscriptionConsumer struct {
	manager *SubscriptionManager
	publisher channels.Publisher
	baseFilter *OrderFilter
	exchangeLookup ExchangeLookup
}

func NewSubscriptionConsumer(manager *SubscriptionManager, publisher channels.Publisher, filter *OrderFilter, lookup ExchangeLookup) (*SubscriptionConsumer) {
	return &SubscriptionConsumer{manager, publisher, filter, lookup}
}

type SubscriptionMessage struct {
	Type      string `json:"type"`
	Channel   string `json:"channel"`
	RequestID string `json:"requestId"`
	Payload   *OrderFilter
}

func sendError(publisher channels.Publisher, err error) {
	message := SubscriptionUpdate{
		Type: "error",
		Channel: "orders",
		RequestID: "unknown",
		Payload: []interface{}{err.Error()},
	}
	data, _ := json.Marshal(message)
	publisher.Publish(string(data))
}

func (consumer *SubscriptionConsumer) Consume(delivery channels.Delivery) {
	incoming := &SubscriptionMessage{}
	if err := json.Unmarshal([]byte(delivery.Payload()), incoming); err != nil {
		log.Printf("Error parsing JSON: %v", err.Error())
		sendError(consumer.publisher, err)
		return
	}
	baseFilterFn := func(order *db.Order) (bool) { return true }
	payloadFilterFn := func(order *db.Order) (bool) { return true }
	if consumer.baseFilter != nil {
		var err error
		baseFilterFn, err = consumer.baseFilter.GetFilter(consumer.exchangeLookup)
		if err != nil {
			log.Printf("Error getting filter %v", err.Error())
			sendError(consumer.publisher, err)
			return
		}
	}
	if incoming.Payload != nil {
		var err error
		payloadFilterFn, err = incoming.Payload.GetFilter(consumer.exchangeLookup)
		if err != nil {
			log.Printf("Error getting filter %v", err.Error())
			sendError(consumer.publisher, err)
			return
		}
	}
	subscription := Subscription{
		publisher: consumer.publisher,
		requestID: incoming.RequestID,
		filter: func(order *db.Order) (bool) { return baseFilterFn(order) && payloadFilterFn(order) },
	}
	consumer.manager.Add(subscription)
}
