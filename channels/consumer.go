package channels

const (
  purgeBatchSize      = 100
  prefetchLimit       = 5
)

type Consumer interface {
  Consume(Delivery)
}

type ConsumerChannel interface {
  AddConsumer(Consumer)
  StartConsuming() error
  StopConsuming() error
  ReturnAllUnacked() error
  PurgeRejected() error
}
