package producer

import (
	"github.com/mauricioromagnollo/flowkafka/internal/schemaregistry"
	"github.com/riferrei/srclient"
	"github.com/stretchr/testify/mock"
)

// schemaRegistryStub is a mock implementation of schemaregistry.SchemaRegistry.
// Only GetLatestSchema is stubbed because it is the single method the producer
// calls; the embedded interface covers the rest of the contract and panics if
// an unexpected method is ever reached.
type schemaRegistryStub struct {
	schemaregistry.SchemaRegistry
	mock.Mock
}

func (m *schemaRegistryStub) GetLatestSchema(subject string) (*srclient.Schema, error) {
	args := m.Called(subject)

	schema, _ := args.Get(0).(*srclient.Schema)

	return schema, args.Error(1)
}
