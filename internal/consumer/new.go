package consumer

import (
	kafkago "github.com/segmentio/kafka-go"
)

// consumerClient wraps Kafka reading with consumer group and manual commit support.
type consumerClient struct {
	reader *kafkago.Reader
	cfg    Config
}

// NewConsumer creates a new Kafka consumer with the given configuration.
func NewConsumer(config Config) Consumer {
	dialer := buildDialer(config)

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: config.Brokers,
		GroupID: config.GroupID,
		Topic:   config.TopicName,
		Dialer:  dialer,
	})

	return &consumerClient{
		reader: reader,
		cfg:    config,
	}
}
