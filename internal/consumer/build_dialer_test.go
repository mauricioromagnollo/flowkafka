package consumer

import (
	"crypto/tls"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/segmentio/kafka-go/sasl/plain"
)

var _ = Describe("buildDialer", func() {
	It("should return a dialer without SASL and TLS when sasl password is not provided", func() {
		config := Config{
			SaslUsername: "any-username",
		}

		dialer := buildDialer(config)

		Expect(dialer).ToNot(BeNil())
		Expect(dialer.Timeout).To(Equal(10 * time.Second))
		Expect(dialer.DualStack).To(BeTrue())
		Expect(dialer.SASLMechanism).To(BeNil())
		Expect(dialer.TLS).To(BeNil())
	})

	It("should return a dialer with SASL and TLS when sasl password is provided", func() {
		config := Config{
			SaslUsername: "any-username",
			SaslPassword: "any-password",
		}

		dialer := buildDialer(config)

		Expect(dialer).ToNot(BeNil())
		Expect(dialer.Timeout).To(Equal(10 * time.Second))
		Expect(dialer.DualStack).To(BeTrue())
		Expect(dialer.SASLMechanism).To(Equal(plain.Mechanism{
			Username: "any-username",
			Password: "any-password",
		}))
		Expect(dialer.TLS).ToNot(BeNil())
		Expect(dialer.TLS.MinVersion).To(Equal(uint16(tls.VersionTLS12)))
	})
})
