package brokerclient

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kafkago "github.com/segmentio/kafka-go"
)

var _ = Describe("NewBrokerClient", func() {
	It("should return a non-nil broker client when config is empty", func() {
		client := NewBrokerClient(Config{})

		Expect(client).ToNot(BeNil())
	})

	// Regression: NewBrokerClient used to index cfg.Brokers[0] without checking
	// the slice length, which panicked with "index out of range" for an empty
	// broker list.
	It("should not panic when brokers is empty", func() {
		Expect(func() { NewBrokerClient(Config{}) }).ToNot(Panic())
	})

	It("should return a client with a nil address when brokers is empty", func() {
		client, ok := NewBrokerClient(Config{}).(*brokerClient)

		Expect(ok).To(BeTrue())
		Expect(client.client.Addr).To(BeNil())
	})

	It("should return a client with a nil address when brokers is an empty slice", func() {
		client, ok := NewBrokerClient(Config{Brokers: []string{}}).(*brokerClient)

		Expect(ok).To(BeTrue())
		Expect(client.client.Addr).To(BeNil())
	})

	It("should return a client with the broker address when a single broker is provided", func() {
		client, ok := NewBrokerClient(Config{Brokers: []string{"localhost:9092"}}).(*brokerClient)

		Expect(ok).To(BeTrue())
		Expect(client.client.Addr).ToNot(BeNil())
		Expect(client.client.Addr.String()).To(Equal(kafkago.TCP("localhost:9092").String()))
	})

	It("should return a client with all broker addresses when multiple brokers are provided", func() {
		brokers := []string{"localhost:9092", "localhost:9093", "localhost:9094"}

		client, ok := NewBrokerClient(Config{Brokers: brokers}).(*brokerClient)

		Expect(ok).To(BeTrue())
		Expect(client.client.Addr).ToNot(BeNil())
		address := client.client.Addr.String()
		for _, broker := range brokers {
			Expect(address).To(ContainSubstring(broker))
		}
	})

	It("should return a client with the configured timeout when timeout is provided", func() {
		client, ok := NewBrokerClient(Config{Timeout: 3 * time.Second}).(*brokerClient)

		Expect(ok).To(BeTrue())
		Expect(client.client.Timeout).To(Equal(3 * time.Second))
	})

	It("should return a client with the configured transport when transport is provided", func() {
		transport := &kafkago.Transport{DialTimeout: 5 * time.Second}

		client, ok := NewBrokerClient(Config{Transport: transport}).(*brokerClient)

		Expect(ok).To(BeTrue())
		Expect(client.client.Transport).To(BeIdenticalTo(transport))
	})

	// Regression: assigning a nil *kafkago.Transport straight into the
	// RoundTripper interface field yields a non-nil interface wrapping a nil
	// pointer, which panics on the first request instead of falling back to
	// kafkago.DefaultTransport.
	It("should return a client with a truly nil transport when transport is not provided", func() {
		client, ok := NewBrokerClient(Config{}).(*brokerClient)

		Expect(ok).To(BeTrue())
		Expect(client.client.Transport).To(BeNil())
		Expect(client.client.Transport == nil).To(BeTrue())
	})

	It("should return a client with a truly nil transport when a nil transport pointer is provided", func() {
		var nilTransport *kafkago.Transport

		client, ok := NewBrokerClient(Config{Transport: nilTransport}).(*brokerClient)

		Expect(ok).To(BeTrue())
		Expect(client.client.Transport == nil).To(BeTrue())
	})

	It("should return a client holding the given config when config is provided", func() {
		cfg := Config{
			Brokers:   []string{"localhost:9092"},
			TopicName: "any-topic",
			Timeout:   time.Second,
		}

		client, ok := NewBrokerClient(cfg).(*brokerClient)

		Expect(ok).To(BeTrue())
		Expect(client.cfg).To(Equal(cfg))
	})
})
