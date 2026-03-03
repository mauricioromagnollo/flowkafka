package schemaregistry

import "github.com/riferrei/srclient"

// Config holds the configuration parameters for connecting to a Kafka Schema Registry.
// It contains the endpoint URL and SASL credentials required for authentication.
type Config struct {
	Endpoint     string
	SaslUsername string
	SaslPassword string
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
