package flowkafka

import (
	"context"

	"github.com/mauricioromagnollo/flowkafka/internal/consumer"
)

// Consumer is an interface that defines the methods for interacting with a Kafka consumer.
type Consumer = consumer.Consumer

// ConsumerConfig holds the configuration options for the Kafka consumer.
type ConsumerConfig = consumer.Config

type consumerWrapper struct {
	client consumer.Consumer
}

// NewConsumer creates a new Kafka consumer with the given configuration.
func NewConsumer(config ConsumerConfig) Consumer {
	return consumerWrapper{
		client: consumer.NewConsumer(config),
	}
}

func (c consumerWrapper) Consume(ctx context.Context, handler func(msg Message) error) error {
	return c.client.Consume(ctx, handler)
}

func (c consumerWrapper) ConsumeMessages(ctx context.Context, msgsChan chan<- Message) error {
	return c.client.ConsumeMessages(ctx, msgsChan)
}

func (c consumerWrapper) Close() error {
	return c.client.Close()
}

func (c consumerWrapper) ValidateConnection(ctx context.Context) error {
	return c.client.ValidateConnection(ctx)
}
