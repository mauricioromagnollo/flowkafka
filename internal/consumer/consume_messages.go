package consumer

import (
	"context"
	"fmt"
)

// ConsumeMessages reads messages from the Kafka topic and sends them to the
// provided channel. It blocks until the context is canceled, at which point
// it returns nil. If a fetch or commit error occurs it returns immediately.
// The caller owns the channel and is responsible for closing it after
// ConsumeMessages returns.
func (c *consumerClient) ConsumeMessages(ctx context.Context, msgsChan chan<- Message) error {
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

		select {
		case msgsChan <- msg:
		case <-ctx.Done():
			return nil
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			return fmt.Errorf("failed to commit message: %w", err)
		}
	}
}
