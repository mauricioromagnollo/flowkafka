package flowkafka

import (
	"github.com/mauricioromagnollo/flowkafka/internal/schemaregistry"
	"github.com/riferrei/srclient"
)

// SchemaRegistryClient defines the interface for interacting with a schema registry service.
// It provides methods for retrieving schemas, validating connections, and creating new schemas.
// Implementations of this interface should handle communication with schema registry backends
// such as Confluent Schema Registry or compatible services.
type SchemaRegistryClient interface {
	// GetLatestSchema retrieves the latest schema for a given subject from the schema registry.
	GetLatestSchema(subject string) (*srclient.Schema, error)
	// ValidateConnection checks if the connection to the schema registry is valid.
	ValidateConnection() error
	// CreateNewSchema creates a new schema for a given subject in the schema registry.
	CreateNewSchema(subject, schema string, schemaType SchemaType) error
}

// SchemaType is a type alias for srclient.SchemaType that represents the format
// of a schema in the schema registry (e.g., AVRO, JSON, PROTOBUF).
type SchemaType srclient.SchemaType

// SchemaType constants represent the different types of schemas.
const (
	SchemaTypeAvro     SchemaType = "AVRO"
	SchemaTypeJSON     SchemaType = "JSON"
	SchemaTypeProtobuf SchemaType = "PROTOBUF"
)

// SchemaRegistryConfig holds the configuration parameters for connecting to a Kafka Schema Registry.
// It contains the endpoint URL and SASL credentials required for authentication.
type SchemaRegistryConfig struct {
	Endpoint     string
	SaslUsername string
	SaslPassword string
}

type schemaRegistryWrapper struct {
	client *schemaregistry.Client
}

// NewSchemaRegistryClient creates a new SchemaRegistryClient instance using the provided configuration.
// It initializes the underlying schema registry client and returns a wrapper that implements the SchemaRegistryClient interface.
func NewSchemaRegistryClient(config SchemaRegistryConfig) SchemaRegistryClient {
	client := schemaregistry.NewSchemaRegistry(schemaregistry.Config{
		Endpoint:     config.Endpoint,
		SaslUsername: config.SaslUsername,
		SaslPassword: config.SaslPassword,
	})

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
	return w.client.CreateNewSchema(subject, schema, schemaregistry.SchemaType(schemaType))
}
