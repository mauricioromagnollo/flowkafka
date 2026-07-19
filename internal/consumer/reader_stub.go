package consumer

import (
	"context"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
)

type readerStub struct {
	Reader
	mock.Mock
}

func (m *readerStub) FetchMessage(ctx context.Context) (kafkago.Message, error) {
	args := m.Called(ctx)

	return args.Get(0).(kafkago.Message), args.Error(1)
}

func (m *readerStub) CommitMessages(ctx context.Context, msgs ...kafkago.Message) error {
	args := m.Called(ctx, msgs)

	return args.Error(0)
}

func (m *readerStub) Close() error {
	args := m.Called()

	return args.Error(0)
}
