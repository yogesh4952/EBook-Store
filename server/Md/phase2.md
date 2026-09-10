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
