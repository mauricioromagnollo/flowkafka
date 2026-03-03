package kafka

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/hamba/avro/v2"
	kafkago "github.com/segmentio/kafka-go"
)

// ProduceAvroMessage produces a message to the specified Kafka topic, encoding the message value using Avro serialization.
func (c *Client) ProduceAvroMessage(ctx context.Context, topic string, key []byte, msg any) error {
	if c.config.SchemaRegistry == nil {
		return fmt.Errorf("schema registry client is nil")
	}

	subject := topic + "-value"
	latest, err := c.config.SchemaRegistry.GetLatestSchema(subject)
	if err != nil {
		return fmt.Errorf("failed to get latest schema for subject %s: %w", subject, err)
	}

	codec, err := avro.Parse(latest.Schema())
	if err != nil {
		return fmt.Errorf("failed to parse avro schema: %w", err)
	}

	avroPayload, err := avro.Marshal(codec, msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message to avro: %w", err)
	}

	schemaID := latest.ID()
	if schemaID < 0 {
		return fmt.Errorf("invalid schema id: %d", schemaID)
	}

	value := make([]byte, 5+len(avroPayload))
	value[0] = 0
	binary.BigEndian.PutUint32(value[1:5], uint32(schemaID)) //nolint:gosec // G115: schemaID is validated as non-negative above
	copy(value[5:], avroPayload)

	return c.writer.WriteMessages(ctx, kafkago.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}
