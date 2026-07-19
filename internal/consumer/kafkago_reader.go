package consumer

import (
	"context"

	kafkago "github.com/segmentio/kafka-go"
)

type kafkaGoReader struct {
	reader *kafkago.Reader
}

func newKafkaGoReader(cfg Config, dialer *kafkago.Dialer) Reader {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Topic:   cfg.TopicName,
		Dialer:  dialer,
	})

	return &kafkaGoReader{reader: reader}
}

func (r *kafkaGoReader) FetchMessage(ctx context.Context) (kafkago.Message, error) {
	return r.reader.FetchMessage(ctx)
}

func (r *kafkaGoReader) CommitMessages(ctx context.Context, msgs ...kafkago.Message) error {
	return r.reader.CommitMessages(ctx, msgs...)
}

func (r *kafkaGoReader) Close() error {
	return r.reader.Close()
}
