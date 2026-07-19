package transport

import (
	"crypto/tls"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/segmentio/kafka-go/sasl/plain"
)

var _ = Describe("NewTransport", func() {
	It("should return a non-nil transport when config is empty", func() {
		transport := NewTransport(Config{})

		Expect(transport).ToNot(BeNil())
	})

	It("should return a transport with 10s dial timeout when config is empty", func() {
		transport := NewTransport(Config{})

		Expect(transport.DialTimeout).To(Equal(10 * time.Second))
	})

	It("should return a transport without SASL and TLS when sasl password is empty", func() {
		transport := NewTransport(Config{
			SaslUsername: "any-username",
			SaslPassword: "",
		})

		Expect(transport.SASL).To(BeNil())
		Expect(transport.TLS).To(BeNil())
	})

	It("should return a transport with plain SASL mechanism when sasl password is provided", func() {
		transport := NewTransport(Config{
			SaslUsername: "any-username",
			SaslPassword: "any-password",
		})

		mechanism, ok := transport.SASL.(plain.Mechanism)

		Expect(ok).To(BeTrue())
		Expect(mechanism.Username).To(Equal("any-username"))
		Expect(mechanism.Password).To(Equal("any-password"))
	})

	It("should return a transport with TLS 1.2 as minimum version when sasl password is provided", func() {
		transport := NewTransport(Config{
			SaslUsername: "any-username",
			SaslPassword: "any-password",
		})

		Expect(transport.TLS).ToNot(BeNil())
		Expect(transport.TLS.MinVersion).To(Equal(uint16(tls.VersionTLS12)))
	})

	It("should return a transport with 10s dial timeout when sasl password is provided", func() {
		transport := NewTransport(Config{
			SaslUsername: "any-username",
			SaslPassword: "any-password",
		})

		Expect(transport.DialTimeout).To(Equal(10 * time.Second))
	})

	// The implementation only checks SaslPassword, so an empty username still
	// enables SASL. This spec locks that real behavior.
	It("should return a transport with SASL configured when sasl password is provided and username is empty", func() {
		transport := NewTransport(Config{
			SaslUsername: "",
			SaslPassword: "any-password",
		})

		mechanism, ok := transport.SASL.(plain.Mechanism)

		Expect(ok).To(BeTrue())
		Expect(mechanism.Username).To(BeEmpty())
		Expect(mechanism.Password).To(Equal("any-password"))
		Expect(transport.TLS).ToNot(BeNil())
	})
})
