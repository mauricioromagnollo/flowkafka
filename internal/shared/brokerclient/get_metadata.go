package brokerclient

import (
	"context"
	"fmt"

	kafkago "github.com/segmentio/kafka-go"
)

func (c *brokerClient) GetMetadata(ctx context.Context) (*ClientMetadata, error) {
	resp, err := c.client.Metadata(ctx, &kafkago.MetadataRequest{
		Topics: []string{c.cfg.TopicName},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect producer to kafka broker %s: %w", c.cfg.Brokers[0], err)
	}

	return &ClientMetadata{
		Topics: func() []TopicMetadata {
			topics := make([]TopicMetadata, len(resp.Topics))
			for i, topic := range resp.Topics {
				topics[i] = TopicMetadata{
					Name:  topic.Name,
					Error: topic.Error,
				}
			}
			return topics
		}(),
	}, nil
}
