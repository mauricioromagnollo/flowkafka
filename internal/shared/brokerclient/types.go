package brokerclient

import (
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

// ClientMetadata represents the metadata information about the Kafka cluster and topics.
type ClientMetadata struct {
	Topics []TopicMetadata
}

// TopicMetadata represents the metadata information about a Kafka topic.
type TopicMetadata struct {
	Name  string
	Error error
}

// Config holds the configuration options for the Kafka producer.
// Brokers is a list of Kafka broker addresses.
// SaslUsername is the username for SASL authentication.
// SaslPassword is the password for SASL authentication.
// TopicName is the name of the Kafka topic to produce messages to.
// SchemaRegistry is the client for interacting with the Schema Registry.
type Config struct {
	Brokers   []string
	TopicName string
	Timeout   time.Duration `default:"10s"`
	Transport *kafkago.Transport
}
