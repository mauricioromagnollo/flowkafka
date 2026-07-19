package flowkafka

import (
	"context"

	"github.com/mauricioromagnollo/flowkafka/internal/consumer"
	"github.com/mauricioromagnollo/flowkafka/internal/shared/types"
	"github.com/stretchr/testify/mock"
)

// consumerStub is a mock implementation of consumer.Consumer used to assert that
// consumerWrapper delegates every call to the internal client. It lives in a
// _test.go file so the published package does not depend on testify/mock.
type consumerStub struct {
	consumer.Consumer
	mock.Mock
}

func (m *consumerStub) Consume(ctx context.Context, handler func(msg types.Message) error) error {
	args := m.Called(ctx, handler)

	return args.Error(0)
}

func (m *consumerStub) ConsumeMessages(ctx context.Context, msgsChan chan<- types.Message) error {
	args := m.Called(ctx, msgsChan)

	return args.Error(0)
}

func (m *consumerStub) Close() error {
	args := m.Called()

	return args.Error(0)
}

func (m *consumerStub) ValidateConnection(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}
