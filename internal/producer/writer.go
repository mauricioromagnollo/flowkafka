package producer

import (
	"context"
)

// Writer defines the interface for writing messages to Kafka.
type Writer interface {
	// WriteMessages writes a message to Kafka.
	WriteMessages(ctx context.Context, msg Message) error
	// Close closes the writer and releases any resources.
	Close() error
}
