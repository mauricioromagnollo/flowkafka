package brokerclient

import (
	"net"

	kafkago "github.com/segmentio/kafka-go"
)

type brokerClient struct {
	cfg    Config
	client *kafkago.Client
}

// NewBrokerClient creates a new Kafka broker client with the given configuration.
// When no broker address is configured the client is built without an address,
// so that construction never panics. Any request made through such a client
// fails at call time instead.
func NewBrokerClient(cfg Config) BrokerClient {
	var addr net.Addr
	if len(cfg.Brokers) > 0 {
		addr = kafkago.TCP(cfg.Brokers...)
	}

	client := &kafkago.Client{
		Addr:    addr,
		Timeout: cfg.Timeout,
	}

	// Assigning a nil *kafkago.Transport to the RoundTripper interface field
	// would produce a non-nil interface holding a nil pointer, which panics on
	// the first request instead of falling back to kafkago.DefaultTransport.
	if cfg.Transport != nil {
		client.Transport = cfg.Transport
	}

	return &brokerClient{
		cfg:    cfg,
		client: client,
	}
}
