package consumer

import (
	"context"
	"errors"

	"github.com/mauricioromagnollo/flowkafka/internal/shared/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kafkago "github.com/segmentio/kafka-go"
)

var _ = Describe("Consume", func() {
	var mockReader *readerStub

	BeforeEach(func() {
		mockReader = new(readerStub)
	})

	AfterEach(func() {
		mockReader.AssertExpectations(GinkgoT())
	})

	It("should return error when fail to fetch message and context is still alive", func() {
		consumer := &consumerClient{
			reader: mockReader,
			cfg:    Config{TopicName: "any-topic"},
		}
		ctx := context.Background()

		mockReader.On("FetchMessage", ctx).
			Return(kafkago.Message{}, errors.New("any error reason returned by the reader")).Once()

		err := consumer.Consume(ctx, func(msg types.Message) error { return nil })

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to fetch message from topic any-topic"))
		Expect(err.Error()).To(ContainSubstring("any error reason returned by the reader"))
	})

	It("should return context error when fail to fetch message and context is canceled", func() {
		consumer := &consumerClient{
			reader: mockReader,
			cfg:    Config{TopicName: "any-topic"},
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		mockReader.On("FetchMessage", ctx).
			Return(kafkago.Message{}, context.Canceled).Once()

		err := consumer.Consume(ctx, func(msg types.Message) error { return nil })

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("context error:"))
		Expect(errors.Is(err, context.Canceled)).To(BeTrue())
	})

	It("should return error and not commit when handler returns an error", func() {
		consumer := &consumerClient{
			reader: mockReader,
			cfg:    Config{TopicName: "any-topic"},
		}
		ctx := context.Background()

		mockReader.On("FetchMessage", ctx).
			Return(kafkago.Message{Value: []byte("any-value")}, nil).Once()

		err := consumer.Consume(ctx, func(msg types.Message) error {
			return errors.New("any error reason returned by the handler")
		})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("message handler error:"))
		Expect(err.Error()).To(ContainSubstring("any error reason returned by the handler"))
		mockReader.AssertNotCalled(GinkgoT(), "CommitMessages")
	})

	It("should return error when fail to commit message", func() {
		consumer := &consumerClient{
			reader: mockReader,
			cfg:    Config{TopicName: "any-topic"},
		}
		ctx := context.Background()
		kafkaMsg := kafkago.Message{Value: []byte("any-value")}

		mockReader.On("FetchMessage", ctx).Return(kafkaMsg, nil).Once()
		mockReader.On("CommitMessages", ctx, []kafkago.Message{kafkaMsg}).
			Return(errors.New("any error reason returned by the reader")).Once()

		err := consumer.Consume(ctx, func(msg types.Message) error { return nil })

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to commit message:"))
		Expect(err.Error()).To(ContainSubstring("any error reason returned by the reader"))
	})

	It("should deliver the mapped message to the handler and commit it when message is fetched successfully", func() {
		consumer := &consumerClient{
			reader: mockReader,
			cfg:    Config{TopicName: "any-topic"},
		}
		ctx := context.Background()
		headers := []kafkago.Header{{Key: "any-header", Value: []byte("any-header-value")}}
		kafkaMsg := kafkago.Message{
			Key:       []byte("any-key"),
			Value:     []byte("any-value"),
			Headers:   headers,
			Partition: 3,
			Offset:    42,
		}

		mockReader.On("FetchMessage", ctx).Return(kafkaMsg, nil).Once()
		mockReader.On("CommitMessages", ctx, []kafkago.Message{kafkaMsg}).Return(nil).Once()
		mockReader.On("FetchMessage", ctx).
			Return(kafkago.Message{}, errors.New("stop the loop")).Once()

		var received []types.Message
		err := consumer.Consume(ctx, func(msg types.Message) error {
			received = append(received, msg)
			return nil
		})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("stop the loop"))
		Expect(received).To(HaveLen(1))
		Expect(received[0].Key).To(Equal([]byte("any-key")))
		Expect(received[0].Value).To(Equal([]byte("any-value")))
		Expect(received[0].Headers).To(Equal(headers))
		Expect(received[0].Partition).To(Equal(3))
		Expect(received[0].Offset).To(Equal(int64(42)))
	})
})
