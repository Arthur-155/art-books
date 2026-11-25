# Exemplos executáveis do SDK — SDK gofi/v0.2.18 (c5413b9)

> Gerado do SDK — não edite. Programas completos, compilados e testados na CI do SDK:
> leia o código no checkout, na versão que o projeto fixou.

## netx/grpcx — gRPC example

`.gofi/gofi-sdk-go/examples/grpc/`

A `people.v1.PeopleService` that moves a `Person` entity and its photo through every gRPC call shape, with two servers to choose from and one client that works with both:

- `.gofi/gofi-sdk-go/examples/grpc/generate.go`

## iam — login example

`.gofi/gofi-sdk-go/examples/iam/login/`

Login with the `iam` package in its three usual shapes, with **no database**: users, sessions and tokens live in memory.

- `.gofi/gofi-sdk-go/examples/iam/login/auth.go`
- `.gofi/gofi-sdk-go/examples/iam/login/handlers.go`
- `.gofi/gofi-sdk-go/examples/iam/login/main.go`
- `.gofi/gofi-sdk-go/examples/iam/login/users.go`

## Kafka consumer

`.gofi/gofi-sdk-go/examples/msq/kafka/consumer/`

Processes the `OrderCreated` events from `orders`. Invalid orders are retried and then dead-lettered to `orders-dlq`, which a second handler logs. Part of the [Kafka example](../README.md).

- `.gofi/gofi-sdk-go/examples/msq/kafka/consumer/main.go`

## Kafka producer

`.gofi/gofi-sdk-go/examples/msq/kafka/producer/`

Publishes an `OrderCreated` event to `orders` every 2 seconds until `Ctrl+C`. Part of the [Kafka example](../README.md); run the consumer too to see the messages processed.

- `.gofi/gofi-sdk-go/examples/msq/kafka/producer/main.go`

## RabbitMQ consumer

`.gofi/gofi-sdk-go/examples/msq/rabbitmq/consumer/`

Processes the `OrderCreated` events from `orders`. Invalid orders are retried and then dead-lettered to `orders-dlq`, which a second handler logs. Part of the [RabbitMQ example](../README.md).

- `.gofi/gofi-sdk-go/examples/msq/rabbitmq/consumer/main.go`

## RabbitMQ producer

`.gofi/gofi-sdk-go/examples/msq/rabbitmq/producer/`

Publishes an `OrderCreated` event to `orders` every 2 seconds until `Ctrl+C`. Part of the [RabbitMQ example](../README.md); run the consumer too to see the messages processed.

- `.gofi/gofi-sdk-go/examples/msq/rabbitmq/producer/main.go`

## Amazon SQS consumer

`.gofi/gofi-sdk-go/examples/msq/sqs/consumer/`

Processes the `OrderCreated` events from `orders`. Invalid orders are retried and then dead-lettered to `orders-dlq`, which a second handler logs. Part of the [Amazon SQS example](../README.md).

- `.gofi/gofi-sdk-go/examples/msq/sqs/consumer/main.go`

## Amazon SQS producer

`.gofi/gofi-sdk-go/examples/msq/sqs/producer/`

Publishes an `OrderCreated` event to `orders` every 2 seconds until `Ctrl+C`. Part of the [Amazon SQS example](../README.md); run the consumer too to see the messages processed.

- `.gofi/gofi-sdk-go/examples/msq/sqs/producer/main.go`

## netx/httpx — API example

`.gofi/gofi-sdk-go/examples/netx/api/`

A minimal HTTP API built with `netx/httpx`: several handlers, public and private routes, a global middleware and an auth middleware, all wired in `main.go` through the gofi `httpserver` component.

- `.gofi/gofi-sdk-go/examples/netx/api/main.go`

## obs — observability example

`.gofi/gofi-sdk-go/examples/obs/`

An instrumented service that exports **traces, metrics and logs** over OTLP to a local Grafana stack. It covers what a real service has to instrument:

- `.gofi/gofi-sdk-go/examples/obs/loadgen.go`
- `.gofi/gofi-sdk-go/examples/obs/main.go`
- `.gofi/gofi-sdk-go/examples/obs/runner.go`

## sqln — dynamic filter API

`.gofi/gofi-sdk-go/examples/sqln/filter-api/`

An HTTP API where the **client chooses the filters**: the JSON body is turned into SQL by `sqln`, through an allowlist (`sqln.FilterMapping`). Fields, operators or sort columns outside the allowlist are rejected with `400`.

- `.gofi/gofi-sdk-go/examples/sqln/filter-api/main.go`

## sqln — full-text search job

`.gofi/gofi-sdk-go/examples/sqln/fulltext/`

A **job** (no HTTP server): it opens PostgreSQL through `gofi`, applies the migrations (the catalog, then its full-text index), runs a few full-text searches, prints the results and exits.

- `.gofi/gofi-sdk-go/examples/sqln/fulltext/main.go`

## sqln — search job

`.gofi/gofi-sdk-go/examples/sqln/search/`

A **job** (no HTTP server): it opens PostgreSQL through `gofi`, applies the migrations, runs a few `sqln` criteria queries, prints the results and exits.

- `.gofi/gofi-sdk-go/examples/sqln/search/main.go`

