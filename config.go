package flowkafka

// KafkaConfig holds the configuration for connecting to a Kafka cluster.
type KafkaConfig struct {
	Addresses    []string // Addresses is a list of Kafka broker addresses.
	GroupID      string   // GroupID is the consumer group ID.
	SaslUsername string   // SaslUsername is the username for SASL authentication.
	SaslPassword string   // SaslPassword is the password for SASL authentication.
}
