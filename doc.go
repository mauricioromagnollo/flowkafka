// Package flowkafka provides a high-level abstraction for working with
// Apache Kafka in Go applications.
//
// It wraps [github.com/segmentio/kafka-go] behind clean interfaces so you can
// produce, consume, and manage Kafka messages with minimal boilerplate.
//
// The package exposes three independent clients:
//
//   - [Producer] — publish JSON messages with Publish or Avro messages in the
//     Confluent wire format with PublishAvro, validate broker connectivity,
//     and manage the writer lifecycle.
//   - [Consumer] — consume messages with consumer group support and manual
//     commit, either through a handler function with Consume or through a
//     channel with ConsumeMessages.
//   - [SchemaRegistry] — create, query, and validate schemas against any
//     Confluent-compatible Schema Registry.
//
// All three are defined as interfaces, making them straightforward to mock in
// unit tests.
//
// Create clients with [NewProducer], [NewConsumer], and [NewSchemaRegistry].
//
// See the project README for runnable examples.
package flowkafka
