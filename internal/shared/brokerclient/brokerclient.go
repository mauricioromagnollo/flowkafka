package brokerclient

import "context"

// BrokerClient defines the interface for interacting with Kafka brokers to fetch metadata and validate connections.
type BrokerClient interface {
	// GetMetadata retrieves metadata about the Kafka cluster, including topics and their statuses.
	GetMetadata(ctx context.Context) (*ClientMetadata, error)
}
