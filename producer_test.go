package flowkafka

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Producer", func() {
	var (
		mockProducer *producerStub
		wrapper      producerWrapper
		ctx          context.Context
	)

	BeforeEach(func() {
		mockProducer = new(producerStub)
		wrapper = producerWrapper{client: mockProducer}
		ctx = context.Background()
	})

	AfterEach(func() {
		mockProducer.AssertExpectations(GinkgoT())
	})

	Describe("NewProducer", func() {
		It("should return a producer when a config is provided", func() {
			producer := NewProducer(ProducerConfig{
				Brokers:   []string{"localhost:9092"},
				TopicName: "any-topic-name",
			})

			Expect(producer).ToNot(BeNil())

			wrapper, ok := producer.(producerWrapper)
			Expect(ok).To(BeTrue())
			Expect(wrapper.client).ToNot(BeNil())
		})
	})

	Describe("ValidateConnection", func() {
		It("should return nil when the client returns nil", func() {
			mockProducer.On("ValidateConnection", ctx).Return(nil)

			err := wrapper.ValidateConnection(ctx)

			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error when the client returns an error", func() {
			expectedErr := errors.New("any error reason returned by the client")
			mockProducer.On("ValidateConnection", ctx).Return(expectedErr)

			err := wrapper.ValidateConnection(ctx)

			Expect(err).To(MatchError(expectedErr))
		})
	})

	Describe("Publish", func() {
		It("should return nil when the client returns nil", func() {
			key := []byte("any-message-key")
			msg := map[string]string{"any-field": "any-value"}
			mockProducer.On("Publish", ctx, key, msg).Return(nil)

			err := wrapper.Publish(ctx, key, msg)

			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error when the client returns an error", func() {
			key := []byte("any-message-key")
			msg := map[string]string{"any-field": "any-value"}
			expectedErr := errors.New("any error reason returned by the client")
			mockProducer.On("Publish", ctx, key, msg).Return(expectedErr)

			err := wrapper.Publish(ctx, key, msg)

			Expect(err).To(MatchError(expectedErr))
		})
	})

	Describe("PublishAvro", func() {
		It("should return nil when the client returns nil", func() {
			key := []byte("any-message-key")
			msg := map[string]string{"any-field": "any-value"}
			mockProducer.On("PublishAvro", ctx, key, msg).Return(nil)

			err := wrapper.PublishAvro(ctx, key, msg)

			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error when the client returns an error", func() {
			key := []byte("any-message-key")
			msg := map[string]string{"any-field": "any-value"}
			expectedErr := errors.New("any error reason returned by the client")
			mockProducer.On("PublishAvro", ctx, key, msg).Return(expectedErr)

			err := wrapper.PublishAvro(ctx, key, msg)

			Expect(err).To(MatchError(expectedErr))
		})
	})

	Describe("Close", func() {
		It("should return nil when the client returns nil", func() {
			mockProducer.On("Close").Return(nil)

			err := wrapper.Close()

			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error when the client returns an error", func() {
			expectedErr := errors.New("any error reason returned by the client")
			mockProducer.On("Close").Return(expectedErr)

			err := wrapper.Close()

			Expect(err).To(MatchError(expectedErr))
		})
	})

	Describe("HasSchemaRegistry", func() {
		It("should return true when the client returns true", func() {
			mockProducer.On("HasSchemaRegistry").Return(true)

			Expect(wrapper.HasSchemaRegistry()).To(BeTrue())
		})

		It("should return false when the client returns false", func() {
			mockProducer.On("HasSchemaRegistry").Return(false)

			Expect(wrapper.HasSchemaRegistry()).To(BeFalse())
		})
	})
})
