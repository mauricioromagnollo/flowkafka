package producer

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

	It("should return error when broker client fails to get metadata", func() {
		producer := &producerClient{
			cfg:          Config{TopicName: "any-topic-name"},
			brokerClient: mockBrokerClient,
		}
		ctx := context.Background()

		mockBrokerClient.On("GetMetadata", ctx).
			Return((*brokerclient.ClientMetadata)(nil), errors.New("any error reason returned by the broker client"))

		err := producer.ValidateConnection(ctx)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("producer failed to connect to kafka broker"))
		Expect(err.Error()).To(ContainSubstring("any error reason returned by the broker client"))
	})

	It("should return error when the topic is found with an error", func() {
		producer := &producerClient{
			cfg:          Config{TopicName: "any-topic-name"},
			brokerClient: mockBrokerClient,
		}
		ctx := context.Background()

		mockBrokerClient.On("GetMetadata", ctx).Return(&brokerclient.ClientMetadata{
			Topics: []brokerclient.TopicMetadata{
				{Name: "any-topic-name", Error: errors.New("any error reason returned by the topic")},
			},
		}, nil)

		err := producer.ValidateConnection(ctx)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("producer failed to find topic any-topic-name"))
		Expect(err.Error()).To(ContainSubstring("any error reason returned by the topic"))
	})

	It("should return nil when the topic is found without an error", func() {
		producer := &producerClient{
			cfg:          Config{TopicName: "any-topic-name"},
			brokerClient: mockBrokerClient,
		}
		ctx := context.Background()

		mockBrokerClient.On("GetMetadata", ctx).Return(&brokerclient.ClientMetadata{
			Topics: []brokerclient.TopicMetadata{
				{Name: "any-other-topic-name", Error: nil},
				{Name: "any-topic-name", Error: nil},
			},
		}, nil)

		err := producer.ValidateConnection(ctx)

		Expect(err).ToNot(HaveOccurred())
	})

	It("should return error when the topic is not present in the metadata", func() {
		producer := &producerClient{
			cfg:          Config{TopicName: "any-topic-name"},
			brokerClient: mockBrokerClient,
		}
		ctx := context.Background()

		mockBrokerClient.On("GetMetadata", ctx).Return(&brokerclient.ClientMetadata{
			Topics: []brokerclient.TopicMetadata{
				{Name: "any-other-topic-name", Error: nil},
			},
		}, nil)

		err := producer.ValidateConnection(ctx)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("producer failed to find topic any-topic-name"))
	})

	It("should return error when the metadata has an empty topic list", func() {
		producer := &producerClient{
			cfg:          Config{TopicName: "any-topic-name"},
			brokerClient: mockBrokerClient,
		}
		ctx := context.Background()

		mockBrokerClient.On("GetMetadata", ctx).Return(&brokerclient.ClientMetadata{
			Topics: []brokerclient.TopicMetadata{},
		}, nil)

		err := producer.ValidateConnection(ctx)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("producer failed to find topic any-topic-name"))
	})
})
