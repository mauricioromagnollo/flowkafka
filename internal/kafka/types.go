package kafka

import "github.com/mauricioromagnollo/flowkafka/internal/schemaregistry"

// Config represents configuration for one Kafka client by cluster.
type Config struct {
	Addresses      []string
	GroupID        string
	SaslUsername   string
	SaslPassword   string
	SchemaRegistry schemaregistry.SchemaRegistry
}

// Message represents a Kafka message with its associated metadata.
// It encapsulates the key-value pair of the message along with headers,
// partition information, and offset for message tracking and processing.
//
// Fields:
//   - Key: The message key as a byte array, used for partitioning and message identification
//   - Value: The actual message payload as a byte array
//   - Headers: A map of custom headers attached to the message, with string keys and byte array values
//   - Partition: The Kafka partition number from which this message was consumed
//   - Offset: The position of the message within its partition, used for ordering and tracking
type Message struct {
	Key       []byte
	Value     []byte
	Headers   map[string][]byte
	Partition int
	Offset    int64
}
