package kafka

import (
	"context"
)

// ReadMessage reads a message from the specified Kafka topic, returning the message along with its key, value, headers, partition, and offset.
func (c *Client) ReadMessage(ctx context.Context, topic string) (Message, error) {
	r := c.getOrCreateReader(topic)

	m, err := r.FetchMessage(ctx)
	if err != nil {
		return Message{}, err
	}

	headers := make(map[string][]byte, len(m.Headers))
	for _, h := range m.Headers {
		headers[h.Key] = h.Value
	}

	return Message{
		Key:       m.Key,
		Value:     m.Value,
		Headers:   headers,
		Partition: m.Partition,
		Offset:    m.Offset,
	}, nil
}
