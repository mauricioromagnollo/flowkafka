package consumer

import (
	"context"
	"fmt"
)

// Consume reads messages from the Kafka topic in a blocking loop.
// Each message is passed to the handler function. Messages are committed
// only after the handler returns nil. The loop exits when the context is
// canceled or the handler returns an error.
func (c *consumerClient) Consume(ctx context.Context, handler func(msg Message) error) error {
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return fmt.Errorf("context error: %w", ctx.Err())
			}
			return fmt.Errorf("failed to fetch message from topic %s: %w", c.cfg.TopicName, err)
		}

		msg := Message{
			Key:       m.Key,
			Value:     m.Value,
			Headers:   m.Headers,
			Partition: m.Partition,
			Offset:    m.Offset,
		}

		if err := handler(msg); err != nil {
			return fmt.Errorf("message handler error: %w", err)
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			return fmt.Errorf("failed to commit message: %w", err)
		}
	}
}
