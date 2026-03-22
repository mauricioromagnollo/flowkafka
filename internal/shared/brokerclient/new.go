package brokerclient

import (
	kafkago "github.com/segmentio/kafka-go"
)

type brokerClient struct {
	cfg    Config
	client *kafkago.Client
}

// NewBrokerClient creates a new Kafka broker client with the given configuration.
func NewBrokerClient(cfg Config) BrokerClient {
	client := &kafkago.Client{
		Addr:      kafkago.TCP(cfg.Brokers[0]),
		Transport: cfg.Transport,
		Timeout:   cfg.Timeout,
	}

	return &brokerClient{
		cfg:    cfg,
		client: client,
	}
}
