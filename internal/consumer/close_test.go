package consumer

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Close", func() {
	var mockReader *readerStub

	BeforeEach(func() {
		mockReader = new(readerStub)
	})

	AfterEach(func() {
		mockReader.AssertExpectations(GinkgoT())
	})

	It("should return nil when reader is nil", func() {
		consumer := &consumerClient{
			reader: nil,
			cfg:    Config{},
		}

		err := consumer.Close()

		Expect(err).To(BeNil())
	})

	It("should return nil when reader is provided and Close is called successfully", func() {
		consumer := &consumerClient{
			reader: mockReader,
			cfg:    Config{},
		}

		mockReader.On("Close").Return(nil)

		err := consumer.Close()

		Expect(err).To(BeNil())
	})

	It("should return error when reader Close returns an error", func() {
		consumer := &consumerClient{
			reader: mockReader,
			cfg:    Config{},
		}

		mockReader.On("Close").Return(errors.New("any error reason returned by the reader"))

		err := consumer.Close()

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("any error reason returned by the reader"))
	})
})
