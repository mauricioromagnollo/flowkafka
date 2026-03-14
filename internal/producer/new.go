package producer

import (
	"github.com/mauricioromagnollo/flowkafka/internal/schemaregistry"
)

// ProducerMessage wraps Kafka publishing with retry and clean shutdown support.
type producerClient struct {
	writer         Writer
	cfg            Config
	schemaRegistry schemaregistry.SchemaRegistry
}

// NewProducer creates a new Kafka producer with the given configuration.
func NewProducer(cfg Config) Producer {
	transport := buildTransport(cfg)
	writer := newKafkaGoWriter(cfg, transport)

	return &producerClient{
		writer:         writer,
		cfg:            cfg,
		schemaRegistry: cfg.SchemaRegistry,
	}
}
