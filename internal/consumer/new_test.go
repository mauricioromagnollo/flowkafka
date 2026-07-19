package consumer

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewConsumer", func() {
	It("should return a consumer with all fields populated", func() {
		config := Config{
			Brokers:   []string{"localhost:9092"},
			GroupID:   "any-group",
			TopicName: "any-topic",
		}

		consumer := NewConsumer(config)

		Expect(consumer).ToNot(BeNil())

		client, ok := consumer.(*consumerClient)
		Expect(ok).To(BeTrue())
		Expect(client.reader).ToNot(BeNil())
		Expect(client.cfg).To(Equal(config))
		Expect(client.transport).ToNot(BeNil())
		Expect(client.brokerClient).ToNot(BeNil())
	})

	// Pre-existing behavior: kafka-go validates the reader config eagerly and panics
	// when the broker list is empty. brokerclient.NewBrokerClient guards against this
	// case, but the reader does not. This spec documents the current behavior.
	It("should panic when no broker is provided", func() {
		config := Config{
			Brokers:   []string{},
			TopicName: "any-topic",
		}

		Expect(func() { NewConsumer(config) }).
			To(PanicWith(MatchError(ContainSubstring("cannot create a new kafka reader with an empty list of broker addresses"))))
	})

	It("should return a consumer with a SASL transport when sasl credentials are provided", func() {
		config := Config{
			Brokers:      []string{"localhost:9092"},
			GroupID:      "any-group",
			TopicName:    "any-topic",
			SaslUsername: "any-username",
			SaslPassword: "any-password",
		}

		consumer := NewConsumer(config)

		client, ok := consumer.(*consumerClient)
		Expect(ok).To(BeTrue())
		Expect(client.transport).ToNot(BeNil())
		Expect(client.transport.SASL).ToNot(BeNil())
		Expect(client.transport.TLS).ToNot(BeNil())
	})
})
