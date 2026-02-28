package flowkafka_test

import (
	"testing"

	"github.com/mauricioromagnollo/flowkafka"
)

func TestPublicAPICompiles(t *testing.T) {
	_ = flowkafka.KafkaConfig{}
	_ = flowkafka.Message{}
}
