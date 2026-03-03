package kafka

import (
	kafkago "github.com/segmentio/kafka-go"
)

// getOrCreateReader retrieves an existing Kafka reader for the specified topic or creates a new one if it doesn't exist.
// This method is thread-safe and uses a mutex to prevent concurrent map access.
// The reader is configured with manual commit mode (CommitInterval: 0), which disables auto-commit.
// This means messages must be explicitly committed after successful processing by calling CommitMessages.
func (c *Client) getOrCreateReader(topic string) *kafkago.Reader {
	c.mu.Lock()
	defer c.mu.Unlock()

	if r, ok := c.readers[topic]; ok {
		return r
	}

	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        c.config.Addresses,
		Topic:          topic,
		GroupID:        c.config.GroupID,
		Dialer:         c.dialer,
		MinBytes:       1e3,
		MaxBytes:       10e6,
		CommitInterval: 0, // 0 disables auto-commit, messages must be committed manually after processing
	})

	c.readers[topic] = r
	return r
}
