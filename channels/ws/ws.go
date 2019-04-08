package ws

import (
	"context"
	"fmt"
	"github.com/notegio/openrelay/channels"
	"github.com/notegio/openrelay/types"
	"github.com/notegio/openrelay/pool"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"net/http"
	"log"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type websocketDelivery struct {
	payload string
}

func (delivery *websocketDelivery) Payload() string {
	return delivery.payload
}

func (delivery *websocketDelivery) Ack() bool {
	// websocketDeliveris have no ack, reject, or return, so these are no-ops
	return true
}
func (delivery *websocketDelivery) Reject() bool {
	return true
}
func (delivery *websocketDelivery) Return() bool {
	return true
}

type WebsocketChannel struct {
	open bool
	conn *websocket.Conn
	payloads chan []byte
	consumers []channels.Consumer
	Filter string
	quit chan struct{}
	cleanup func()
}

func (pub *WebsocketChannel) Publish(payload string) bool {
	pub.payloads <- []byte(payload)
	return true
}

func (consumerChannel *WebsocketChannel) AddConsumer(consumer channels.Consumer) bool {
	consumerChannel.consumers = append(consumerChannel.consumers, consumer)
	return true
}
func (consumerChannel *WebsocketChannel) StartConsuming() bool {
	go func () {
		defer consumerChannel.cleanup()
		for {
			select {
			case _ = <-consumerChannel.quit:
				return
			default:
			}
			_, p, err := consumerChannel.conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			for _, consumer := range consumerChannel.consumers {
				consumer.Consume(&websocketDelivery{string(p)})
			}
		}
	}()
	return true
}

func (consumerChannel *WebsocketChannel) StopConsuming() bool {
	consumerChannel.quit <- struct{}{}
	return true
}
func (consumerChannel *WebsocketChannel) ReturnAllUnacked() int {
	return 0
}
func (consumerChannel *WebsocketChannel) PurgeRejected() int {
	return 0
}
func (consumerChannel *WebsocketChannel) Publisher() channels.Publisher {
	return consumerChannel
}

func GetChannels(port uint, db *gorm.DB, cleanup func()) (<-chan *WebsocketChannel, func() (error)) {
	outChan := make(chan *WebsocketChannel)
	handler := pool.PoolDecorator(db, func (w http.ResponseWriter, r *http.Request, p types.Pool) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
		wsChannel := &WebsocketChannel{true, conn, make(chan []byte), []channels.Consumer{}, p.QueryString(), make(chan struct{}), cleanup}
		outChan <- wsChannel
		for payload := range wsChannel.payloads {
			if err := conn.WriteMessage(websocket.BinaryMessage, payload); err != nil {
				log.Println(err)
				return
			}
		}
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", port),
		Handler: mux,
	}
	go func() {
		log.Printf("%v", srv.ListenAndServe())
	}()
	return outChan, func() (error) { return srv.Shutdown(context.Background()) }
}
