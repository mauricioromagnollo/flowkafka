package schemaregistry

import (
	"github.com/riferrei/srclient"
)

// SchemaRegistry defines the interface for interacting with a schema registry service.
// It provides methods for retrieving schemas, validating connections, and creating new schemas.
// Implementations of this interface should handle communication with schema registry backends
// such as Confluent Schema Registry or compatible services.
type SchemaRegistry interface {
	// GetLatestSchema retrieves the latest schema for a given subject from the schema registry.
	GetLatestSchema(subject string) (*srclient.Schema, error)
	// ValidateConnection checks if the connection to the schema registry is valid.
	ValidateConnection() error
	// CreateNewSchema creates a new schema for a given subject in the schema registry.
	CreateNewSchema(subject, schema string, schemaType SchemaType) error
}
