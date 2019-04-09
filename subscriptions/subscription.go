package subscriptions

import (
	"encoding/json"
	"github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/search"
	"log"
)

type Subscription struct {
	publisher channels.Publisher
	filter func(*db.Order) bool
	requestID string
	internalId int64
}

type SubscriptionUpdate struct {
	Type      string `json:"type"`
	Channel   string `json:"channel"`
	RequestID string `json:"requestId"`
	Payload   []interface{} `json:"payload"`
}


func (subscription *Subscription) Publish(order *db.Order) (bool) {
	if subscription.filter(order) {
		formatted := search.GetFormattedOrder(*order)
		message := &SubscriptionUpdate{
			Type: "update",
			Channel: "orders",
			RequestID: subscription.requestID,
			Payload: []interface{}{formatted},
		}
		data, err := json.Marshal(message)
		if err != nil {
			log.Printf(err.Error())
		}
		return subscription.publisher.Publish(string(data))
	}
	return true
}
