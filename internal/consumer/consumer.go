package consumer

import (
	"context"
)

// Consumer is an interface that defines the methods for interacting with a Kafka consumer.
type Consumer interface {
	// Consume reads messages from the Kafka topic in a blocking loop.
	// Each message is passed to the handler function. Messages are committed
	// only after the handler returns nil. The loop exits when the context is
	// canceled or the handler returns an error.
	Consume(ctx context.Context, handler func(msg Message) error) error
	// ConsumeMessages reads messages from the Kafka topic and sends them to the
	// provided channel. It blocks until the context is canceled. The caller
	// owns the channel and must close it after ConsumeMessages returns.
	ConsumeMessages(ctx context.Context, msgsChan chan<- Message) error
	// Close closes the consumer and releases any resources it holds.
	Close() error
	// ValidateConnection checks if the consumer can connect to the Kafka cluster.
	ValidateConnection(ctx context.Context) error
}
