package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	kafkago "github.com/segmentio/kafka-go"
)

// ProduceJSONMessage produces a message to the specified Kafka topic, encoding the message value as JSON.
func (c *Client) ProduceJSONMessage(ctx context.Context, topic string, key []byte, msg any) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return c.writer.WriteMessages(ctx, kafkago.Message{
		Topic: topic,
		Key:   key,
		Value: b,
	})
}
