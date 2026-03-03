package flowkafka_test

import (
	"testing"

	"github.com/mauricioromagnollo/flowkafka"
)

func TestPublicAPICompiles(_ *testing.T) {
	// kafka
	_ = flowkafka.KafkaConfig{}
	_ = flowkafka.Message{}
	_ = flowkafka.KafkaClient(nil)
	_ = flowkafka.NewKafkaClient(flowkafka.KafkaConfig{})

	// schemaregistry
	_ = flowkafka.SchemaRegistryConfig{}
	_ = flowkafka.SchemaRegistryClient(nil)
	_ = flowkafka.NewSchemaRegistryClient(flowkafka.SchemaRegistryConfig{})
}
