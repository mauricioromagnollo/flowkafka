package consumer

import (
	"crypto/tls"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func buildDialer(config Config) *kafkago.Dialer {
	dialer := &kafkago.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	if config.SaslPassword != "" {
		dialer.SASLMechanism = plain.Mechanism{
			Username: config.SaslUsername,
			Password: config.SaslPassword,
		}

		dialer.TLS = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return dialer
}
