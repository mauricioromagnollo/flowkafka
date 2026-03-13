package flowkafka

import (
	"context"

	"github.com/mauricioromagnollo/flowkafka/internal/producer"
)

// Producer is an interface that defines the methods for interacting with a Kafka producer.
type Producer = producer.Producer

// ProducerConfig holds the configuration options for the Kafka producer.
type ProducerConfig = producer.Config

type producerWrapper struct {
	client producer.Producer
}

// NewProducer creates a new Kafka producer with the given configuration.
func NewProducer(config ProducerConfig) Producer {
	return producerWrapper{
		client: producer.NewProducer(config),
	}
}

func (p producerWrapper) ValidateConnection(ctx context.Context) error {
	return p.client.ValidateConnection(ctx)
}

func (p producerWrapper) Publish(ctx context.Context, key []byte, msg any) error {
	return p.client.Publish(ctx, key, msg)
}

func (p producerWrapper) PublishAvro(ctx context.Context, key []byte, msg any) error {
	return p.client.PublishAvro(ctx, key, msg)
}

func (p producerWrapper) Close() error {
	return p.client.Close()
}

func (p producerWrapper) HasSchemaRegistry() bool {
	return p.client.HasSchemaRegistry()
}

// RequiredAcks defines the acknowledgment level required from the Kafka cluster for a message to be considered successfully sent.
type RequiredAcks = producer.RequiredAcks

const (
	// RequiredAcksNone means the producer does not wait for any acknowledgment from the Kafka cluster. This provides the lowest latency but the highest risk of message loss.
	RequiredAcksNone RequiredAcks = producer.RequiredAcksNone
	// RequiredAcksOne means the producer waits for an acknowledgment from the leader broker only. This provides a balance between latency and durability.
	RequiredAcksOne RequiredAcks = producer.RequiredAcksOne
	// RequiredAcksAll means the producer waits for acknowledgments from all in-sync replicas. This provides the highest level of durability but also the highest latency.
	RequiredAcksAll RequiredAcks = producer.RequiredAcksAll
)
