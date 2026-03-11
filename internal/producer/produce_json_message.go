package producer

import (
	"context"
	"encoding/json"
	"fmt"

	kafkago "github.com/segmentio/kafka-go"
)

// ProduceMessage sends one message to a Kafka topic.
// It takes a context, the topic name, and the message to be sent.
// It returns an error if there was a failure in sending the message.
func (c *producerClient) ProduceJSONMessage(ctx context.Context, key []byte, msg any) (err error) {
	message, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message to publish: %w", err)
	}

	if err := c.writer.WriteMessages(ctx, kafkago.Message{
		Value: message,
		Key:   key,
	}); err != nil {
		return fmt.Errorf("failed to publish message in kafka: %w", err)
	}

	return nil
}
