package flowkafka

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFlowkafka(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flowkafka Suite")
}
