package producer

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Publish", func() {
	var mockWriter *writerStub

	BeforeEach(func() {
		mockWriter = new(writerStub)
	})

	AfterEach(func() {
		mockWriter.AssertExpectations(GinkgoT())
	})

	It("should return error when fail to marshal message", func() {
		config := Config{}
		producer := &producerClient{
			writer:         mockWriter,
			cfg:            config,
			schemaRegistry: nil,
		}
		unmarshalableMsg := make(chan int)
		ctx := context.Background()

		err := producer.Publish(ctx, nil, unmarshalableMsg)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to marshal message to publish"))
	})

	It("should return error when fail to write message", func() {
		config := Config{}
		producer := &producerClient{
			writer:         mockWriter,
			cfg:            config,
			schemaRegistry: nil,
		}
		msg := map[string]string{"key": "value"}
		ctx := context.Background()

		mockWriter.On("WriteMessages", ctx, Message{
			Value: []byte(`{"key":"value"}`),
			Key:   nil,
		}).Return(fmt.Errorf("any error reason returned by the writer"))

		err := producer.Publish(ctx, nil, msg)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to publish message in kafka: any error reason returned by the writer"))
	})

	It("should return nil when message is published successfully", func() {
		config := Config{}
		producer := &producerClient{
			writer:         mockWriter,
			cfg:            config,
			schemaRegistry: nil,
		}
		msg := map[string]string{"key": "value"}
		msgKey := []byte("message-key")
		ctx := context.Background()

		mockWriter.On("WriteMessages", ctx, Message{
			Value: []byte(`{"key":"value"}`),
			Key:   msgKey,
		}).Return(nil)

		err := producer.Publish(ctx, msgKey, msg)

		Expect(err).ToNot(HaveOccurred())
	})
})
