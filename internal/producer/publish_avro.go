package producer

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"

	avro "github.com/hamba/avro/v2"
)

// PublishAvro sends a message to a Kafka topic using Avro serialization
// with Schema Registry. The message is serialized in the Confluent wire format:
// [magic_byte(1)][schema_id(4)][avro_payload(n)]
func (c *producerClient) PublishAvro(ctx context.Context, key []byte, msg any) error {
	if !c.HasSchemaRegistry() {
		return fmt.Errorf("schema registry is not configured")
	}

	// Get the latest schema from Schema Registry
	latestSchema, err := c.schemaRegistry.GetLatestSchema(c.cfg.TopicName + "-value")
	if err != nil {
		return fmt.Errorf("failed to get latest schema for topic %s: %w", c.cfg.TopicName, err)
	}

	// Parse the Avro schema
	codec, err := avro.Parse(latestSchema.Schema())
	if err != nil {
		return fmt.Errorf("failed to parse schema: %w", err)
	}

	// Marshal the message to Avro binary format
	avroPayload, err := avro.Marshal(codec, msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message to Avro: %w", err)
	}

	// Build the Confluent wire format:
	// [magic_byte(1 byte)][schema_id(4 bytes)][avro_payload(n bytes)]
	schemaID := latestSchema.ID()
	if schemaID < 0 || schemaID > math.MaxUint32 {
		return fmt.Errorf("schema ID %d is out of range for uint32", schemaID)
	}
	messageValue := make([]byte, 5+len(avroPayload))
	messageValue[0] = 0                                             // Magic byte
	binary.BigEndian.PutUint32(messageValue[1:5], uint32(schemaID)) // #nosec G115
	copy(messageValue[5:], avroPayload)

	// Send message to Kafka using the persistent writer
	if err := c.writer.WriteMessages(ctx, Message{
		Value: messageValue,
		Key:   key,
	}); err != nil {
		return fmt.Errorf("failed to publish message to Kafka topic %s: %w", c.cfg.TopicName, err)
	}

	return nil
}
