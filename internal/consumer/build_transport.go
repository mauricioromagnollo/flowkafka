package consumer

import (
	"crypto/tls"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func buildTransport(config Config) *kafkago.Transport {
	transport := &kafkago.Transport{
		DialTimeout: 10 * time.Second,
	}

	if config.SaslPassword != "" {
		transport.SASL = plain.Mechanism{
			Username: config.SaslUsername,
			Password: config.SaslPassword,
		}

		transport.TLS = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return transport
}
