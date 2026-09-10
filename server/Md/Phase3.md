# Phase 3: Payment Integration (Distributed State & Security)

**Duration:** 2–3 weeks

Payments are not just API calls — they are distributed state machines operating over untrusted networks.

## Core Principle

Moving money requires two-phase commit logic and cryptographic verification. You can never trust the client, and you can't fully trust the network.

## Prerequisites

- **HTTP Protocol**: Headers, status codes (200 vs 201 vs 202), POST bodies.
- **Cryptography Basics**: Hashing (SHA-256) and HMAC (Hash-based Message Authentication Code).
- **Webhooks**: A webhook is just the payment provider sending an HTTP POST request to your server.

## Learning Progression

1. **State Machine**: Draw your order states on paper. `Draft → Pending Payment → Paid → Processing → Shipped`. Never allow an order to jump from Draft to Shipped.

2. **Manual Webhook Handling**: When Stripe/Razorpay sends a webhook to your Gin endpoint, do not use their SDK to verify it. Write the raw Go code using `crypto/hmac` to verify the signature in the HTTP header yourself. This teaches you how security actually works.

3. **Idempotency Keys**: Generate a unique ID for every payment attempt. If the network drops and the frontend retries, the payment gateway uses that ID to know it's a duplicate and doesn't charge the card twice.

## How to Integrate

1. Look at the `/webhook/payment` endpoint plan in the README. Do **not** use the Stripe/Razorpay SDK to verify the webhook.
2. **Exercise**: Write raw Go code using `crypto/hmac` and `sha256` to verify the signature in the HTTP header yourself.
3. Add an `idempotency_key` (string, unique) column to the orders table. Before processing a webhook, check if an order with that key already exists. If yes, return 200 OK immediately without processing again.
