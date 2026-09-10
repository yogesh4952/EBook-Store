Phase 4: Benchmarking & Profiling (Performance)
You want to know if your Go code is fast.

    The Abstraction: "Just use JMeter or look at the server CPU."
    First Principle: Performance is about identifying bottlenecks. Is your system CPU-bound (doing heavy math/encryption) or I/O-bound (waiting for the database or network)?
    Prerequisites to learn first:
        Linux Basics: Understand CPU cores, RAM, and File Descriptors.
        Latency vs. Throughput: Latency is how long one request takes. Throughput is how many requests per second. They are inversely related.
        Database Indexing: You cannot benchmark an app if your SQL queries are doing full table scans.
    Your Learning Progression:
        Go Native Profiling (pprof): This is built into Go. Add net/http/pprof to your Gin router. Learn how to generate a CPU Flame Graph and a Memory Heap Profile. See exactly which line of Go code is allocating memory.
        Load Testing: Use k6 (it's written in Go, so you can read its source code!) or wrk. Write a script that hammers your /api/cart endpoint with 1,000 concurrent users.
        The RED Method: Learn to measure your API by Rate (requests per sec), Errors (percentage failing), and Duration (latency).
