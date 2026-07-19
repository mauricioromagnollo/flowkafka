# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `internal/shared` package tree to remove duplication between the producer and the
  consumer:
  - `internal/shared/transport` — builds the `kafka-go` transport with optional
    SASL/PLAIN authentication and TLS (minimum TLS 1.2), shared by both clients.
  - `internal/shared/brokerclient` — broker metadata client used by connection
    validation, with a `client_stub.go` test double.
  - `internal/shared/types` — the shared `Message` type.
- `Writer` interface and `kafkago_writer.go` adapter in `internal/producer`, isolating
  the producer from the `segmentio/kafka-go` writer and making it stubbable.
- Ginkgo v2 + Gomega test suites for the producer package (`producer_suite_test.go`),
  covering `Publish`, `Close`, and `HasSchemaRegistry`, with a `testify/mock` based
  `writerStub`.
- Tests for the public API surface and for the `Message` type.
- Optional producer tuning options on `ProducerConfig`: `MaxAttempts`, `BatchTimeout`,
  `WriteTimeout`, `ReadTimeout`, `RequiredAcks`, and `ClientTimeout`.
- Public `RequiredAcks` type with the `RequiredAcksNone`, `RequiredAcksOne`, and
  `RequiredAcksAll` constants.

### Changed

- **Breaking:** renamed the producer methods `ProduceJSONMessage` to `Publish` and
  `ProduceAvroMessage` to `PublishAvro`, on both the `Producer` interface and the root
  package wrapper.
- Moved the public `Message` alias from `message.go` to `types.go` in the root package.
- Connection validation for the producer and the consumer now goes through the shared
  broker metadata client instead of each package building its own dialer and transport.

### Removed

- Per-package `build_transport.go` files in `internal/producer` and `internal/consumer`,
  superseded by `internal/shared/transport`.

## [1.0.0] - 2026-03-11

First stable release. The public API is exposed by the root `flowkafka` package as thin
wrappers over the `internal` implementations.

### Added

- **Producer** (`NewProducer`, `ProducerConfig`):
  - `ProduceJSONMessage` for publishing JSON-encoded messages.
  - `ProduceAvroMessage` for publishing Avro messages in the Confluent wire format via a
    Schema Registry.
  - `ValidateConnection` to check broker reachability.
  - `HasSchemaRegistry` to report whether a Schema Registry client is configured.
  - `Close` for clean shutdown.
- **Consumer** (`NewConsumer`, `ConsumerConfig`):
  - `Consume` — handler-based consumption; a message is committed only after the handler
    returns `nil`.
  - `ConsumeMessages` — channel-based consumption for pipeline-style workers.
  - `ValidateConnection` and `Close`.
- **Schema Registry** (`NewSchemaRegistry`, `SchemaRegistryConfig`) backed by
  `riferrei/srclient`:
  - `GetLatestSchema`, `CreateNewSchema`, and `ValidateConnection`.
  - `SchemaType` with the `SchemaTypeAvro`, `SchemaTypeJSON`, and `SchemaTypeProtobuf`
    constants.
- SASL/PLAIN authentication with TLS (minimum TLS 1.2) enabled automatically whenever a
  SASL password is configured, for the producer, the consumer, and the Schema Registry
  client.
- Package documentation (`doc.go`) for the root package and every internal package.
- `docs/README.md` with installation instructions and Quick Start examples for the
  producer, both consumer styles, and Avro with Schema Registry.
- Project tooling: `Makefile` with test, coverage, race, lint, and Docker Compose
  targets; full `golangci-lint` configuration; `.editorconfig`; pull request template;
  and a CI workflow running lint and tests on every pull request.
- MIT License.

### Changed

- Reorganized the codebase into `internal/producer`, `internal/consumer`, and
  `internal/schemaregistry`, with one exported method per file, replacing the earlier
  single `kafka` package.

### Removed

- The combined `kafka` client package and the unused Kafka configuration and message
  structs from the initial prototype.

## [0.1.0] - 2026-02-28

Initial tagged prototype, exposing only the public types.

### Added

- Public `KafkaConfig` type.
- Public `Message` type.
- Go module `github.com/mauricioromagnollo/flowkafka` with the Go version pinned via
  `.go-version` and `.tool-versions`.

[Unreleased]: https://github.com/mauricioromagnollo/flowkafka/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/mauricioromagnollo/flowkafka/compare/v0.1.0...v1.0.0
[0.1.0]: https://github.com/mauricioromagnollo/flowkafka/releases/tag/v0.1.0
