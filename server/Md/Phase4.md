# Phase 4: Benchmarking & Profiling (Performance)

**Duration:** 1–2 weeks

Know if your Go code is fast — and where it isn't.

## Core Principle

You cannot optimize what you cannot measure. Latency (time per request) and throughput (requests per second) are different metrics.

## Prerequisites

- **Linux Basics**: CPU cores, RAM, file descriptors.
- **Latency vs Throughput**: Latency is how long one request takes. Throughput is how many requests per second. They are inversely related.
- **Database Indexing**: You cannot benchmark an app if your SQL queries are doing full table scans.

## Learning Progression

1. **Go Native Profiling (`pprof`)**: Built into Go. Add `net/http/pprof` to your Gin router. Learn how to generate a CPU flame graph and a memory heap profile. See exactly which line of Go code is allocating memory.

2. **Load Testing**: Use [k6](https://k6.io) (written in Go, so you can read its source code!) or `wrk`. Write a script that hammers your `/api/cart` endpoint with 1,000 concurrent users.

3. **The RED Method**: Measure your API by **Rate** (requests/sec), **Errors** (percentage failing), and **Duration** (latency).

## How to Integrate

1. Add `import _ "net/http/pprof"` to `cmd/api/main.go`. Expose it on a separate internal port (e.g., `:8081`).

2. Install k6:
   ```bash
   go install go.k6.io/k6@latest
   ```

3. Write a script to hit `/api/books` with 100, 500, and 1000 concurrent virtual users.

4. **Exercise**: While k6 is running, generate a CPU profile:
   ```bash
   go tool pprof http://localhost:8081/debug/pprof/profile?seconds=10
   ```
   Open it in the web UI. Find the exact line of Go code causing the most allocations or CPU time, and fix it (e.g., fixing an N+1 query in GORM).
