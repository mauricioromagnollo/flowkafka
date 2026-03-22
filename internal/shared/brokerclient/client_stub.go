package brokerclient

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// BrokerClientStub is a mock implementation of the BrokerClient interface for testing purposes.
type BrokerClientStub struct {
	BrokerClient
	mock.Mock
}

func (m *BrokerClientStub) GetMetadata(ctx context.Context) (*ClientMetadata, error) {
	args := m.Called(ctx)

	return args.Get(0).(*ClientMetadata), args.Error(1)
}
