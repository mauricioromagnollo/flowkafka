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
}
