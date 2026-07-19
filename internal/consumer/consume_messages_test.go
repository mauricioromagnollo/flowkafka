package consumer

import (
	"context"
	"errors"

	"github.com/mauricioromagnollo/flowkafka/internal/shared/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("ConsumeMessages", func() {
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
		msgsChan := make(chan types.Message, 1)

		mockReader.On("FetchMessage", ctx).
			Return(kafkago.Message{}, errors.New("any error reason returned by the reader")).Once()

		err := consumer.ConsumeMessages(ctx, msgsChan)

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
		msgsChan := make(chan types.Message, 1)

		mockReader.On("FetchMessage", ctx).
			Return(kafkago.Message{}, context.Canceled).Once()

		err := consumer.ConsumeMessages(ctx, msgsChan)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("context error:"))
		Expect(errors.Is(err, context.Canceled)).To(BeTrue())
	})

	It("should send the mapped message to the channel when message is fetched successfully", func() {
		consumer := &consumerClient{
			reader: mockReader,
			cfg:    Config{TopicName: "any-topic"},
		}
		ctx := context.Background()
		msgsChan := make(chan types.Message, 1)
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

		err := consumer.ConsumeMessages(ctx, msgsChan)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("stop the loop"))
		Expect(msgsChan).To(HaveLen(1))

		received := <-msgsChan
		Expect(received.Key).To(Equal([]byte("any-key")))
		Expect(received.Value).To(Equal([]byte("any-value")))
		Expect(received.Headers).To(Equal(headers))
		Expect(received.Partition).To(Equal(3))
		Expect(received.Offset).To(Equal(int64(42)))
	})

	It("should return nil when context is canceled while sending the message to the channel", func() {
		consumer := &consumerClient{
			reader: mockReader,
			cfg:    Config{TopicName: "any-topic"},
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		msgsChan := make(chan types.Message)

		mockReader.On("FetchMessage", mock.Anything).
			Return(kafkago.Message{Value: []byte("any-value")}, nil).
			Run(func(args mock.Arguments) { cancel() }).
			Once()

		errChan := make(chan error, 1)
		go func() {
			defer GinkgoRecover()
			errChan <- consumer.ConsumeMessages(ctx, msgsChan)
		}()

		Eventually(errChan).Should(Receive(BeNil()))
	})

	It("should return error when fail to commit message", func() {
		consumer := &consumerClient{
			reader: mockReader,
			cfg:    Config{TopicName: "any-topic"},
		}
		ctx := context.Background()
		msgsChan := make(chan types.Message, 1)
		kafkaMsg := kafkago.Message{Value: []byte("any-value")}

		mockReader.On("FetchMessage", ctx).Return(kafkaMsg, nil).Once()
		mockReader.On("CommitMessages", ctx, []kafkago.Message{kafkaMsg}).
			Return(errors.New("any error reason returned by the reader")).Once()

		err := consumer.ConsumeMessages(ctx, msgsChan)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to commit message:"))
		Expect(err.Error()).To(ContainSubstring("any error reason returned by the reader"))
	})
})
