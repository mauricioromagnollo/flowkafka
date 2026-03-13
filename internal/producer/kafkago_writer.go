package producer

import (
	"context"

	kafkago "github.com/segmentio/kafka-go"
)

type kafkaGoWriter struct {
	writer *kafkago.Writer
}

func newKafkaGoWriter(cfg Config, transport *kafkago.Transport) Writer {
	writer := &kafkago.Writer{
		Addr:         kafkago.TCP(cfg.Brokers...),
		Topic:        cfg.TopicName,
		Balancer:     &kafkago.LeastBytes{},
		RequiredAcks: kafkago.RequiredAcks(cfg.RequiredAcks),
		MaxAttempts:  cfg.MaxAttempts,
		BatchTimeout: cfg.BatchTimeout,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		Transport:    transport,
	}

	return &kafkaGoWriter{writer: writer}
}

func (w *kafkaGoWriter) WriteMessages(ctx context.Context, msg Message) error {
	return w.writer.WriteMessages(ctx, kafkago.Message{
		Value: msg.Value,
		Key:   msg.Key,
	})
}

func (w *kafkaGoWriter) Close() error {
	return w.writer.Close()
}
