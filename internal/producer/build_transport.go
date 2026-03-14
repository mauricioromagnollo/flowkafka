package producer

import (
	"crypto/tls"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// buildTransport creates a Kafka transport with optional SASL/PLAIN authentication and TLS configuration based on the provided producer configuration.
func buildTransport(cfg Config) *kafkago.Transport {
	transport := &kafkago.Transport{
		DialTimeout: 10 * time.Second,
	}

	if cfg.SaslPassword != "" {
		transport.SASL = plain.Mechanism{
			Username: cfg.SaslUsername,
			Password: cfg.SaslPassword,
		}

		transport.TLS = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return transport
}
