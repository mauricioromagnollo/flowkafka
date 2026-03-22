package consumer

import (
	"github.com/mauricioromagnollo/flowkafka/internal/shared/brokerclient"
	"github.com/mauricioromagnollo/flowkafka/internal/shared/transport"
	kafkago "github.com/segmentio/kafka-go"
)

// consumerClient wraps Kafka reading with consumer group and manual commit support.
type consumerClient struct {
	reader       *kafkago.Reader
	cfg          Config
	transport    *kafkago.Transport
	brokerClient brokerclient.BrokerClient
}

// NewConsumer creates a new Kafka consumer with the given configuration.
func NewConsumer(config Config) Consumer {
	dialer := buildDialer(config)
	transport := transport.NewTransport(transport.Config{
		SaslUsername: config.SaslUsername,
		SaslPassword: config.SaslPassword,
	})
	brokerClient := brokerclient.NewBrokerClient(brokerclient.Config{
		Brokers:   config.Brokers,
		TopicName: config.TopicName,
		Transport: transport,
	})
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: config.Brokers,
		GroupID: config.GroupID,
		Topic:   config.TopicName,
		Dialer:  dialer,
	})

	return &consumerClient{
		reader:       reader,
		cfg:          config,
		transport:    transport,
		brokerClient: brokerClient,
	}
}
