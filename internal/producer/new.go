package producer

import (
	"time"

	"github.com/mauricioromagnollo/flowkafka/internal/schemaregistry"
	kafkago "github.com/segmentio/kafka-go"
)

// ProducerMessage wraps Kafka publishing with retry and clean shutdown support.
type producerClient struct {
	writer         *kafkago.Writer
	cfg            Config
	schemaRegistry schemaregistry.SchemaRegistry
}

// NewProducer creates a new Kafka producer with the given configuration.
func NewProducer(config Config) Producer {
	transport := buildTransport(config)

	writer := &kafkago.Writer{
		Addr:         kafkago.TCP(config.Brokers...),
		Topic:        config.TopicName,
		Balancer:     &kafkago.LeastBytes{},
		RequiredAcks: kafkago.RequireAll,
		MaxAttempts:  10,
		BatchTimeout: 1 * time.Second,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		Transport:    transport,
	}

	return &producerClient{
		writer:         writer,
		cfg:            config,
		schemaRegistry: config.SchemaRegistry,
	}
}
