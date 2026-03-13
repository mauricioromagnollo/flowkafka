package producer

import (
	"time"

	"github.com/mauricioromagnollo/flowkafka/internal/schemaregistry"
)

// Config holds the configuration options for the Kafka producer.
// Brokers is a list of Kafka broker addresses.
// SaslUsername is the username for SASL authentication.
// SaslPassword is the password for SASL authentication.
// TopicName is the name of the Kafka topic to produce messages to.
// SchemaRegistry is the client for interacting with the Schema Registry.
type Config struct {
	Brokers        []string
	SaslUsername   string
	SaslPassword   string
	TopicName      string
	SchemaRegistry schemaregistry.SchemaRegistry
	// Optional Props
	MaxAttempts  int           `default:"10"`
	BatchTimeout time.Duration `default:"1s"`
	WriteTimeout time.Duration `default:"30s"`
	ReadTimeout  time.Duration `default:"30s"`
	RequiredAcks RequiredAcks  `default:"-1"`
}

// RequiredAcks defines the acknowledgment level required from the Kafka cluster for a message to be considered successfully sent.
type RequiredAcks int

const (
	// RequiredAcksNone means the producer does not wait for any acknowledgment from the Kafka cluster. This provides the lowest latency but the highest risk of message loss.
	RequiredAcksNone RequiredAcks = 0
	// RequiredAcksOne means the producer waits for an acknowledgment from the leader broker only. This provides a balance between latency and durability.
	RequiredAcksOne RequiredAcks = 1
	// RequiredAcksAll means the producer waits for acknowledgments from all in-sync replicas. This provides the highest level of durability but also the highest latency.
	RequiredAcksAll RequiredAcks = -1
)

// Message represents a Kafka message.
// Key is the key of the message.
// Value is the value of the message.
// Headers contains any additional headers associated with the message.
// Partition is the partition number of the message.
// Offset is the offset of the message within the partition.
type Message struct {
	Key       []byte
	Value     []byte
	Headers   any
	Partition int
	Offset    int64
}
