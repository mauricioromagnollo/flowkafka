package flowkafka

import (
	"github.com/mauricioromagnollo/flowkafka/internal/schemaregistry"
	"github.com/riferrei/srclient"
)

// SchemaRegistry defines the interface for interacting with a schema registry service.
// It provides methods for retrieving schemas, validating connections, and creating new schemas.
// Implementations of this interface should handle communication with schema registry backends
// such as Confluent Schema Registry or compatible services.
type SchemaRegistry = schemaregistry.SchemaRegistry

// SchemaRegistryConfig holds the configuration parameters for connecting to a Kafka Schema Registry.
// It contains the endpoint URL and SASL credentials required for authentication.
type SchemaRegistryConfig = schemaregistry.Config

// SchemaType is a type alias for srclient.SchemaType that represents the format
// of a schema in the schema registry (e.g., AVRO, JSON, PROTOBUF).
type SchemaType = schemaregistry.SchemaType

type schemaRegistryWrapper struct {
	client *schemaregistry.Client
}

// NewSchemaRegistry creates a new SchemaRegistry instance using the provided configuration.
// It initializes the underlying schema registry client and returns a wrapper that implements the SchemaRegistry interface.
func NewSchemaRegistry(config SchemaRegistryConfig) SchemaRegistry {
	client := schemaregistry.NewSchemaRegistry(config)

	return &schemaRegistryWrapper{
		client: client,
	}
}

func (w *schemaRegistryWrapper) GetLatestSchema(subject string) (*srclient.Schema, error) {
	return w.client.GetLatestSchema(subject)
}

func (w *schemaRegistryWrapper) ValidateConnection() error {
	return w.client.ValidateConnection()
}

func (w *schemaRegistryWrapper) CreateNewSchema(subject, schema string, schemaType SchemaType) error {
	return w.client.CreateNewSchema(subject, schema, schemaType)
}
