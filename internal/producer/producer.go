package producer

import (
	"context"
)

// Producer is an interface that defines the methods for interacting with a Kafka producer.
type Producer interface {
	// ValidateConnection checks if the producer can connect to the Kafka cluster.
	// It takes a context and returns an error if the connection is not valid.
	ValidateConnection(ctx context.Context) error
	// Publish sends one JSON message to a Kafka topic.
	// It takes a context, the key, and the message to be sent.
	// It returns an error if there was a failure in sending the message.
	Publish(ctx context.Context, key []byte, msg any) error
	// PublishAvro sends one Avro serialized message to a Kafka topic using Schema Registry (Confluent wire format).
	// It takes a context, the key, and the message to be sent.
	// It returns an error if there was a failure in sending the message.
	PublishAvro(ctx context.Context, key []byte, msg any) error
	// Close closes the producer and releases any resources it holds.
	Close() error
	// HasSchemaRegistry returns true if the producer is configured with a Schema Registry client.
	HasSchemaRegistry() bool
}
