package brokerclient

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BrokerClientStub", func() {
	var stub *BrokerClientStub

	BeforeEach(func() {
		stub = new(BrokerClientStub)
	})

	AfterEach(func() {
		stub.AssertExpectations(GinkgoT())
	})

	It("should satisfy the BrokerClient interface when used as a dependency", func() {
		var client BrokerClient = stub

		Expect(client).ToNot(BeNil())
	})

	It("should return the configured metadata when GetMetadata is expected to succeed", func() {
		ctx := context.Background()
		expected := &ClientMetadata{
			Topics: []TopicMetadata{{Name: "any-topic"}},
		}

		stub.On("GetMetadata", ctx).Return(expected, nil)

		metadata, err := stub.GetMetadata(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(metadata).To(BeIdenticalTo(expected))
		Expect(metadata.Topics).To(HaveLen(1))
		Expect(metadata.Topics[0].Name).To(Equal("any-topic"))
	})

	It("should return the configured error when GetMetadata is expected to fail", func() {
		ctx := context.Background()
		expectedErr := errors.New("any error reason returned by the broker")

		stub.On("GetMetadata", ctx).Return(&ClientMetadata{}, expectedErr)

		metadata, err := stub.GetMetadata(ctx)

		Expect(err).To(MatchError(expectedErr))
		Expect(metadata).To(Equal(&ClientMetadata{}))
	})

	It("should return the metadata for the given context when a custom context is used", func() {
		ctx := context.WithValue(context.Background(), struct{}{}, "any-value")
		expected := &ClientMetadata{}

		stub.On("GetMetadata", ctx).Return(expected, nil)

		metadata, err := stub.GetMetadata(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(metadata).To(BeIdenticalTo(expected))
	})
})
