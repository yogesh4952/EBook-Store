# EBook Store

A full-stack ebook marketplace built with **Go** and **Next.js**. Sellers list books, customers browse and place orders with COD or eSewa payment options.

## Architecture

```
Browser → Next.js BFF (port 3000) → Go REST API (port 8080) → PostgreSQL + Redis
                │
                └── server-side rendering, static assets
```

- The browser never calls the Go API directly. All backend requests are proxied through Next.js route handlers (`app/api/*`), which manage the JWT in an `httpOnly` cookie.
- **Server:** Go 1.26, Gin, GORM, PostgreSQL 16, Redis
- **Client:** Next.js 16 (App Router), React 19, Tailwind CSS v4, Zustand

## Tech Stack

| Layer    | Technology                                      |
|----------|-------------------------------------------------|
| API      | Go, Gin, GORM                                   |
| Auth     | OTP via email (Redis-backed), JWT (HS256)       |
| Database | PostgreSQL 16, Redis (OTP storage)              |
| Frontend | Next.js 16, React 19, TypeScript, Tailwind v4   |
| State    | Zustand                                         |
| Logging  | Zerolog with request-ID correlation             |
| Docs     | Swagger (swaggo)                                |

## Project Structure

```
EBook-Store/
├── client/                         # Next.js frontend
│   ├── app/
│   │   ├── (public)/               # Public pages (homepage, catalog)
│   │   ├── auth/                   # Login (OTP-based), register
│   │   └── api/                    # BFF route handlers (proxy to Go)
│   ├── components/                 # Book cards, navbar, OTP input
│   ├── lib/config.ts               # Backend URL, endpoint map
│   ├── store/                      # Zustand auth store
│   └── proxy.ts                    # Route protection middleware
│
├── server/                         # Go REST API
│   ├── cmd/api/main.go             # Entry point, route registration, DI
│   ├── internal/
│   │   ├── auth/                   # OTP login, register, JWT issuance
│   │   ├── book/                   # Book CRUD (seller-owned)
│   │   ├── order/                  # Order placement, listing
│   │   ├── address/                # User address management
│   │   ├── user/                   # User model and listing
│   │   ├── middleware/             # JWT auth, role-based access, logging
│   │   ├── initializers/           # DB, Redis, env bootstrapping
│   │   ├── db/                     # Migration and seed scripts
│   │   ├── cart/                   # (stub — not yet implemented)
│   │   └── sellers/                # (stub — not yet implemented)
│   ├── pkg/
│   │   ├── utils/                  # JWT, OTP generation, pagination
│   │   └── logger/                 # Zerolog + GORM adapter
│   ├── docs/                       # Generated Swagger files
│   └── data/book.json              # Seed data (~50 books)
│
├── docker-compose.yml              # PostgreSQL 16
└── README.md
```

## Prerequisites

- Go 1.26+
- Node.js 20+
- PostgreSQL 16
- Redis 7+
- Docker (optional, for database)

## Quick Start

### 1. Start PostgreSQL

```bash
docker compose up -d
```

This starts PostgreSQL 16 on port 5432 with credentials `yogesh/yogesh` and database `ebookstore`.

### 2. Configure the server

```bash
cd server
cp .env.example .env   # if available, or create manually
```

Required environment variables:

| Variable       | Default                                              | Description         |
|----------------|------------------------------------------------------|---------------------|
| `PORT`         | `8080`                                               | Server bind port    |
| `DSN`          | `postgres://yogesh@localhost:5432/ebookstore`        | PostgreSQL DSN      |
| `jwt_secret`   | —                                                    | JWT signing key     |
| `SMTP_EMAIL`   | —                                                    | SMTP sender email   |
| `SMTP_PASSWORD`| —                                                    | SMTP app password   |
| `SMTP_HOST`    | `smtp.gmail.com`                                     | SMTP host           |
| `SMTP_PORT`    | `587`                                                | SMTP port           |
| `REDIS_HOST`   | `localhost`                                          | Redis host          |
| `REDIS_PORT`   | `6379`                                               | Redis port          |

### 3. Run the server

```bash
cd server
go run ./cmd/api
# → http://localhost:8080
# → Swagger UI at http://localhost:8080/swagger/index.html
```

For development with hot-reload:

```bash
air
```

### 4. Seed the database

The database is auto-migrated on startup. To seed sample book data:

```bash
go run ./internal/db/seed
```

### 5. Run the client

```bash
cd client
npm install
npm run dev
# → http://localhost:3000
```

The client points to `http://localhost:8080` by default. To change, edit `client/lib/config.ts`.

## API Endpoints

### Auth

| Method | Path                  | Auth | Description                          |
|--------|-----------------------|------|--------------------------------------|
| POST   | `/api/auth/send-otp`  | —    | Send 6-digit OTP to email            |
| POST   | `/api/auth/login`     | —    | Verify OTP, return JWT               |
| POST   | `/api/auth/register`  | —    | Register new user (auto-creates seller profile for seller role) |

### Books

| Method | Path                   | Auth                   | Description                    |
|--------|------------------------|------------------------|--------------------------------|
| GET    | `/api/book/list-books` | —                      | List books (paginated)         |
| POST   | `/api/book/publish-book` | Seller/Admin         | Publish a new book             |
| PATCH  | `/api/book/update-book`  | Seller/Admin (owner) | Update book (atomic, owner-only) |

### Orders

| Method | Path                    | Auth   | Description              |
|--------|-------------------------|--------|--------------------------|
| POST   | `/api/order/place-order` | Bearer | Place an order           |
| GET    | `/api/order/list-user-order` | Bearer | List orders for current user |

### Addresses

| Method | Path                     | Auth   | Description            |
|--------|--------------------------|--------|------------------------|
| POST   | `/api/address/add-address` | Bearer | Add a delivery address |

### Other

| Method | Path           | Description       |
|--------|----------------|-------------------|
| GET    | `/swagger/*`   | Swagger UI        |

## Authentication Flow

```
1. Client → POST /api/auth/send-otp { email }
2. Server generates 6-digit OTP, stores in Redis (5-min TTL), emails it
3. Client → POST /api/auth/login { email, otp }
4. Server verifies OTP against Redis, issues JWT (1h expiry)
5. BFF stores JWT in httpOnly cookie (accessToken)
6. Subsequent requests: cookie → BFF reads token → forwards as Authorization: Bearer → Go API
```

JWT claims: `user_id`, `email`, `role` (admin/seller/customer).

### Role-Based Access

- `AuthRequired` middleware validates the JWT and injects user context.
- `AuthorizeRoles("seller", "admin")` restricts routes by role.

## Database Schema

```
users ──< orders ──< order_items
  │
  ├──< books (via sellers)
  │
  ├──< addresses
  │
  └──< sellers
```

Key models:

- **User**: email, phone, role (admin/seller/customer)
- **Seller**: linked to User, unique seller number
- **Book**: title, author, genre, price, stock status, linked to Seller
- **Order**: unique order code, payment method/status, shipping address, total
- **OrderItem**: book + quantity + unit price, composite key on (order_id, book_id)
- **UserAddress**: city + delivery address, linked to User

## Development

### Run tests

```bash
cd server
go test ./...
```

### Generate Swagger docs

```bash
cd server
swag init -g cmd/api/main.go
```

### Lint

```bash
cd client
npm run lint
```

## Current Limitations

- Cart functionality is stubbed (not implemented)
- Seller management endpoints are stubbed
- Payment integration (eSewa) is modeled but not wired to a gateway
- `/api/auth/me` endpoint is referenced by the client but not implemented in the server
- Social login (Google, Facebook) buttons are UI placeholders
- Registration page is an empty placeholder
