# Phase 2: Queue Management (Async Processing)

**Duration:** 2–3 weeks

Decouple "Add to Cart" or "Send Email" logic from synchronous request handling.

## Core Principle

A message queue is a durable, distributed linked list (or log). Its two jobs:

- **Backpressure**: If consumers are slow, messages pile up safely instead of crashing the server.
- **Decoupling**: The producer doesn't need to know the consumer exists.

## Prerequisites

- **TCP/IP & Sockets**: A queue is just a server listening on a TCP port.
- **Blocking vs Non-blocking I/O**: How a program waits for data without freezing.
- **Idempotency**: An operation can be applied multiple times without changing the result beyond the initial application. Crucial for queues.

## Learning Progression

| Level | Stack | What You Learn |
|-------|-------|----------------|
| 1 | Go channels + goroutine worker | The pattern (in-memory) |
| 2 | Redis Lists (`LPUSH` / `BRPOP`) | Persistence and polling |
| 3 | RabbitMQ | Exchanges, routing keys, ACKs |

Skip Kafka for now — it's too complex for a solo learner and abstracts away the basic queue mechanics.

**The "Aha!" Moment**: Intentionally crash your consumer while processing a message. Learn how RabbitMQ's NACK (Negative Acknowledgment) re-queues the message so it isn't lost.

## How to Integrate

| Week | Task |
|------|------|
| 1 | Create `pkg/worker`. Build a worker pool using Go channels that processes a slice of "email tasks". |
| 2 | Replace the channel with Redis `LPUSH` (producer) and `BRPOP` (consumer). |
| 3 | Install RabbitMQ via Docker. When `PlaceOrder` succeeds, publish an `OrderCreated` event to a RabbitMQ exchange instead of calling an email function synchronously. |

**Exercise**: Intentionally `panic` in your email worker. Use RabbitMQ NACK and Dead Letter Queues (DLQ) so the message isn't lost and can be retried.
