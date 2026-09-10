Phase 5: Caching & Database (The Glue)
Since you are using GORM, the abstraction is very thick.

    The Abstraction: "GORM handles the database, Redis handles the cache."
    First Principle: The database is the Source of Truth (Slow, consistent, durable). The cache is a Lie (Fast, eventually consistent, volatile).
    Prerequisites to learn first:
        SQL Transactions: BEGIN, COMMIT, ROLLBACK. Understand ACID properties.
        Cache Invalidation Strategies: Write-Through, Write-Back, and Cache-Aside.
    Your Learning Exercise:
    Implement the Cache-Aside pattern manually.
        User requests Product.
        Check Redis. If hit, return.
        If miss, query Postgres via GORM.
        Save to Redis with a TTL (Time To Live).
        Return to user.
        Then, figure out what happens when the Admin updates the product price. How do you invalidate the Redis cache so the user doesn't see the old price?
