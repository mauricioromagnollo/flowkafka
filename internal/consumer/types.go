package consumer

// Config holds the configuration options for the Kafka consumer.
type Config struct {
	Brokers      []string // Brokers is a list of Kafka broker addresses.
	GroupID      string   // GroupID is the consumer group ID.
	TopicName    string   // TopicName is the name of the Kafka topic to consume messages from.
	SaslUsername string   // SaslUsername is the username for SASL authentication.
	SaslPassword string   // SaslPassword is the password for SASL authentication.
}

// Message represents a Kafka message received by the consumer.
type Message struct {
	Key       []byte // Key is the key of the message.
	Value     []byte // Value is the value of the message.
	Headers   any    // Headers contains any additional headers associated with the message.
	Partition int    // Partition is the partition number of the message.
	Offset    int64  // Offset is the offset of the message within the partition.
}
