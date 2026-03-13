# **flowkafka**

[![Go Reference](https://pkg.go.dev/badge/github.com/mauricioromagnollo/flowkafka.svg)](https://pkg.go.dev/github.com/mauricioromagnollo/flowkafka)

A production-ready Kafka worker toolkit for Go with clean architecture support.

## Installation

```sh
go get github.com/mauricioromagnollo/flowkafka
```

## Quick Start

### Producer

```go
producer := flowkafka.NewProducer(flowkafka.ProducerConfig{
    Brokers:      []string{"localhost:9092"},
    TopicName:    "my-topic",
    SaslUsername: "", // optional
    SaslPassword: "", // optional
})
defer producer.Close()

// Validate the connection
if err := producer.ValidateConnection(ctx); err != nil {
    log.Fatal(err)
}

// Send a JSON message
err := producer.Publish(ctx, []byte("key"), map[string]any{"event": "signup"})
if err != nil {
    log.Fatal(err)
}
```

### Consumer (handler-based)

```go
consumer := flowkafka.NewConsumer(flowkafka.ConsumerConfig{
    Brokers:   []string{"localhost:9092"},
    GroupID:   "my-group",
    TopicName: "my-topic",
})
defer consumer.Close()

err := consumer.Consume(ctx, func(msg flowkafka.Message) error {
    fmt.Printf("key=%s value=%s\n", msg.Key, msg.Value)
    return nil // returning nil commits the message
})
```

### Consumer (channel-based)

```go
consumer := flowkafka.NewConsumer(flowkafka.ConsumerConfig{
    Brokers:   []string{"localhost:9092"},
    GroupID:   "my-group",
    TopicName: "my-topic",
})
defer consumer.Close()

msgsChan := make(chan flowkafka.Message)

go func() {
    for msg := range msgsChan {
        fmt.Printf("key=%s value=%s\n", msg.Key, msg.Value)
    }
}()

err := consumer.ConsumeMessages(ctx, msgsChan)
```

### Avro (with Schema Registry)

```go
sr := flowkafka.NewSchemaRegistry(flowkafka.SchemaRegistryConfig{
    Endpoint: "http://localhost:8081",
})

producer := flowkafka.NewProducer(flowkafka.ProducerConfig{
    Brokers:        []string{"localhost:9092"},
    TopicName:      "my-topic",
    SchemaRegistry: sr,
})

err := producer.PublishAvro(ctx, []byte("key"), myAvroStruct)
if err != nil {
    log.Fatal(err)
}
```

## License

[MIT](LICENSE)
