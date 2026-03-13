package producer

import (
	"context"
	"fmt"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

// ValidateConnection tests connectivity to Kafka brokers.
func (c *producerClient) ValidateConnection(ctx context.Context) error {
	transport := buildTransport(c.cfg)

	client := &kafkago.Client{
		Addr:      kafkago.TCP(c.cfg.Brokers[0]),
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Metadata(ctx, &kafkago.MetadataRequest{
		Topics: []string{c.cfg.TopicName},
	})
	if err != nil {
		return fmt.Errorf("failed to connect producer to kafka broker %s: %w", c.cfg.Brokers[0], err)
	}

	for _, topic := range resp.Topics {
		if topic.Name == c.cfg.TopicName && topic.Error == nil {
			return nil
		}
	}

	return fmt.Errorf("producer failed to find topic %s", c.cfg.TopicName)
}
