package transport

// Config holds the configuration for creating a Kafka transport, including optional SASL/PLAIN authentication credentials.
type Config struct {
	SaslUsername string
	SaslPassword string
}
