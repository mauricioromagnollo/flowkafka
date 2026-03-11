package producer

import "github.com/mauricioromagnollo/flowkafka/internal/schemaregistry"

// Config holds the configuration options for the Kafka producer.
type Config struct {
	Brokers        []string                      // Brokers is a list of Kafka broker addresses.
	SaslUsername   string                        // SaslUsername is the username for SASL authentication.
	SaslPassword   string                        // SaslPassword is the password for SASL authentication.
	TopicName      string                        // TopicName is the name of the Kafka topic to produce messages to.
	SchemaRegistry schemaregistry.SchemaRegistry // SchemaRegistry is the client for interacting with the Schema Registry.
}

// Message represents a Kafka message.
type Message struct {
	Key       []byte // Key is the key of the message.
	Value     []byte // Value is the value of the message.
	Headers   any    // Headers contains any additional headers associated with the message.
	Partition int    // Partition is the partition number of the message.
	Offset    int64  // Offset is the offset of the message within the partition.
}
