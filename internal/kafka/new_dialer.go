package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func newDialer(config Config) *kafkago.Dialer {
	dialer := &kafkago.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	if config.SaslPassword != "" {
		dialer.SASLMechanism = plain.Mechanism{
			Username: config.SaslUsername,
			Password: config.SaslPassword,
		}

		rootCAs, err := x509.SystemCertPool()
		if err != nil {
			rootCAs = nil
		}

		dialer.TLS = &tls.Config{
			RootCAs:    rootCAs,
			MinVersion: tls.VersionTLS12,
		}
	}

	return dialer
}
