package schemaregistry

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/riferrei/srclient"
)

var _ = Describe("Types", func() {
	It("should return AVRO when SchemaTypeAvro is evaluated", func() {
		Expect(string(SchemaTypeAvro)).To(Equal("AVRO"))
	})

	It("should return JSON when SchemaTypeJSON is evaluated", func() {
		Expect(string(SchemaTypeJSON)).To(Equal("JSON"))
	})

	It("should return PROTOBUF when SchemaTypeProtobuf is evaluated", func() {
		Expect(string(SchemaTypeProtobuf)).To(Equal("PROTOBUF"))
	})

	It("should return srclient.Avro when SchemaTypeAvro is converted to srclient.SchemaType", func() {
		Expect(srclient.SchemaType(SchemaTypeAvro)).To(Equal(srclient.Avro))
	})

	It("should return srclient.Json when SchemaTypeJSON is converted to srclient.SchemaType", func() {
		Expect(srclient.SchemaType(SchemaTypeJSON)).To(Equal(srclient.Json))
	})

	It("should return srclient.Protobuf when SchemaTypeProtobuf is converted to srclient.SchemaType", func() {
		Expect(srclient.SchemaType(SchemaTypeProtobuf)).To(Equal(srclient.Protobuf))
	})

	It("should return the provided values when Config is built", func() {
		config := Config{
			Endpoint:     "http://localhost:8081",
			SaslUsername: "any-username",
			SaslPassword: "any-password",
		}

		Expect(config.Endpoint).To(Equal("http://localhost:8081"))
		Expect(config.SaslUsername).To(Equal("any-username"))
		Expect(config.SaslPassword).To(Equal("any-password"))
	})
})
