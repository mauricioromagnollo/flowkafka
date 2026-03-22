package types

// Message represents a Kafka message received by the consumer.
type Message struct {
	Key       []byte // Key is the key of the message.
	Value     []byte // Value is the value of the message.
	Headers   any    // Headers contains any additional headers associated with the message.
	Partition int    // Partition is the partition number of the message.
	Offset    int64  // Offset is the offset of the message within the partition.
}
