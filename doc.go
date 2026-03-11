// Package flowkafka provides a high-level abstraction for working with
// Apache Kafka in Go applications.
//
// It wraps [segmentio/kafka-go] behind clean interfaces so you can produce,
// consume, and manage Kafka messages with minimal boilerplate.
//
// The package exposes two independent clients:
//
//   - [KafkaClient] — produce JSON or Avro messages, consume with manual
//     commit, validate broker connectivity, and manage reader lifecycle.
//   - [SchemaRegistryClient] — create, query, and validate schemas against
//     any Confluent-compatible Schema Registry.
//
// Both clients are defined as interfaces, making them straightforward to
// mock in unit tests.
//
// See the project README and the examples/ directory for runnable demos.
package flowkafka
