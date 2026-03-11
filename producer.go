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

func (p producerWrapper) ProduceJSONMessage(ctx context.Context, key []byte, msg any) error {
	return p.client.ProduceJSONMessage(ctx, key, msg)
}

func (p producerWrapper) ProduceAvroMessage(ctx context.Context, key []byte, msg any) error {
	return p.client.ProduceAvroMessage(ctx, key, msg)
}

func (p producerWrapper) Close() error {
	return p.client.Close()
}

func (p producerWrapper) HasSchemaRegistry() bool {
	return p.client.HasSchemaRegistry()
}
