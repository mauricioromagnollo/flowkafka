package flowkafka

import (
	"context"

	"github.com/mauricioromagnollo/flowkafka/internal/kafka"
)

// KafkaClient provides an interface for interacting with Apache Kafka message broker.
type KafkaClient interface {
	// ValidateConnection checks if the Kafka cluster is reachable and the connection can be established.
	ValidateConnection(ctx context.Context) error
	// ReadMessage fetches one message from the given topic (manual commit).
	ReadMessage(ctx context.Context, topic string) (Message, error)
	// CommitMessage commits the message offset for the given topic.
	CommitMessage(ctx context.Context, topic string, msg Message) error
	// ProduceAvroMessage produces a message serialized as Avro using Schema Registry (Confluent wire format).
	ProduceAvroMessage(ctx context.Context, topic string, key []byte, msg any) error
	// ProduceJSONMessage produces a message serialized as JSON.
	ProduceJSONMessage(ctx context.Context, topic string, key []byte, msg any) error
	// Close closes the Kafka connection.
	Close() error
	// CloseReader closes the reader for the given topic, unblocking any pending ReadMessage calls.
	CloseReader(topic string) error
}

// Message represents a Kafka message.
type Message struct {
	Key       []byte // Key is the key of the message.
	Value     []byte // Value is the value of the message.
	Headers   any    // Headers contains any additional headers associated with the message.
	Partition int    // Partition is the partition number of the message.
	Offset    int64  // Offset is the offset of the message within the partition.
}

// KafkaConfig holds the configuration for connecting to a Kafka cluster.
type KafkaConfig struct {
	Addresses    []string // Addresses is a list of Kafka broker addresses.
	GroupID      string   // GroupID is the consumer group ID.
	SaslUsername string   // SaslUsername is the username for SASL authentication.
	SaslPassword string   // SaslPassword is the password for SASL authentication.
}

// kafkaWrapper adapts the internal kafka client to the public Kafka interface,
// converting between the public and internal types at the boundary.
type kafkaWrapper struct {
	client *kafka.Client
}

// NewKafkaClient creates a new Kafka client from the given public configuration.
func NewKafkaClient(cfg KafkaConfig) KafkaClient {
	return &kafkaWrapper{
		client: kafka.NewKafka(kafka.Config{
			Addresses:    cfg.Addresses,
			GroupID:      cfg.GroupID,
			SaslUsername: cfg.SaslUsername,
			SaslPassword: cfg.SaslPassword,
		}),
	}
}

func (w *kafkaWrapper) ValidateConnection(ctx context.Context) error {
	return w.client.ValidateConnection(ctx)
}

func (w *kafkaWrapper) ReadMessage(ctx context.Context, topic string) (Message, error) {
	m, err := w.client.ReadMessage(ctx, topic)
	if err != nil {
		return Message{}, err
	}

	return Message{
		Key:       m.Key,
		Value:     m.Value,
		Headers:   m.Headers,
		Partition: m.Partition,
		Offset:    m.Offset,
	}, nil
}

func (w *kafkaWrapper) CommitMessage(ctx context.Context, topic string, msg Message) error {
	return w.client.CommitMessage(ctx, topic, kafka.Message{
		Key:       msg.Key,
		Value:     msg.Value,
		Partition: msg.Partition,
		Offset:    msg.Offset,
	})
}

func (w *kafkaWrapper) ProduceAvroMessage(ctx context.Context, topic string, key []byte, msg any) error {
	return w.client.ProduceAvroMessage(ctx, topic, key, msg)
}

func (w *kafkaWrapper) ProduceJSONMessage(ctx context.Context, topic string, key []byte, msg any) error {
	return w.client.ProduceJSONMessage(ctx, topic, key, msg)
}

func (w *kafkaWrapper) Close() error {
	return w.client.Close()
}

func (w *kafkaWrapper) CloseReader(topic string) error {
	return w.client.CloseReader(topic)
}
