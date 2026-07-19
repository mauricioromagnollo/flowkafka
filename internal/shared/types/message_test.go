package types

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kafkago "github.com/segmentio/kafka-go"
)

var _ = Describe("Message", func() {
	It("should return the configured values when all fields are provided", func() {
		message := Message{
			Key:       []byte("any-key"),
			Value:     []byte("any-value"),
			Headers:   map[string]string{"any-header": "any-header-value"},
			Partition: 3,
			Offset:    42,
		}

		Expect(message.Key).To(Equal([]byte("any-key")))
		Expect(message.Value).To(Equal([]byte("any-value")))
		Expect(message.Headers).To(Equal(map[string]string{"any-header": "any-header-value"}))
		Expect(message.Partition).To(Equal(3))
		Expect(message.Offset).To(Equal(int64(42)))
	})

	It("should return empty values when the message is a zero value", func() {
		var message Message

		Expect(message.Key).To(BeNil())
		Expect(message.Value).To(BeNil())
		Expect(message.Headers).To(BeNil())
		Expect(message.Partition).To(Equal(0))
		Expect(message.Offset).To(Equal(int64(0)))
	})

	It("should accept kafka-go headers when headers are set by the consumer or producer", func() {
		headers := []kafkago.Header{
			{Key: "any-header", Value: []byte("any-header-value")},
		}
		message := Message{Headers: headers}

		typedHeaders, ok := message.Headers.([]kafkago.Header)

		Expect(ok).To(BeTrue())
		Expect(typedHeaders).To(HaveLen(1))
		Expect(typedHeaders[0].Key).To(Equal("any-header"))
		Expect(typedHeaders[0].Value).To(Equal([]byte("any-header-value")))
	})

	It("should keep the value nil when only the key is provided", func() {
		message := Message{Key: []byte("any-key")}

		Expect(message.Key).To(Equal([]byte("any-key")))
		Expect(message.Value).To(BeNil())
		Expect(message.Headers).To(BeNil())
	})
})
