Phase 1: Go Concurrency & Context (The Prerequisites for Everything)
Before you touch queues, payments, or benchmarking, you must deeply understand how Go handles doing multiple things at once.

    The Abstraction: "Just use go func() and sync.WaitGroup."
    First Principle: Goroutines are not OS threads; they are user-space threads multiplexed onto OS threads by the Go Scheduler.
    Prerequisites to learn first:
        1. Channels: Understand buffered vs. unbuffered channels. Understand that channels are for communicating state, while Mutexes (sync.Mutex) are for protecting state.
        2. context.Context: This is mandatory. You cannot build payment gateways or queues without it. Context is how you pass deadlines, cancellation signals, and request-scoped values across API boundaries and goroutines.
    Your Learning Exercise: Build a simple background worker in Go that processes a slice of numbers. Use a channel to feed it work. Then, use context.WithTimeout to kill the worker if it takes too long. Do not use any external worker libraries yet.
