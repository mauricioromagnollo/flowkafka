package kafka

import (
	"context"
	"fmt"
)

// ValidateConnection validates the connection to Kafka by attempting to retrieve the list of brokers.
func (c *Client) ValidateConnection(ctx context.Context) error {
	conn, err := c.dialer.DialContext(ctx, "tcp", c.config.Addresses[0])
	if err != nil {
		return fmt.Errorf("failed to connect to kafka: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	_, err = conn.Brokers()
	if err != nil {
		return fmt.Errorf("failed to get kafka brokers: %w", err)
	}

	return nil
}
