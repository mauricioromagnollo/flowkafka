package producer

import (
	"context"
	"time"

	"github.com/mauricioromagnollo/flowkafka/internal/shared/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kafkago "github.com/segmentio/kafka-go"
)

var _ = Describe("newKafkaGoWriter", func() {
	It("should return a writer built from the provided config", func() {
		transport := &kafkago.Transport{}
		config := Config{
			Brokers:      []string{"localhost:9092", "localhost:9093"},
			TopicName:    "any-topic-name",
			MaxAttempts:  7,
			BatchTimeout: 1 * time.Second,
			WriteTimeout: 2 * time.Second,
			ReadTimeout:  3 * time.Second,
			RequiredAcks: RequiredAcksAll,
		}

		writer := newKafkaGoWriter(config, transport)

		Expect(writer).ToNot(BeNil())

		kafkaWriter, ok := writer.(*kafkaGoWriter)
		Expect(ok).To(BeTrue())
		Expect(kafkaWriter.writer.Addr).To(Equal(kafkago.TCP("localhost:9092", "localhost:9093")))
		Expect(kafkaWriter.writer.Topic).To(Equal("any-topic-name"))
		Expect(kafkaWriter.writer.MaxAttempts).To(Equal(7))
		Expect(kafkaWriter.writer.BatchTimeout).To(Equal(1 * time.Second))
		Expect(kafkaWriter.writer.WriteTimeout).To(Equal(2 * time.Second))
		Expect(kafkaWriter.writer.ReadTimeout).To(Equal(3 * time.Second))
		Expect(kafkaWriter.writer.RequiredAcks).To(Equal(kafkago.RequireAll))
		Expect(kafkaWriter.writer.Transport).To(BeIdenticalTo(transport))
		Expect(kafkaWriter.writer.Balancer).To(BeAssignableToTypeOf(&kafkago.LeastBytes{}))
	})

	It("should return a writer with the required acks translated from the config", func() {
		config := Config{TopicName: "any-topic-name", RequiredAcks: RequiredAcksOne}

		writer := newKafkaGoWriter(config, &kafkago.Transport{})

		kafkaWriter, ok := writer.(*kafkaGoWriter)
		Expect(ok).To(BeTrue())
		Expect(kafkaWriter.writer.RequiredAcks).To(Equal(kafkago.RequireOne))
	})

	// KNOWN DIVERGENCE: Config.RequiredAcks is documented with a `default:"-1"`
	// tag (RequiredAcksAll), but nothing in the package reads struct tags.
	// newKafkaGoWriter builds a &kafkago.Writer{} struct literal, so it also
	// bypasses the coercion kafkago.NewWriter applies (RequiredAcks == 0 is
	// rewritten to RequireAll there). The result is that a zero-value Config
	// produces RequiredAcks == 0, which kafka-go interprets as RequireNone
	// (fire-and-forget), not RequireAll. This spec locks the ACTUAL current
	// behavior, not the behavior implied by the tag.
	It("should return required acks none when the config is a zero value", func() {
		config := Config{}

		writer := newKafkaGoWriter(config, &kafkago.Transport{})

		kafkaWriter, ok := writer.(*kafkaGoWriter)
		Expect(ok).To(BeTrue())
		Expect(kafkaWriter.writer.RequiredAcks).To(Equal(kafkago.RequireNone))
		Expect(kafkaWriter.writer.RequiredAcks).ToNot(Equal(kafkago.RequireAll))
	})

	It("should return nil address when the brokers list is empty", func() {
		config := Config{TopicName: "any-topic-name"}

		writer := newKafkaGoWriter(config, &kafkago.Transport{})

		kafkaWriter, ok := writer.(*kafkaGoWriter)
		Expect(ok).To(BeTrue())
		Expect(kafkaWriter.writer.Addr).To(BeNil())
	})
})

var _ = Describe("kafkaGoWriter", func() {
	// newUnreachableWriter builds a writer pointed at a closed port so the specs
	// exercise the kafka-go delegation without reaching a real broker.
	newUnreachableWriter := func() *kafkaGoWriter {
		writer := newKafkaGoWriter(Config{
			Brokers:      []string{"127.0.0.1:1"},
			TopicName:    "any-topic-name",
			MaxAttempts:  1,
			WriteTimeout: 10 * time.Millisecond,
			ReadTimeout:  10 * time.Millisecond,
		}, &kafkago.Transport{DialTimeout: 10 * time.Millisecond})

		kafkaWriter, ok := writer.(*kafkaGoWriter)
		Expect(ok).To(BeTrue())

		return kafkaWriter
	}

	Describe("WriteMessages", func() {
		It("should return error when the context is already canceled", func() {
			writer := newUnreachableWriter()
			defer func() { Expect(writer.Close()).To(Succeed()) }()

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			err := writer.WriteMessages(ctx, types.Message{
				Key:   []byte("any-message-key"),
				Value: []byte("any-message-value"),
			})

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Close", func() {
		It("should return nil when the underlying writer is closed", func() {
			writer := newUnreachableWriter()

			err := writer.Close()

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
