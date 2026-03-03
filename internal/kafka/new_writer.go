package kafka

import (
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

func newWriter(config Config, dialer *kafkago.Dialer) *kafkago.Writer {
	writer := &kafkago.Writer{
		Addr:         kafkago.TCP(config.Addresses...),
		Balancer:     &kafkago.LeastBytes{},
		RequiredAcks: kafkago.RequireAll,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Transport: &kafkago.Transport{
			Dial: dialer.DialFunc,
			SASL: dialer.SASLMechanism,
			TLS:  dialer.TLS,
		},
	}

	return writer
}
