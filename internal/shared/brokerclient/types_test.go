package brokerclient

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kafkago "github.com/segmentio/kafka-go"
)

var _ = Describe("ClientMetadata", func() {
	It("should return the configured topics when topics are provided", func() {
		metadata := ClientMetadata{
			Topics: []TopicMetadata{
				{Name: "any-topic"},
				{Name: "other-topic"},
			},
		}

		Expect(metadata.Topics).To(HaveLen(2))
		Expect(metadata.Topics[0].Name).To(Equal("any-topic"))
		Expect(metadata.Topics[1].Name).To(Equal("other-topic"))
	})

	It("should return nil topics when the metadata is a zero value", func() {
		var metadata ClientMetadata

		Expect(metadata.Topics).To(BeNil())
	})
})

var _ = Describe("TopicMetadata", func() {
	It("should return the configured name and error when both are provided", func() {
		expectedErr := errors.New("any error reason returned by the broker")
		topic := TopicMetadata{
			Name:  "any-topic",
			Error: expectedErr,
		}

		Expect(topic.Name).To(Equal("any-topic"))
		Expect(topic.Error).To(MatchError(expectedErr))
	})

	It("should return an empty name and a nil error when the topic is a zero value", func() {
		var topic TopicMetadata

		Expect(topic.Name).To(BeEmpty())
		Expect(topic.Error).To(BeNil())
	})
})

var _ = Describe("Config", func() {
	It("should return the configured values when all fields are provided", func() {
		transport := &kafkago.Transport{}
		cfg := Config{
			Brokers:   []string{"localhost:9092"},
			TopicName: "any-topic",
			Timeout:   5 * time.Second,
			Transport: transport,
		}

		Expect(cfg.Brokers).To(Equal([]string{"localhost:9092"}))
		Expect(cfg.TopicName).To(Equal("any-topic"))
		Expect(cfg.Timeout).To(Equal(5 * time.Second))
		Expect(cfg.Transport).To(BeIdenticalTo(transport))
	})

	// The `default:"10s"` struct tag on Timeout is documentation only: nothing
	// in the package reads struct tags at runtime, so the zero value stays 0
	// and kafkago.Client applies no timeout.
	It("should return zero values when the config is a zero value", func() {
		var cfg Config

		Expect(cfg.Brokers).To(BeNil())
		Expect(cfg.TopicName).To(BeEmpty())
		Expect(cfg.Timeout).To(Equal(time.Duration(0)))
		Expect(cfg.Transport).To(BeNil())
	})

	It("should not apply the documented default timeout when timeout is not provided", func() {
		client, ok := NewBrokerClient(Config{Brokers: []string{"localhost:9092"}}).(*brokerClient)

		Expect(ok).To(BeTrue())
		Expect(client.client.Timeout).To(Equal(time.Duration(0)))
	})
})
