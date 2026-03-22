package producer

import (
	"context"

	"github.com/mauricioromagnollo/flowkafka/internal/shared/types"
	"github.com/stretchr/testify/mock"
)

type writerStub struct {
	Writer
	mock.Mock
}

func (m *writerStub) WriteMessages(ctx context.Context, msg types.Message) error {
	args := m.Called(ctx, msg)

	return args.Error(0)
}

func (m *writerStub) Close() error {
	args := m.Called()

	return args.Error(0)
}
