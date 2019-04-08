package subscriptions

import (
	"github.com/notegio/openrelay/db"
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
		manager.Remove(failures[len(failures) - j])
	}
}

func (manager *SubscriptionManager) Remove(index int) {
	lastIndex := len(manager.subscriptions) - 1
	manager.subscriptions[index] = manager.subscriptions[lastIndex]
	manager.subscriptions = manager.subscriptions[:lastIndex]
}

func (manager *SubscriptionManager) Add(subscription Subscription) {
	manager.subscriptions = append(manager.subscriptions, subscription)
}
