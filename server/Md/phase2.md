Phase 2: Queue Management (Asynchronous Processing)
You want to decouple your "Add to Cart" or "Send Email" logic.

    The Abstraction: "Just use RabbitMQ, Kafka, or AWS SQS."
    First Principle: A message queue is just a durable, distributed linked list (or log). Its core job is to provide Backpressure (if consumers are slow, messages pile up safely instead of crashing the server) and Decoupling (Producer doesn't need to know the Consumer exists).
    Prerequisites to learn first:
        TCP/IP & Sockets: Understand that a queue is just a server listening on a TCP port.
        Blocking vs. Non-blocking I/O: Understand how a program waits for data without freezing.
        Idempotency: The concept that an operation can be applied multiple times without changing the result beyond the initial application. (Crucial for queues).
    Your Learning Progression:
        Level 1 (In-Memory): Build a queue using Go channels and a background Goroutine worker. (Teaches the pattern).
        Level 2 (Persistent): Use Redis Lists (LPUSH to add, BRPOP to consume). Redis blocks the connection until a message arrives. (Teaches persistence and polling).
        Level 3 (The Real Deal): Introduce RabbitMQ (Skip Kafka for now; Kafka is too complex for a solo learner and abstracts away the basic queue mechanics). Learn about Exchanges, Routing Keys, and Acknowledgments (ACKs).
    The "Aha!" Moment: Intentionally crash your consumer service while it's processing a message. Learn how RabbitMQ's "Negative Acknowledgment" (NACK) re-queues the message so it isn't lost.

Phase 2: Queue Management (Async Processing)

    Expected Time: 2–3 Weeks
    First Principle to Learn: A queue is just a durable list that provides backpressure and decoupling.
    How to Integrate:
        Week 1 (In-Memory): Create pkg/worker. Build a simple worker pool using Go channels that processes a slice of "email tasks".
        Week 2 (Redis): Replace the channel with Redis LPUSH (producer) and BRPOP (consumer).
        Week 3 (RabbitMQ): Install RabbitMQ via Docker. When PlaceOrder succeeds, instead of calling an email function synchronously, publish an OrderCreated event to a RabbitMQ exchange.
        Exercise: Intentionally panic in your email worker. Learn how to use RabbitMQ NACK (Negative Acknowledgment) and Dead Letter Queues (DLQ) so the message isn't lost and can be retried.
