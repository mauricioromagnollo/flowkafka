package brokerclient

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Stub is a mock implementation of the BrokerClient interface for testing purposes.
type Stub struct {
	BrokerClient
	mock.Mock
}

// GetMetadata records the call and returns the metadata and error configured on the stub.
func (m *Stub) GetMetadata(ctx context.Context) (*ClientMetadata, error) {
	args := m.Called(ctx)

	return args.Get(0).(*ClientMetadata), args.Error(1)
}
