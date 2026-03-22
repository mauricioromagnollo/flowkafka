package producer

import (
	"context"

	"github.com/mauricioromagnollo/flowkafka/internal/shared/types"
)

// Writer defines the interface for writing messages to Kafka.
type Writer interface {
	// WriteMessages writes a message to Kafka.
	WriteMessages(ctx context.Context, msg types.Message) error
	// Close closes the writer and releases any resources.
	Close() error
}
