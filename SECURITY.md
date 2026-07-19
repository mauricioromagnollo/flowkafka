# Security Policy

`flowkafka` is a Go library that talks to Apache Kafka brokers and Confluent-compatible
Schema Registries. Because it handles broker credentials and transport configuration,
security reports are taken seriously and handled privately until a fix is available.

## Supported Versions

Only the latest released minor version receives security fixes. Because this is a
single-maintainer project, older lines are not backported — please stay current.

| Version | Supported |
| --- | --- |
| `1.0.x` (latest release) | ✅ |
| `0.1.x` | ❌ |
| Unreleased `main` | ✅ (fixes land here first) |

Upgrading is a one-line change:

```sh
go get github.com/mauricioromagnollo/flowkafka@latest
```

## Reporting a Vulnerability

**Please do not open a public GitHub issue, pull request, or discussion for a security
vulnerability.** Public disclosure before a fix exists puts every user of the library at
risk.

### Preferred channel — GitHub Private Vulnerability Reporting

1. Go to <https://github.com/mauricioromagnollo/flowkafka/security>.
2. Click **Report a vulnerability**.
3. Fill in the advisory form with as much detail as you can.

This keeps the report, the discussion, and the eventual CVE and advisory in one place,
and lets you be credited automatically when the advisory is published.

### Alternative channel — email

If you cannot use GitHub Private Vulnerability Reporting, email
<mauricioromagnollo@gmail.com> with `[SECURITY] flowkafka` in the subject line. If you
would like to encrypt the report, say so in a first plain-text message and a key will be
exchanged.

## What to Include in a Report

The more of this you can provide, the faster the fix:

- **Description** of the vulnerability and the security impact (what an attacker gains).
- **Affected versions** — module version, tag, or commit SHA.
- **Environment** — Go version, OS, Kafka distribution and version, whether SASL, TLS, or
  a Schema Registry is in play.
- **Reproduction steps** — ideally a minimal Go program that demonstrates the issue.
- **Proof of concept**, if you have one.
- **Suggested mitigation or patch**, if you have one in mind.
- **Disclosure preferences** — how you would like to be credited, and any deadline you
  are working to.

## Response Timeline

This project is maintained by one person in their own time. The commitments below are
what can realistically be honored:

| Stage | Target |
| --- | --- |
| Acknowledgment of your report | Within 48 hours |
| Initial assessment (valid / not valid, severity) | Within 7 days |
| Status update while a fix is being developed | At least every 7 days |
| Fix released for a confirmed high-severity issue | Target 30 days from confirmation |

If you have not heard back within 72 hours, please send a polite follow-up — mail does
get filtered occasionally.

## Disclosure Policy

This project follows **coordinated disclosure**:

1. You report privately through one of the channels above.
2. The report is confirmed and a fix is developed in private (a temporary private fork,
   when using GitHub Security Advisories).
3. A patched version is released and tagged.
4. A GitHub Security Advisory is published, crediting you unless you ask to stay
   anonymous, and a CVE is requested when the issue warrants one.
5. Public disclosure happens after the fix is released, or after **90 days** from the
   report, whichever comes first.

Reporters who follow this process will not be pursued or penalized in any way for their
research. There is no paid bug bounty for this project.

## Security Considerations for Users of the Library

The following are not vulnerabilities in `flowkafka` itself, but they are the most common
ways deployments using it get compromised. Please review them before shipping.

### SASL/PLAIN sends credentials over the wire

`flowkafka` configures SASL/PLAIN when `SaslUsername` and `SaslPassword` are set on
`ProducerConfig`, `ConsumerConfig`, or `SchemaRegistryConfig`. **SASL/PLAIN transmits the
username and password as effectively cleartext** — it only becomes safe when the
connection underneath it is encrypted.

The library enables TLS (minimum version TLS 1.2) on the Kafka transport whenever a SASL
password is provided, so credentials are not sent in the clear against a TLS-enabled
listener. Two consequences follow:

- **Never point a SASL-configured client at a `PLAINTEXT` listener.** If the broker's
  listener does not speak TLS, the handshake will fail rather than silently downgrade —
  do not "fix" this by stripping credentials or terminating TLS somewhere untrusted.
- **Use an `https://` endpoint for the Schema Registry** whenever you pass SASL
  credentials to `SchemaRegistryConfig`. Plain `http://` sends basic-auth credentials in
  the clear.

Prefer a stronger mechanism (SCRAM-SHA-256/512, mTLS) at the broker level where your
platform offers it.

### Never hardcode credentials

Do not embed broker usernames, passwords, or Schema Registry credentials in source code,
container images, or committed configuration files.

```go
// Do NOT do this
producer := flowkafka.NewProducer(flowkafka.ProducerConfig{
    Brokers:      []string{"broker.example.com:9093"},
    TopicName:    "orders",
    SaslUsername: "svc-orders",
    SaslPassword: "hunter2", // committed to git, forever
})

// Do this instead
producer := flowkafka.NewProducer(flowkafka.ProducerConfig{
    Brokers:      strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
    TopicName:    os.Getenv("KAFKA_TOPIC"),
    SaslUsername: os.Getenv("KAFKA_SASL_USERNAME"),
    SaslPassword: os.Getenv("KAFKA_SASL_PASSWORD"),
})
```

Load credentials from environment variables injected at runtime, or from a secret manager
(AWS Secrets Manager, HashiCorp Vault, GCP Secret Manager, Kubernetes Secrets). Rotate
them regularly, and give each service its own principal with the narrowest ACLs that let
it do its job — produce to its topics, consume from its group, and nothing more.

### Message content and logging

- Errors returned by this library are wrapped with `%w` and may include broker addresses
  and topic names. Treat error strings as operational data; do not echo them to end
  users.
- `Message.Value` may contain personal or regulated data. Avoid logging raw payloads, and
  scrub them if you must.
- Untrusted message payloads are still untrusted input after deserialization. Validate
  them in your handler before acting on them.

### Dependencies

`flowkafka` depends on `segmentio/kafka-go`, `riferrei/srclient`, and `hamba/avro`. Keep
your module graph current and run `govulncheck ./...` in your own CI:

```sh
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

If a vulnerability in one of these dependencies affects `flowkafka` users, report it
through the channels above so a bumped release can be cut.
