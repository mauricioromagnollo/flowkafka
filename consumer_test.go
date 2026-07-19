package flowkafka

import (
	"context"
	"errors"

	"github.com/mauricioromagnollo/flowkafka/internal/shared/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Consumer", func() {
	var (
		mockConsumer *consumerStub
		wrapper      consumerWrapper
		ctx          context.Context
	)

	BeforeEach(func() {
		mockConsumer = new(consumerStub)
		wrapper = consumerWrapper{client: mockConsumer}
		ctx = context.Background()
	})

	AfterEach(func() {
		mockConsumer.AssertExpectations(GinkgoT())
	})

	Describe("NewConsumer", func() {
		It("should return a consumer when a config is provided", func() {
			consumer := NewConsumer(ConsumerConfig{
				Brokers:   []string{"localhost:9092"},
				GroupID:   "any-group-id",
				TopicName: "any-topic-name",
			})

			Expect(consumer).ToNot(BeNil())

			newWrapper, ok := consumer.(consumerWrapper)
			Expect(ok).To(BeTrue())
			Expect(newWrapper.client).ToNot(BeNil())
		})
	})

	Describe("Consume", func() {
		It("should return nil when the client returns nil", func() {
			handler := func(_ types.Message) error { return nil }

			var receivedHandler func(msg types.Message) error
			mockConsumer.On("Consume", ctx, mock.Anything).
				Run(func(args mock.Arguments) {
					receivedHandler = args.Get(1).(func(msg types.Message) error)
				}).
				Return(nil)

			err := wrapper.Consume(ctx, handler)

			Expect(err).ToNot(HaveOccurred())
			Expect(receivedHandler).ToNot(BeNil())
			Expect(receivedHandler(types.Message{})).To(BeNil())
		})

		It("should return error when the client returns an error", func() {
			expectedErr := errors.New("any error reason returned by the client")
			mockConsumer.On("Consume", ctx, mock.Anything).Return(expectedErr)

			err := wrapper.Consume(ctx, func(_ types.Message) error { return nil })

			Expect(err).To(MatchError(expectedErr))
		})
	})

	Describe("ConsumeMessages", func() {
		It("should return nil when the client returns nil", func() {
			msgsChan := make(chan types.Message, 1)
			defer close(msgsChan)

			var receivedChan chan<- types.Message
			mockConsumer.On("ConsumeMessages", ctx, mock.Anything).
				Run(func(args mock.Arguments) {
					receivedChan = args.Get(1).(chan<- types.Message)
				}).
				Return(nil)

			err := wrapper.ConsumeMessages(ctx, msgsChan)

			Expect(err).ToNot(HaveOccurred())
			Expect(receivedChan).ToNot(BeNil())

			receivedChan <- types.Message{Key: []byte("any-message-key")}
			Expect(<-msgsChan).To(Equal(types.Message{Key: []byte("any-message-key")}))
		})

		It("should return error when the client returns an error", func() {
			msgsChan := make(chan types.Message)
			defer close(msgsChan)

			expectedErr := errors.New("any error reason returned by the client")
			mockConsumer.On("ConsumeMessages", ctx, mock.Anything).Return(expectedErr)

			err := wrapper.ConsumeMessages(ctx, msgsChan)

			Expect(err).To(MatchError(expectedErr))
		})
	})

	Describe("Close", func() {
		It("should return nil when the client returns nil", func() {
			mockConsumer.On("Close").Return(nil)

			err := wrapper.Close()

			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error when the client returns an error", func() {
			expectedErr := errors.New("any error reason returned by the client")
			mockConsumer.On("Close").Return(expectedErr)

			err := wrapper.Close()

			Expect(err).To(MatchError(expectedErr))
		})
	})

	Describe("ValidateConnection", func() {
		It("should return nil when the client returns nil", func() {
			mockConsumer.On("ValidateConnection", ctx).Return(nil)

			err := wrapper.ValidateConnection(ctx)

			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error when the client returns an error", func() {
			expectedErr := errors.New("any error reason returned by the client")
			mockConsumer.On("ValidateConnection", ctx).Return(expectedErr)

			err := wrapper.ValidateConnection(ctx)

			Expect(err).To(MatchError(expectedErr))
		})
	})
})
