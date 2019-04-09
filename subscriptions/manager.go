package subscriptions

import (
	"github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/channels"
	"log"
)

type SubscriptionManager struct {
	subscriptions []Subscription
}

func (manager *SubscriptionManager) Publish(order *db.Order) {
	failures := []int{}
	for i, subscription := range manager.subscriptions {
		if !subscription.Publish(order) {
			failures = append(failures, i)
		}
	}
	for j := range failures {
		// By going backwards through the failures, we'll never having ordering
		// issues as we remove items
		manager.Remove(failures[len(failures) - (j+1)])
	}
}

func (manager *SubscriptionManager) Remove(index int) {
	lastIndex := len(manager.subscriptions) - 1
	log.Printf("Removing subscription: %v", manager.subscriptions[index].requestID)
	manager.subscriptions[index] = manager.subscriptions[lastIndex]
	manager.subscriptions = manager.subscriptions[:lastIndex]
}

func (manager *SubscriptionManager) PruneByPublisher(channel channels.Publisher)  {
	toPrune := []int{}
	for i, subscription := range manager.subscriptions {
		if channel == subscription.publisher {
			toPrune = append(toPrune, i)
		}
	}
	for j := range toPrune {
		// By going backwards through the toPrune, we'll never having ordering
		// issues as we remove items
		manager.Remove(toPrune[len(toPrune) - (j+1)])
	}
}

func (manager *SubscriptionManager) Add(subscription Subscription) {
	manager.subscriptions = append(manager.subscriptions, subscription)
}
