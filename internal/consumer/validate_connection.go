package consumer

import (
	"context"
	"fmt"
)

// ValidateConnection tests connectivity to Kafka brokers.
func (c *consumerClient) ValidateConnection(ctx context.Context) error {
	// Use the client to fetch metadata from the Kafka broker to validate the connection
	metadata, err := c.brokerClient.GetMetadata(ctx)
	if err != nil {
		return fmt.Errorf("consumer failed to connect to kafka broker: %w", err)
	}

	// Check if the topic exists and has no errors
	for _, topic := range metadata.Topics {
		if topic.Name == c.cfg.TopicName {
			if topic.Error != nil {
				return fmt.Errorf("consumer failed to find topic %s: %w", c.cfg.TopicName, topic.Error)
			}
			return nil // Topic found and has no errors
		}
	}

	return fmt.Errorf("consumer failed to find topic %s", c.cfg.TopicName)
}
