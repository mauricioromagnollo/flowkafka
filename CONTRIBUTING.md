# Contributing to flowkafka

First off, thank you for taking the time to contribute. `flowkafka` is a Kafka worker
toolkit for Go, and every bug report, documentation fix, and pull request helps make it
better.

This document explains how to get set up, the conventions the codebase follows, and what
is expected of a pull request.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Reporting a Bug](#reporting-a-bug)
- [Requesting a Feature](#requesting-a-feature)
- [Development Setup](#development-setup)
- [Makefile Commands](#makefile-commands)
- [Code Style](#code-style)
- [Testing Guidelines](#testing-guidelines)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Adding a New Producer or Consumer Method](#adding-a-new-producer-or-consumer-method)
- [Release Process](#release-process)

## Code of Conduct

This project and everyone participating in it is governed by the
[Code of Conduct](CODE_OF_CONDUCT.md), which is based on the Contributor Covenant v2.1.
By participating, you are expected to uphold this code. Please report unacceptable
behavior to <mauricioromagnollo@gmail.com>.

## Reporting a Bug

Before opening a new report, please search the
[existing issues](https://github.com/mauricioromagnollo/flowkafka/issues) — someone may
have hit the same problem already.

If nothing matches, open a new issue using the bug report template at
[New issue → Bug report](https://github.com/mauricioromagnollo/flowkafka/issues/new/choose)
and include:

- The version of `flowkafka` (module version or commit SHA) and your Go version.
- The Kafka distribution and version you are running against (Apache Kafka, Redpanda,
  MSK, Confluent Cloud, ...), and whether SASL/TLS or a Schema Registry is involved.
- A minimal reproducible example — the smallest `main.go` that shows the problem.
- What you expected to happen and what actually happened, with the full error message.
  Because every error in this library is wrapped with `%w`, please paste the entire
  chain, not just the outermost message.

**Do not report security vulnerabilities through public issues.** Follow
[SECURITY.md](SECURITY.md) instead.

## Requesting a Feature

Open a feature request at
[New issue → Feature request](https://github.com/mauricioromagnollo/flowkafka/issues/new/choose)
and describe:

- The problem you are trying to solve, not only the API you have in mind.
- How you would use it — a snippet of the call site you wish existed.
- Any alternatives you considered, including whether the underlying
  [`segmentio/kafka-go`](https://github.com/segmentio/kafka-go) API already covers it.

`flowkafka` deliberately keeps a small public surface. Proposals that add configuration
knobs which map one-to-one to the underlying driver are less likely to be accepted than
proposals that remove boilerplate from user code.

Please open an issue and get agreement on the design **before** submitting a large pull
request. Small fixes and documentation improvements can go straight to a PR.

## Development Setup

### Prerequisites

| Tool | Version | Why |
| --- | --- | --- |
| Go | `1.26.0` (see [`.go-version`](.go-version) / [`.tool-versions`](.tool-versions)) | Building and testing |
| Docker | any recent version | Runs `golangci-lint` via `make lint` and the local Kafka stack via `make up` |
| GNU Make | any recent version | Task runner |
| CGO toolchain | system default | Required by `make test-race` (`CGO_ENABLED=1`) |

If you use [asdf](https://asdf-vm.com/) or [mise](https://mise.jdx.dev/), the pinned Go
version is picked up automatically from `.tool-versions`.

### Fork and clone

```sh
# 1. Fork https://github.com/mauricioromagnollo/flowkafka on GitHub, then:
git clone https://github.com/<your-username>/flowkafka.git
cd flowkafka

# 2. Keep a reference to the canonical repository
git remote add upstream https://github.com/mauricioromagnollo/flowkafka.git

# 3. Download dependencies
go mod download
```

### Verify your setup

```sh
make test
make lint
```

Both must pass on a clean checkout before you start changing anything. If they do not,
open an issue — that is a bug in itself.

### Local Kafka

Some manual verification requires a running broker. `make up` starts the local stack
(Kafka and Schema Registry) with Docker Compose in detached mode, and `make down` stops
it while keeping images, networks, and volumes around for the next run.

## Makefile Commands

Run `make help` at any time to print the list below directly from the Makefile.

| Command | What it does |
| --- | --- |
| `make help` | Print all available commands |
| `make test` | Run the full test suite (`go test -v ./...`) |
| `make test-coverage` | Run tests with race detection and print per-function coverage |
| `make test-coverage-web` | Same as above, then open the HTML coverage report |
| `make test-race` | Run only the data race detector (`CGO_ENABLED=1 go test -race ./...`) |
| `make lint` | Run `golangci-lint` inside Docker, pinned to the same version CI uses |
| `make up` | Start the local Kafka stack in detached mode |
| `make down` | Stop the containers, keeping images, networks, and volumes |
| `make stop` | Stop containers without removing them |
| `make clear` | Remove containers, images, networks, and volumes |

`make lint` runs the linter in a container (`golangci/golangci-lint:v2.9.0`) so that your
local result matches CI exactly. You do not need `golangci-lint` installed on your host.

## Code Style

The linter configuration in [`.golangci.yml`](.golangci.yml) is the source of truth.
Beyond it, the codebase follows a few explicit conventions — please match them, since a
consistent layout is what makes this repository quick to navigate.

### One exported method per file (internal packages)

Inside `internal/`, every exported method of a type lives in its own file, and the file
name is the method name in `snake_case`:

```
internal/producer/publish.go              -> func (c *producerClient) Publish(...)
internal/producer/publish_avro.go         -> func (c *producerClient) PublishAvro(...)
internal/producer/validate_connection.go  -> func (c *producerClient) ValidateConnection(...)
internal/producer/close.go                -> func (c *producerClient) Close(...)
internal/producer/has_schema_registry.go  -> func (c *producerClient) HasSchemaRegistry(...)
```

Supporting files follow the same idea:

- `<package>.go` (e.g. `producer.go`, `consumer.go`) holds the package's interface.
- `new.go` holds the constructor and the concrete struct.
- `types.go` holds `Config` and other shared types and constants.
- `*_stub.go` holds test doubles.

### `doc.go` in every package

Every package — including new ones — must have a `doc.go` containing the package comment:

```go
// Package producer provides a Kafka producer implementation that supports retry logic and clean shutdown.
// It allows producing messages to Kafka topics with optional integration with a schema registry for message serialization.
package producer
```

The `revive` `exported` rule is enabled, so every exported identifier also needs its own
doc comment starting with the identifier's name.

### Error handling

Every returned error must be wrapped with context using `%w`, never `%v` and never
returned bare from an external call:

```go
if err := c.writer.WriteMessages(ctx, message); err != nil {
    return fmt.Errorf("failed to publish message in kafka: %w", err)
}
```

Guidelines for the message itself:

- Lowercase, no trailing punctuation.
- Describe the operation that failed: `failed to marshal message to publish`.
- Do not include the underlying error text manually — `%w` already carries it.

`errorlint`, `nilerr`, and `nilnil` are enabled and will catch most deviations.

### Formatting and linting

- Code must be `gofmt`/`gofumpt` clean; imports are grouped by `gci` and `goimports`.
- `make lint` must pass with zero issues. CI runs the exact same command.
- Dot imports are only allowed for `github.com/onsi/ginkgo/v2` and
  `github.com/onsi/gomega` in test files.
- `//nolint` directives are permitted but should be rare — prefer fixing the finding.

### Public API surface

The root package (`flowkafka`) is a thin façade. It only contains:

- Type aliases re-exporting internal types (`type Producer = producer.Producer`).
- Constructors (`NewProducer`, `NewConsumer`, `NewSchemaRegistry`).
- Wrapper structs that delegate straight to the internal client.

No business logic belongs in the root package. Anything that can fail, branch, or be
tested goes into `internal/`.

## Testing Guidelines

The project uses **[Ginkgo v2](https://onsi.github.io/ginkgo/) +
[Gomega](https://onsi.github.io/gomega/)** for specs and
**[testify/mock](https://github.com/stretchr/testify)** for stubs.

### Layout

Tests live in the **same package** as the code under test (white-box testing) for
everything under `internal/`. This is deliberate: it lets specs construct the unexported
client struct directly and inject stubs.

Each internal package has:

1. **One suite bootstrap file**, `<package>_suite_test.go`, that wires Ginkgo into
   `go test`:

   ```go
   package producer

   import (
       "testing"

       . "github.com/onsi/ginkgo/v2"
       . "github.com/onsi/gomega"
   )

   func TestProducer(t *testing.T) {
       RegisterFailHandler(Fail)
       RunSpecs(t, "Producer Suite")
   }
   ```

2. **One spec file per method**, named after the file it tests (`publish.go` →
   `publish_test.go`), with a single top-level `Describe` named after the method:

   ```go
   var _ = Describe("Publish", func() {
       var mockWriter *writerStub

       BeforeEach(func() {
           mockWriter = new(writerStub)
       })

       AfterEach(func() {
           mockWriter.AssertExpectations(GinkgoT())
       })

       It("should return error when fail to write message", func() {
           // ...
       })
   })
   ```

3. **Stubs in `*_stub.go` files** (not `_test.go`), embedding the interface plus
   `mock.Mock`:

   ```go
   type writerStub struct {
       Writer
       mock.Mock
   }

   func (m *writerStub) WriteMessages(ctx context.Context, msg types.Message) error {
       args := m.Called(ctx, msg)

       return args.Error(0)
   }
   ```

Tests for the root package (`import_test.go` and friends) exist to assert that the public
API compiles and that the aliases stay wired to the internal types.

### Expectations for new contributions

- Every new method or bug fix must come with specs. A bug fix should include a spec that
  fails before the fix.
- Cover the **error paths**, not just the happy path. Each `return fmt.Errorf(...)` in a
  method deserves an `It(...)`.
- Assert on wrapped error messages with `ContainSubstring`, so the wrapping contract is
  itself under test.
- New code should not lower package coverage; check with `make test-coverage`. Aim for
  full coverage of new internal code — the packages are small enough that this is
  realistic.
- Anything spawning goroutines or channels must pass `make test-race`.
- Spec names read as sentences: `It("should return nil when message is published successfully", ...)`.
- Unit tests must not require a running broker. Talk to interfaces and inject stubs.

## Commit Messages

This project follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/),
in English, in the imperative mood, with a lowercase subject and no trailing period.

```
<type>(<optional scope>): <subject>

<optional body>

<optional footer>
```

### Types used in this repository

| Type | Use for |
| --- | --- |
| `feat` | A new feature or public API addition |
| `fix` | A bug fix |
| `docs` | Documentation only |
| `test` | Adding or fixing tests |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `chore` | Repository housekeeping (`.gitignore`, templates, tooling) |
| `ci` | CI configuration and workflows |
| `perf` | Performance improvement |
| `build` | Build system or dependencies (e.g. `build(deps): ...`) |
| `style` | Formatting and linter configuration, no behavior change |

### Examples from the project history

```
feat: add consumer package
fix: update APP_NAME to flowkafka and improve lint command syntax
docs: update README with installation and usage instructions for producer and consumer
test: add missing public api tests
refactor: remove unused Kafka configuration and message structs
chore: add PR template file
ci: add ci pipeline to run in pull requests
build(deps): add required kafka and schemaregistry deps
style: add complete golangci config file
```

### Breaking changes

Mark them with a `!` after the type and a `BREAKING CHANGE:` footer:

```
feat(producer)!: rename ProduceJSONMessage to Publish

BREAKING CHANGE: ProduceJSONMessage and ProduceAvroMessage were renamed to
Publish and PublishAvro. Update call sites accordingly.
```

## Pull Request Process

1. **Sync and branch from `main`.**

   ```sh
   git checkout main
   git pull upstream main
   git checkout -b feat/consumer-dead-letter-topic
   ```

   Branch names follow `<type>/<short-kebab-description>`, using the same types as commit
   messages — for example `fix/producer-close-nil-writer`,
   `docs/open-source-documentation`, `test/consumer-specs`.

2. **One subject per pull request.** Do not mix a refactor, a bug fix, and a docs update
   in the same PR. Small, focused PRs get reviewed and merged faster.

3. **Run the full check locally before pushing:**

   ```sh
   make lint
   make test
   make test-race
   ```

4. **Fill in the pull request template** completely. Link the related issue (`Issue: #123`);
   if no issue exists for a non-trivial change, open one first.

5. **CI must be green.** The [CI workflow](.github/workflows/ci.yaml) runs `make lint` and
   `make test` on every pull request. PRs with a red pipeline are not reviewed.

6. **Update the docs.** If your change touches the public API, update `docs/README.md`
   and the relevant Go doc comments in the same PR. Add an entry under `## [Unreleased]`
   in [`CHANGELOG.md`](CHANGELOG.md).

7. **Respond to review.** Push follow-up commits rather than force-pushing during review,
   so reviewers can see what changed. The maintainer squashes on merge, so the final
   history stays clean.

By submitting a pull request, you agree that your contribution is licensed under the
[MIT License](LICENSE).

## Adding a New Producer or Consumer Method

The layering is intentional: interface → implementation → public delegator. Adding a
method touches exactly three or four files.

Using a hypothetical `PublishBatch` on the producer as an example:

1. **Declare it on the internal interface** in `internal/producer/producer.go`, with a
   doc comment:

   ```go
   type Producer interface {
       // ...
       // PublishBatch sends multiple JSON messages to a Kafka topic in a single write.
       PublishBatch(ctx context.Context, msgs []any) error
   }
   ```

2. **Implement it in a new file** named after the method — `internal/producer/publish_batch.go`:

   ```go
   package producer

   // PublishBatch sends multiple JSON messages to a Kafka topic in a single write.
   func (c *producerClient) PublishBatch(ctx context.Context, msgs []any) error {
       // ... wrap every error with %w
   }
   ```

3. **Write the spec** in `internal/producer/publish_batch_test.go`, covering the happy
   path and each error branch, using the existing `writerStub`. Extend the stub only if
   the new method needs a driver call that is not stubbed yet.

4. **Add the delegator to the root wrapper** in `producer.go`:

   ```go
   func (p producerWrapper) PublishBatch(ctx context.Context, msgs []any) error {
       return p.client.PublishBatch(ctx, msgs)
   }
   ```

   The root wrapper contains no logic — it forwards, and nothing else. If you find
   yourself writing an `if` there, it belongs in `internal/`.

5. If the method introduces new public types or constants, re-export them from the root
   package as aliases (as done for `Message`, `RequiredAcks`, and `SchemaType`), and
   update `docs/README.md` with a usage snippet.

The same steps apply to `internal/consumer` (`consumer.go` → new file → spec →
`consumer.go` at the root) and `internal/schemaregistry`.

## Release Process

Releases are cut by the maintainer.

- The project follows [Semantic Versioning 2.0.0](https://semver.org/):
  - **MAJOR** — breaking change to the public `flowkafka` API.
  - **MINOR** — backwards-compatible functionality.
  - **PATCH** — backwards-compatible bug fixes.
- Releases are published as annotated Git tags in the `vX.Y.Z` form (`v1.0.0`), which is
  what the Go module proxy consumes. There is no separate publishing step: tagging is
  the release.
- Before tagging, the `## [Unreleased]` section of [`CHANGELOG.md`](CHANGELOG.md) is
  promoted to a versioned section with the release date, and the comparison links at the
  bottom of the file are updated.
- Because the module path has no `/vN` suffix, a `v2.0.0` or later release would require
  changing the module path in `go.mod`. Breaking changes are therefore batched and
  avoided whenever an additive design is possible.

## Questions

Open a [discussion or issue](https://github.com/mauricioromagnollo/flowkafka/issues) —
questions about how to use the library are welcome and usually reveal a documentation
gap worth fixing.
