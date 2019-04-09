package subscriptions

import (
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/channels/ws"
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/types"
	"github.com/jinzhu/gorm"
	"log"
)

type WebsocketSubscriptionManager struct {
	manager *SubscriptionManager
}

func NewWebsocketSubscriptionManager() *WebsocketSubscriptionManager {
	return &WebsocketSubscriptionManager{&SubscriptionManager{}}
}

func (subs *WebsocketSubscriptionManager) ListenForSubscriptions(port uint, db *gorm.DB) (func() (error), error) {
	chs, quit := ws.GetChannels(port, db, func(publisher channels.Publisher){
		subs.manager.PruneByPublisher(publisher)
	})
	lookup := dbModule.NewExchangeLookup(db)
	go func() {
		for websocketChannel := range chs {
			ofilter, err := FilterFromQueryString(websocketChannel.Filter)
			if err != nil {
				sendError(websocketChannel, err)
				continue
			}
			subsConsumer := &SubscriptionConsumer{subs.manager, websocketChannel, ofilter, lookup}
			websocketChannel.AddConsumer(subsConsumer)
			websocketChannel.StartConsuming()
		}
	}()
	return quit, nil
}

func (subs *WebsocketSubscriptionManager) Consume(delivery channels.Delivery) {
	order, err := types.OrderFromBytes([]byte(delivery.Payload()))
	if err != nil {
		log.Printf("Error on order: %v", err.Error())
	}
	dbOrder := &dbModule.Order{}
	dbOrder.Order = *order
	dbOrder.Populate()
	subs.manager.Publish(dbOrder)
}
