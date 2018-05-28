package channels

const (
	purgeBatchSize = 100
	prefetchLimit  = 5
)

type Consumer interface {
	Consume(Delivery)
}

type ConsumerChannel interface {
	AddConsumer(Consumer) bool
	StartConsuming() bool
	StopConsuming() bool
	ReturnAllUnacked() int
	PurgeRejected() int
	Publisher() Publisher
}
