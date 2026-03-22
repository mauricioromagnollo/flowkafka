package consumer

// Config holds the configuration options for the Kafka consumer.
type Config struct {
	Brokers      []string // Brokers is a list of Kafka broker addresses.
	GroupID      string   // GroupID is the consumer group ID.
	TopicName    string   // TopicName is the name of the Kafka topic to consume messages from.
	SaslUsername string   // SaslUsername is the username for SASL authentication.
	SaslPassword string   // SaslPassword is the password for SASL authentication.
}
