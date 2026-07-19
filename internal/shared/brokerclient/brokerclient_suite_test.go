package brokerclient

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBrokerClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BrokerClient Suite")
}
