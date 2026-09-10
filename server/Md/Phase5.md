# Phase 5: Caching & Database

**Duration:** 2 weeks

The database is the source of truth. The cache is a lie for the sake of speed. You must manage the lifecycle of that lie.

## Core Principle

- **Database (PostgreSQL)**: Slow, consistent, durable — the source of truth.
- **Cache (Redis)**: Fast, eventually consistent, volatile — a performance optimization.

## Prerequisites

- **SQL Transactions**: `BEGIN`, `COMMIT`, `ROLLBACK`. Understand ACID properties.
- **Cache Invalidation Strategies**: Write-Through, Write-Back, and Cache-Aside.

## Learning Exercise: Cache-Aside Pattern

Implement it manually:

1. User requests a product.
2. Check Redis. If hit → return.
3. If miss → query PostgreSQL via GORM.
4. Save result to Redis with a TTL (Time To Live).
5. Return to user.

Then figure out: **what happens when the admin updates the product price?** How do you invalidate the Redis cache so the user doesn't see the old price?

## How to Integrate

1. Add Redis to your `pkg/database` setup.

2. **Cache-Aside** — in `bookRepo.FindById`:
   ```go
   // 1. Check Redis
   cached, err := r.redis.Get(ctx, "book:"+id).Result()
   if err == nil {
       // unmarshal and return
   }

   // 2. Cache miss — query GORM
   var book Book
   r.db.First(&book, id)

   // 3. Save to Redis with 10-min TTL
   data, _ := json.Marshal(book)
   r.redis.Set(ctx, "book:"+id, data, 10*time.Minute)

   return book, nil
   ```

3. **Invalidation** — in `bookRepo.UpdateBook`, after the GORM update succeeds:
   ```go
   r.redis.Del(ctx, "book:"+id)
   ```
