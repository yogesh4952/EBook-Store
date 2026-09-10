# Phase 1: Go Concurrency & Context

**Duration:** 1–2 weeks

Before touching queues, payments, or benchmarking, you must understand how Go handles concurrent work.

## Core Principle

Goroutines are cheap, but uncontrolled goroutines cause memory leaks. `context.Context` is the only safe way to cancel them.

## Prerequisites

- **Channels**: Buffered vs unbuffered. Channels communicate state; `sync.Mutex` protects state.
- **`context.Context`**: Mandatory for passing deadlines, cancellation signals, and request-scoped values across API boundaries and goroutines.

## How to Integrate

1. **Audit every service and repository function** — ensure `ctx context.Context` is the first argument.
2. In `cmd/api/main.go`, create a root context with a timeout for incoming HTTP requests. Extract it in Gin handlers with `ctx := c.Request.Context()`.
3. **Exercise**: In `GetBookDetails`, if you fetch the Book and Author separately, use `golang.org/x/sync/errgroup` to fetch them concurrently instead of sequentially.

## Learning Exercise

Build a background worker that processes a slice of numbers:

1. Use a channel to feed it work.
2. Use `context.WithTimeout` to kill the worker if it takes too long.
3. No external worker libraries — channels and goroutines only.
