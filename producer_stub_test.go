package flowkafka

import (
	"context"

	"github.com/mauricioromagnollo/flowkafka/internal/producer"
	"github.com/stretchr/testify/mock"
)

// producerStub is a mock implementation of producer.Producer used to assert that
// producerWrapper delegates every call to the internal client. It lives in a
// _test.go file so the published package does not depend on testify/mock.
type producerStub struct {
	producer.Producer
	mock.Mock
}

func (m *producerStub) ValidateConnection(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

func (m *producerStub) Publish(ctx context.Context, key []byte, msg any) error {
	args := m.Called(ctx, key, msg)

	return args.Error(0)
}

func (m *producerStub) PublishAvro(ctx context.Context, key []byte, msg any) error {
	args := m.Called(ctx, key, msg)

	return args.Error(0)
}

func (m *producerStub) Close() error {
	args := m.Called()

	return args.Error(0)
}

func (m *producerStub) HasSchemaRegistry() bool {
	args := m.Called()

	return args.Bool(0)
}
