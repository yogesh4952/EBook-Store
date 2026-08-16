# EBook Store

Online book store (Book Store Nepal). A monorepo with two parts:

- **`client/`** — Next.js 16 (React 19) frontend with a BFF (Backend-for-Frontend) API layer.
- **`server/`** — Go REST API backend (replaces the former .NET server).

## Project layout

```
EBookStore/
├── client/            # Next.js frontend (see client/README.md)
├── server/            # Go REST API backend (see server/README.md)
└── README.md          # this file
```

## Architecture

```
Browser ──▶ Next.js BFF (/api/*) ──▶ Go REST API (server/) ──▶ PostgreSQL
               │
               └─ static pages / assets
```

The frontend never talks to the Go API directly from the browser. All
backend calls go through Next.js route handlers (`client/app/api/*`),
which proxy to the Go server and manage the auth cookie. This keeps the
JWT token in an `httpOnly` cookie and keeps secrets off the client.

## Frontend configuration (static)

The frontend no longer reads any environment variables. All config lives
in one file: **`client/lib/config.ts`**.

| Value              | Default                  | Purpose                          |
| ------------------ | ------------------------ | -------------------------------- |
| `BACKEND_URL`      | `http://localhost:8080`  | Base URL of the Go API           |
| `AUTH_COOKIE_NAME` | `accessToken`            | Cookie that stores the JWT       |
| `API_ENDPOINTS`    | derived from `BACKEND_URL` | Backend endpoint paths          |

To point the app at a different backend, edit `client/lib/config.ts`
only — no `.env` file is needed.

## Backend configuration

The Go API is configured via environment variables (or a `.env` file).
Copy `server/.env.example` to `server/.env` and adjust values:

| Variable            | Default                                             | Purpose              |
| ------------------- | --------------------------------------------------- | -------------------- |
| `SERVER_HOST`       | `0.0.0.0`                                           | Bind host            |
| `SERVER_PORT`       | `8080`                                              | Bind port            |
| `DATABASE_URL`      | `postgres://yogesh:yogesh@localhost:5432/ebookstore`| PostgreSQL DSN       |
| `JWT_SECRET`        | `your-super-secret-key-...` (32+ chars)             | JWT signing key      |
| `JWT_ISSUER`        | `EBookStore`                                        | JWT issuer claim     |
| `JWT_AUDIENCE`      | `EBookStoreUsers`                                   | JWT audience claim   |
| `JWT_EXPIRY_MINUTES`| `60`                                                | Token lifetime       |

## API contract

Routes mirror the old .NET API so the frontend needed no changes.

| Method | Path                 | Auth | Description            |
| ------ | -------------------- | ---- | ---------------------- |
| POST   | `/api/Auth/register` | —    | Create a user          |
| POST   | `/api/Auth/login`    | —    | Login, returns JWT     |
| GET    | `/api/Auth/me`       | Bearer token | Return current user |
| GET    | `/healthz`           | —    | Liveness check         |

### Login request / response

```http
POST /api/Auth/login
Content-Type: application/json

{ "email": "reader@example.com", "password": "secret" }
```

```json
200 OK
{ "token": "<jwt>", "email": "reader@example.com", "firstName": "John" }
```

### Register request

```http
POST /api/Auth/register
Content-Type: application/json

{ "firstName": "John", "lastName": "Doe",
  "email": "reader@example.com", "password": "secret", "role": 0 }
```

`role`: `0` = User, `1` = Vendor, `2` = Admin (default `0`).

## Quick start

### 1. Start the database (PostgreSQL via Docker)

```bash
cd server
docker compose up -d db     # applies ./migrations on first start
```

### 2. Run the Go API

```bash
cd server
cp .env.example .env        # optional, defaults already work
make run                    # or: go run ./cmd/api
# → listening on http://localhost:8080
```

### 3. Run the frontend

```bash
cd client
npm install
npm run dev                 # → http://localhost:3000
```

## What changed in this migration

See **`server/README.md`** for a full write-up of the .NET → Go migration,
the frontend static-config change, and a guide to how each Go folder is
used during development.