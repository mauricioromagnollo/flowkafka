package producer

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kafkago "github.com/segmentio/kafka-go"
)

var _ = Describe("Types", func() {
	It("should return 0 when RequiredAcksNone is evaluated", func() {
		Expect(int(RequiredAcksNone)).To(Equal(0))
	})

	It("should return 1 when RequiredAcksOne is evaluated", func() {
		Expect(int(RequiredAcksOne)).To(Equal(1))
	})

	It("should return -1 when RequiredAcksAll is evaluated", func() {
		Expect(int(RequiredAcksAll)).To(Equal(-1))
	})

	It("should return kafkago.RequireNone when RequiredAcksNone is converted", func() {
		Expect(kafkago.RequiredAcks(RequiredAcksNone)).To(Equal(kafkago.RequireNone))
	})

	It("should return kafkago.RequireOne when RequiredAcksOne is converted", func() {
		Expect(kafkago.RequiredAcks(RequiredAcksOne)).To(Equal(kafkago.RequireOne))
	})

	It("should return kafkago.RequireAll when RequiredAcksAll is converted", func() {
		Expect(kafkago.RequiredAcks(RequiredAcksAll)).To(Equal(kafkago.RequireAll))
	})
})
