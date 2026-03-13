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
}
