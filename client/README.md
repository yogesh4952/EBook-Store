# EBook Store — Frontend (client)

Next.js 16 (React 19) frontend for the EBook Store. The app uses a
Backend-for-Frontend (BFF) pattern: the browser only talks to Next.js
route handlers, which proxy to the Go API and manage the auth cookie.

## Getting started

```bash
npm install
npm run dev          # → http://localhost:3000
```

The frontend is fully **static** — it needs no `.env` and no backend to
start. All configuration lives in **`lib/config.ts`**:

```ts
export const BACKEND_URL = "http://localhost:8080"; // point at the Go API
```

Edit that single file to change where the app talks to.

## Structure

```
app/
├── (public)/          # public pages (uses Navbar layout)
├── auth/login/        # login page
├── api/               # BFF route handlers
│   ├── login/         # POST → proxies to Go /api/Auth/login, sets httpOnly cookie
│   ├── me/            # GET  → proxies to Go /api/Auth/me with the JWT
│   └── logout/        # POST → clears the auth cookie
components/common/     # shared UI (Navbar)
helper/                # small utilities (JWT expiry parsing)
lib/config.ts          # ★ static config — edit values here
store/                 # Zustand state (auth)
middleware.ts          # route protection based on the accessToken cookie
```

## Auth flow

1. Browser POSTs credentials to `/api/login`.
2. The BFF calls the Go API, gets a JWT, and stores it in an `httpOnly`
   cookie named `accessToken`.
3. `/api/me` sends that cookie's token as `Authorization: Bearer` to the
   Go API and returns the user.
4. `middleware.ts` redirects unauthenticated visitors to `/auth/login`.

## Build & lint

```bash
npm run build
npm run lint
```

## Docker

A minimal `Dockerfile` is included for static deployment (`next start`).