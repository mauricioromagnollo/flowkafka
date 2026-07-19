package consumer

import (
	"context"

	kafkago "github.com/segmentio/kafka-go"
)

// Reader defines the interface for reading messages from Kafka.
type Reader interface {
	// FetchMessage reads the next message from Kafka without committing it.
	FetchMessage(ctx context.Context) (kafkago.Message, error)
	// CommitMessages commits the given messages to Kafka.
	CommitMessages(ctx context.Context, msgs ...kafkago.Message) error
	// Close closes the reader and releases any resources.
	Close() error
}
