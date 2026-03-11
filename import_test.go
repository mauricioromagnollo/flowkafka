package flowkafka_test

import (
	"testing"

	"github.com/mauricioromagnollo/flowkafka"
)

func TestPublicAPICompiles(_ *testing.T) {
	// Producer
	_ = flowkafka.ProducerConfig{}
	_ = flowkafka.Producer(nil)
	_ = flowkafka.NewProducer(flowkafka.ProducerConfig{})

	// Schema Registry
	_ = flowkafka.SchemaRegistryConfig{}
	_ = flowkafka.SchemaRegistry(nil)
	_ = flowkafka.NewSchemaRegistry(flowkafka.SchemaRegistryConfig{})
}
