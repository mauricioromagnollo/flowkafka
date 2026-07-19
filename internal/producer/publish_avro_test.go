package producer

import (
	"context"
	"encoding/binary"
	"errors"

	avro "github.com/hamba/avro/v2"
	"github.com/mauricioromagnollo/flowkafka/internal/shared/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/riferrei/srclient"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("PublishAvro", func() {
	const (
		anyTopicName  = "any-topic-name"
		anySubject    = anyTopicName + "-value"
		anyAvroSchema = `{"type":"record","name":"Test","fields":[{"name":"id","type":"string"}]}`
	)

	var (
		mockWriter         *writerStub
		mockSchemaRegistry *schemaRegistryStub
	)

	// newSchema builds a srclient.Schema with a controlled ID and schema payload.
	// srclient.NewSchema is the only exported way to build a Schema outside the
	// package, since every field of the struct is unexported.
	newSchema := func(id int, schema string) *srclient.Schema {
		result, err := srclient.NewSchema(id, schema, srclient.Avro, 1, nil, nil, nil)
		Expect(err).ToNot(HaveOccurred())

		return result
	}

	BeforeEach(func() {
		mockWriter = new(writerStub)
		mockSchemaRegistry = new(schemaRegistryStub)
	})

	AfterEach(func() {
		mockWriter.AssertExpectations(GinkgoT())
		mockSchemaRegistry.AssertExpectations(GinkgoT())
	})

	It("should return error when schema registry is not configured", func() {
		producer := &producerClient{
			writer:         mockWriter,
			cfg:            Config{TopicName: anyTopicName},
			schemaRegistry: nil,
		}
		ctx := context.Background()

		err := producer.PublishAvro(ctx, nil, map[string]any{"id": "any-id"})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("schema registry is not configured"))
	})

	It("should return error when fail to get the latest schema", func() {
		producer := &producerClient{
			writer:         mockWriter,
			cfg:            Config{TopicName: anyTopicName},
			schemaRegistry: mockSchemaRegistry,
		}
		ctx := context.Background()

		mockSchemaRegistry.On("GetLatestSchema", anySubject).
			Return(nil, errors.New("any error reason returned by the schema registry"))

		err := producer.PublishAvro(ctx, nil, map[string]any{"id": "any-id"})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to get latest schema for topic any-topic-name"))
		Expect(err.Error()).To(ContainSubstring("any error reason returned by the schema registry"))
	})

	It("should return error when the schema cannot be parsed", func() {
		producer := &producerClient{
			writer:         mockWriter,
			cfg:            Config{TopicName: anyTopicName},
			schemaRegistry: mockSchemaRegistry,
		}
		ctx := context.Background()

		mockSchemaRegistry.On("GetLatestSchema", anySubject).
			Return(newSchema(1, "any-invalid-avro-schema"), nil)

		err := producer.PublishAvro(ctx, nil, map[string]any{"id": "any-id"})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to parse schema"))
	})

	It("should return error when the message does not match the schema", func() {
		producer := &producerClient{
			writer:         mockWriter,
			cfg:            Config{TopicName: anyTopicName},
			schemaRegistry: mockSchemaRegistry,
		}
		ctx := context.Background()

		mockSchemaRegistry.On("GetLatestSchema", anySubject).
			Return(newSchema(1, anyAvroSchema), nil)

		err := producer.PublishAvro(ctx, nil, map[string]any{"any-unknown-field": "any-value"})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to marshal message to Avro"))
	})

	It("should return error when the schema ID is out of the uint32 range", func() {
		producer := &producerClient{
			writer:         mockWriter,
			cfg:            Config{TopicName: anyTopicName},
			schemaRegistry: mockSchemaRegistry,
		}
		ctx := context.Background()

		mockSchemaRegistry.On("GetLatestSchema", anySubject).
			Return(newSchema(-1, anyAvroSchema), nil)

		err := producer.PublishAvro(ctx, nil, map[string]any{"id": "any-id"})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("schema ID -1 is out of range for uint32"))
	})

	It("should return error when fail to write message", func() {
		producer := &producerClient{
			writer:         mockWriter,
			cfg:            Config{TopicName: anyTopicName},
			schemaRegistry: mockSchemaRegistry,
		}
		ctx := context.Background()

		mockSchemaRegistry.On("GetLatestSchema", anySubject).
			Return(newSchema(1, anyAvroSchema), nil)
		mockWriter.On("WriteMessages", ctx, mock.Anything).
			Return(errors.New("any error reason returned by the writer"))

		err := producer.PublishAvro(ctx, nil, map[string]any{"id": "any-id"})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to publish message to Kafka topic any-topic-name"))
		Expect(err.Error()).To(ContainSubstring("any error reason returned by the writer"))
	})

	It("should return nil when the message is published in the confluent wire format", func() {
		const anySchemaID = 4242

		producer := &producerClient{
			writer:         mockWriter,
			cfg:            Config{TopicName: anyTopicName},
			schemaRegistry: mockSchemaRegistry,
		}
		ctx := context.Background()
		msgKey := []byte("any-message-key")
		msg := map[string]any{"id": "any-id"}

		codec, err := avro.Parse(anyAvroSchema)
		Expect(err).ToNot(HaveOccurred())
		expectedPayload, err := avro.Marshal(codec, msg)
		Expect(err).ToNot(HaveOccurred())

		var writtenMessage types.Message
		mockSchemaRegistry.On("GetLatestSchema", anySubject).
			Return(newSchema(anySchemaID, anyAvroSchema), nil)
		mockWriter.On("WriteMessages", ctx, mock.Anything).
			Run(func(args mock.Arguments) {
				writtenMessage = args.Get(1).(types.Message)
			}).
			Return(nil)

		err = producer.PublishAvro(ctx, msgKey, msg)

		Expect(err).ToNot(HaveOccurred())
		Expect(writtenMessage.Key).To(Equal(msgKey))
		Expect(writtenMessage.Value).To(HaveLen(5 + len(expectedPayload)))
		Expect(writtenMessage.Value[0]).To(Equal(byte(0)))
		Expect(binary.BigEndian.Uint32(writtenMessage.Value[1:5])).To(Equal(uint32(anySchemaID)))
		Expect(writtenMessage.Value[5:]).To(Equal(expectedPayload))
	})
})
