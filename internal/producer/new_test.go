package producer

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewProducer", func() {
	It("should return a producer when a full config is provided", func() {
		schemaRegistry := new(schemaRegistryStub)
		config := Config{
			Brokers:        []string{"localhost:9092"},
			TopicName:      "any-topic-name",
			SaslUsername:   "any-username",
			SaslPassword:   "any-password",
			SchemaRegistry: schemaRegistry,
			MaxAttempts:    5,
			BatchTimeout:   1 * time.Second,
			WriteTimeout:   2 * time.Second,
			ReadTimeout:    3 * time.Second,
			RequiredAcks:   RequiredAcksAll,
			ClientTimeout:  4 * time.Second,
		}

		producer := NewProducer(config)

		Expect(producer).ToNot(BeNil())

		client, ok := producer.(*producerClient)
		Expect(ok).To(BeTrue())
		Expect(client.cfg).To(Equal(config))
		Expect(client.writer).ToNot(BeNil())
		Expect(client.brokerClient).ToNot(BeNil())
		Expect(client.schemaRegistry).To(Equal(schemaRegistry))
	})

	It("should return nil schema registry when it is not configured", func() {
		config := Config{
			Brokers:   []string{"localhost:9092"},
			TopicName: "any-topic-name",
		}

		producer := NewProducer(config)

		client, ok := producer.(*producerClient)
		Expect(ok).To(BeTrue())
		Expect(client.schemaRegistry).To(BeNil())
		Expect(producer.HasSchemaRegistry()).To(BeFalse())
	})

	It("should return a producer when the brokers list is empty", func() {
		config := Config{TopicName: "any-topic-name"}

		Expect(func() { NewProducer(config) }).ToNot(Panic())

		producer := NewProducer(config)

		Expect(producer).ToNot(BeNil())

		client, ok := producer.(*producerClient)
		Expect(ok).To(BeTrue())
		Expect(client.writer).ToNot(BeNil())
		Expect(client.brokerClient).ToNot(BeNil())
	})
})
