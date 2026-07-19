package brokerclient

import (
	"context"
	"net"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kafkago "github.com/segmentio/kafka-go"
	metadataAPI "github.com/segmentio/kafka-go/protocol/metadata"
	"github.com/stretchr/testify/mock"
)

// roundTripperStub is a mock implementation of kafkago.RoundTripper, used to
// exercise GetMetadata without talking to a real Kafka broker.
type roundTripperStub struct {
	mock.Mock
}

func (m *roundTripperStub) RoundTrip(ctx context.Context, addr net.Addr, req kafkago.Request) (kafkago.Response, error) {
	args := m.Called(ctx, addr, req)
	resp, _ := args.Get(0).(kafkago.Response)

	return resp, args.Error(1)
}

var _ = Describe("GetMetadata", func() {
	var mockRoundTripper *roundTripperStub

	newClientWithRoundTripper := func(cfg Config) *brokerClient {
		return &brokerClient{
			cfg: cfg,
			client: &kafkago.Client{
				Addr:      kafkago.TCP("localhost:9092"),
				Timeout:   500 * time.Millisecond,
				Transport: mockRoundTripper,
			},
		}
	}

	BeforeEach(func() {
		mockRoundTripper = new(roundTripperStub)
	})

	It("should return empty topics when the broker returns no topics", func() {
		client := newClientWithRoundTripper(Config{
			Brokers:   []string{"localhost:9092"},
			TopicName: "any-topic",
		})
		ctx := context.Background()

		mockRoundTripper.
			On("RoundTrip", mock.Anything, mock.Anything, mock.Anything).
			Return(&metadataAPI.Response{}, nil)

		metadata, err := client.GetMetadata(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(metadata).ToNot(BeNil())
		Expect(metadata.Topics).To(BeEmpty())
		mockRoundTripper.AssertExpectations(GinkgoT())
	})

	It("should return the mapped topics when the broker returns topics without errors", func() {
		client := newClientWithRoundTripper(Config{
			Brokers:   []string{"localhost:9092"},
			TopicName: "any-topic",
		})
		ctx := context.Background()

		mockRoundTripper.
			On("RoundTrip", mock.Anything, mock.Anything, mock.Anything).
			Return(&metadataAPI.Response{
				Topics: []metadataAPI.ResponseTopic{
					{Name: "any-topic"},
					{Name: "other-topic"},
				},
			}, nil)

		metadata, err := client.GetMetadata(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(metadata.Topics).To(HaveLen(2))
		Expect(metadata.Topics[0].Name).To(Equal("any-topic"))
		Expect(metadata.Topics[0].Error).To(BeNil())
		Expect(metadata.Topics[1].Name).To(Equal("other-topic"))
		Expect(metadata.Topics[1].Error).To(BeNil())
	})

	It("should return the topic error when the broker returns a topic with an error code", func() {
		client := newClientWithRoundTripper(Config{
			Brokers:   []string{"localhost:9092"},
			TopicName: "unknown-topic",
		})
		ctx := context.Background()

		mockRoundTripper.
			On("RoundTrip", mock.Anything, mock.Anything, mock.Anything).
			Return(&metadataAPI.Response{
				Topics: []metadataAPI.ResponseTopic{
					{
						Name:      "unknown-topic",
						ErrorCode: int16(kafkago.UnknownTopicOrPartition),
					},
				},
			}, nil)

		metadata, err := client.GetMetadata(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(metadata.Topics).To(HaveLen(1))
		Expect(metadata.Topics[0].Name).To(Equal("unknown-topic"))
		Expect(metadata.Topics[0].Error).To(MatchError(kafkago.UnknownTopicOrPartition))
	})

	It("should request the configured topic name when metadata is fetched", func() {
		client := newClientWithRoundTripper(Config{
			Brokers:   []string{"localhost:9092"},
			TopicName: "any-topic",
		})
		ctx := context.Background()

		mockRoundTripper.
			On("RoundTrip", mock.Anything, mock.Anything, mock.MatchedBy(func(req kafkago.Request) bool {
				metadataRequest, ok := req.(*metadataAPI.Request)

				return ok && len(metadataRequest.TopicNames) == 1 && metadataRequest.TopicNames[0] == "any-topic"
			})).
			Return(&metadataAPI.Response{}, nil)

		_, err := client.GetMetadata(ctx)

		Expect(err).ToNot(HaveOccurred())
		mockRoundTripper.AssertExpectations(GinkgoT())
	})

	It("should return an error containing the broker addresses when the broker is unreachable", func() {
		client, ok := NewBrokerClient(Config{
			Brokers:   []string{"127.0.0.1:1"},
			TopicName: "any-topic",
			Timeout:   500 * time.Millisecond,
		}).(*brokerClient)
		Expect(ok).To(BeTrue())

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		metadata, err := client.GetMetadata(ctx)

		Expect(err).To(HaveOccurred())
		Expect(metadata).To(BeNil())
		Expect(err.Error()).To(ContainSubstring("failed to connect producer to kafka broker"))
		Expect(err.Error()).To(ContainSubstring("127.0.0.1:1"))
	})

	It("should return an error listing every broker when multiple brokers are unreachable", func() {
		client, ok := NewBrokerClient(Config{
			Brokers:   []string{"127.0.0.1:1", "127.0.0.1:2"},
			TopicName: "any-topic",
			Timeout:   500 * time.Millisecond,
		}).(*brokerClient)
		Expect(ok).To(BeTrue())

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := client.GetMetadata(ctx)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("127.0.0.1:1,127.0.0.1:2"))
	})

	// Regression: the error message used to read cfg.Brokers[0], which panicked
	// when no broker was configured.
	It("should return an error without panicking when brokers is empty", func() {
		client, ok := NewBrokerClient(Config{
			TopicName: "any-topic",
			Timeout:   500 * time.Millisecond,
		}).(*brokerClient)
		Expect(ok).To(BeTrue())

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		var metadata *ClientMetadata
		var err error

		Expect(func() { metadata, err = client.GetMetadata(ctx) }).ToNot(Panic())
		Expect(err).To(HaveOccurred())
		Expect(metadata).To(BeNil())
		Expect(err.Error()).To(ContainSubstring("failed to connect producer to kafka broker"))
	})

	It("should return an error when the round tripper fails", func() {
		client := newClientWithRoundTripper(Config{
			Brokers:   []string{"localhost:9092"},
			TopicName: "any-topic",
		})
		ctx := context.Background()

		mockRoundTripper.
			On("RoundTrip", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, context.DeadlineExceeded)

		metadata, err := client.GetMetadata(ctx)

		Expect(err).To(HaveOccurred())
		Expect(metadata).To(BeNil())
		Expect(err.Error()).To(ContainSubstring("failed to connect producer to kafka broker localhost:9092"))
	})
})
