package producer

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Close", func() {
	var mockWriter *writerStub

	BeforeEach(func() {
		mockWriter = new(writerStub)
	})

	AfterEach(func() {
		mockWriter.AssertExpectations(GinkgoT())
	})

	It("should return nil when writer is nil", func() {
		config := Config{}
		producer := &producerClient{
			writer:         nil,
			cfg:            config,
			schemaRegistry: nil,
		}

		err := producer.Close()

		Expect(err).To(BeNil())
	})

	It("should return nil when writer is provided and Close is called successfully", func() {
		config := Config{}
		producer := &producerClient{
			writer:         mockWriter,
			cfg:            config,
			schemaRegistry: nil,
		}

		mockWriter.On("Close").Return(nil)

		err := producer.Close()

		Expect(err).To(BeNil())
	})

	It("should return error when writer is provided and Close returns an error", func() {
		config := Config{}
		producer := &producerClient{
			writer:         mockWriter,
			cfg:            config,
			schemaRegistry: nil,
		}

		mockWriter.On("Close").Return(errors.New("any error reason returned by the writer"))

		err := producer.Close()

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("any error reason returned by the writer"))
	})
})
