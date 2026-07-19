package flowkafka_test

import (
	"testing"

	"github.com/mauricioromagnollo/flowkafka"
)

func TestPublicAPICompiles(_ *testing.T) {
	// Schema Registry
	_ = flowkafka.SchemaRegistryConfig{}
	_ = flowkafka.SchemaRegistry(nil)
	_ = flowkafka.NewSchemaRegistry(flowkafka.SchemaRegistryConfig{})

	// Producer
	_ = flowkafka.ProducerConfig{}
	_ = flowkafka.Producer(nil)
	_ = flowkafka.NewProducer(flowkafka.ProducerConfig{})

	// Consumer
	_ = flowkafka.ConsumerConfig{}
	_ = flowkafka.Consumer(nil)
	_ = flowkafka.NewConsumer(flowkafka.ConsumerConfig{
		GroupID:   "any_group_id",
		Brokers:   []string{"localhost:9092"},
		TopicName: "any_topic_name",
	})

	// General Types
	_ = flowkafka.Message{}
	_ = flowkafka.SchemaType("AVRO")
	_ = flowkafka.RequiredAcks(0)

	// Required Acks
	_ = flowkafka.RequiredAcksNone
	_ = flowkafka.RequiredAcksOne
	_ = flowkafka.RequiredAcksAll
}

func TestPublicRequiredAcksValues(t *testing.T) {
	testCases := []struct {
		name     string
		acks     flowkafka.RequiredAcks
		expected int
	}{
		{name: "RequiredAcksNone", acks: flowkafka.RequiredAcksNone, expected: 0},
		{name: "RequiredAcksOne", acks: flowkafka.RequiredAcksOne, expected: 1},
		{name: "RequiredAcksAll", acks: flowkafka.RequiredAcksAll, expected: -1},
	}

	for _, testCase := range testCases {
		if got := int(testCase.acks); got != testCase.expected {
			t.Errorf("%s = %d, want %d", testCase.name, got, testCase.expected)
		}
	}
}

// The public types are aliases of the internal ones, so a value produced by the
// constructors must satisfy the exported interfaces without conversion. These
// declarations fail to compile if that ever stops being true.
var (
	_ flowkafka.Producer       = flowkafka.NewProducer(flowkafka.ProducerConfig{})
	_ flowkafka.SchemaRegistry = flowkafka.NewSchemaRegistry(flowkafka.SchemaRegistryConfig{})
	_ flowkafka.Consumer       = flowkafka.NewConsumer(flowkafka.ConsumerConfig{
		GroupID:   "any_group_id",
		Brokers:   []string{"localhost:9092"},
		TopicName: "any_topic_name",
	})
)

func TestPublicConstructorsReturnNonNil(t *testing.T) {
	if producer := flowkafka.NewProducer(flowkafka.ProducerConfig{}); producer == nil {
		t.Fatal("NewProducer returned nil")
	}

	consumer := flowkafka.NewConsumer(flowkafka.ConsumerConfig{
		GroupID:   "any_group_id",
		Brokers:   []string{"localhost:9092"},
		TopicName: "any_topic_name",
	})
	if consumer == nil {
		t.Fatal("NewConsumer returned nil")
	}

	if schemaRegistry := flowkafka.NewSchemaRegistry(flowkafka.SchemaRegistryConfig{}); schemaRegistry == nil {
		t.Fatal("NewSchemaRegistry returned nil")
	}
}
