package producer

import (
	"github.com/mauricioromagnollo/flowkafka/internal/schemaregistry"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HasSchemaRegistry", func() {
	It("should return false when schema registry is not configured", func() {
		config := Config{}
		producer := NewProducer(config)

		hasSchemaRegistry := producer.HasSchemaRegistry()

		Expect(hasSchemaRegistry).To(BeFalse())
	})

	It("should return true when schema registry is configured", func() {
		type MockSchemaRegistry struct {
			schemaregistry.SchemaRegistry
		}
		config := Config{
			SchemaRegistry: &MockSchemaRegistry{},
		}
		producer := NewProducer(config)

		hasSchemaRegistry := producer.HasSchemaRegistry()

		Expect(hasSchemaRegistry).To(BeTrue())
	})
})
