package kafka

import (
	"context"

	kafkago "github.com/segmentio/kafka-go"
)

// CommitMessage commits the given message for the specified topic, marking it as processed.
func (c *Client) CommitMessage(ctx context.Context, topic string, msg Message) error {
	r := c.getOrCreateReader(topic)

	return r.CommitMessages(ctx, kafkago.Message{
		Topic:     topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
		Key:       msg.Key,
		Value:     msg.Value,
	})
}
