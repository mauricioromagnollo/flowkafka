package producer

import (
	"github.com/mauricioromagnollo/flowkafka/internal/schemaregistry"
	"github.com/mauricioromagnollo/flowkafka/internal/shared/brokerclient"
	"github.com/mauricioromagnollo/flowkafka/internal/shared/transport"
)

// ProducerMessage wraps Kafka publishing with retry and clean shutdown support.
type producerClient struct {
	cfg            Config
	writer         Writer
	brokerClient   brokerclient.BrokerClient
	schemaRegistry schemaregistry.SchemaRegistry
}

// NewProducer creates a new Kafka producer with the given configuration.
func NewProducer(cfg Config) Producer {
	transport := transport.NewTransport(transport.Config{
		SaslUsername: cfg.SaslUsername,
		SaslPassword: cfg.SaslPassword,
	})
	writer := newKafkaGoWriter(cfg, transport)
	client := brokerclient.NewBrokerClient(brokerclient.Config{
		Brokers:   cfg.Brokers,
		TopicName: cfg.TopicName,
		Timeout:   cfg.ClientTimeout,
		Transport: transport,
	})

	return &producerClient{
		cfg:            cfg,
		writer:         writer,
		brokerClient:   client,
		schemaRegistry: cfg.SchemaRegistry,
	}
}
