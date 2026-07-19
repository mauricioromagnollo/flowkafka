package consumer

import (
	"context"
	"errors"

	"github.com/mauricioromagnollo/flowkafka/internal/shared/brokerclient"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidateConnection", func() {
	var mockBrokerClient *brokerclient.Stub

	BeforeEach(func() {
		mockBrokerClient = new(brokerclient.Stub)
	})

	AfterEach(func() {
		mockBrokerClient.AssertExpectations(GinkgoT())
	})

	It("should return error when fail to get metadata from the broker", func() {
		consumer := &consumerClient{
			cfg:          Config{TopicName: "any-topic"},
			brokerClient: mockBrokerClient,
		}
		ctx := context.Background()

		mockBrokerClient.On("GetMetadata", ctx).
			Return((*brokerclient.ClientMetadata)(nil), errors.New("any error reason returned by the broker"))

		err := consumer.ValidateConnection(ctx)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("consumer failed to connect to kafka broker"))
		Expect(err.Error()).To(ContainSubstring("any error reason returned by the broker"))
	})

	It("should return error when the topic is found with an error", func() {
		consumer := &consumerClient{
			cfg:          Config{TopicName: "any-topic"},
			brokerClient: mockBrokerClient,
		}
		ctx := context.Background()

		mockBrokerClient.On("GetMetadata", ctx).Return(&brokerclient.ClientMetadata{
			Topics: []brokerclient.TopicMetadata{
				{Name: "any-topic", Error: errors.New("any error reason returned by the topic")},
			},
		}, nil)

		err := consumer.ValidateConnection(ctx)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("consumer failed to find topic any-topic"))
		Expect(err.Error()).To(ContainSubstring("any error reason returned by the topic"))
	})

	It("should return nil when the topic is found without errors", func() {
		consumer := &consumerClient{
			cfg:          Config{TopicName: "any-topic"},
			brokerClient: mockBrokerClient,
		}
		ctx := context.Background()

		mockBrokerClient.On("GetMetadata", ctx).Return(&brokerclient.ClientMetadata{
			Topics: []brokerclient.TopicMetadata{
				{Name: "other-topic", Error: nil},
				{Name: "any-topic", Error: nil},
			},
		}, nil)

		err := consumer.ValidateConnection(ctx)

		Expect(err).To(BeNil())
	})

	It("should return error when the topic is not present in the metadata", func() {
		consumer := &consumerClient{
			cfg:          Config{TopicName: "any-topic"},
			brokerClient: mockBrokerClient,
		}
		ctx := context.Background()

		mockBrokerClient.On("GetMetadata", ctx).Return(&brokerclient.ClientMetadata{
			Topics: []brokerclient.TopicMetadata{
				{Name: "other-topic", Error: nil},
			},
		}, nil)

		err := consumer.ValidateConnection(ctx)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("consumer failed to find topic any-topic"))
	})

	It("should return error when the metadata has an empty topic list", func() {
		consumer := &consumerClient{
			cfg:          Config{TopicName: "any-topic"},
			brokerClient: mockBrokerClient,
		}
		ctx := context.Background()

		mockBrokerClient.On("GetMetadata", ctx).
			Return(&brokerclient.ClientMetadata{Topics: []brokerclient.TopicMetadata{}}, nil)

		err := consumer.ValidateConnection(ctx)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("consumer failed to find topic any-topic"))
	})
})
