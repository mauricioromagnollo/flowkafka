package kafka

import (
	"sync"

	kafkago "github.com/segmentio/kafka-go"
)

// Client is the internal Kafka client.
type Client struct {
	config  Config
	dialer  *kafkago.Dialer
	writer  *kafkago.Writer
	mu      sync.Mutex
	readers map[string]*kafkago.Reader // per topic (lazy)
}

// NewKafka creates a new Kafka client.
func NewKafka(config Config) *Client {
	dialer := newDialer(config)
	writer := newWriter(config, dialer)

	return &Client{
		config:  config,
		dialer:  dialer,
		writer:  writer,
		readers: make(map[string]*kafkago.Reader),
	}
}
