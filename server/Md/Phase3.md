//STARTED

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

<!-- Details -->

🟢 Phase 1: Initiation (Frontend → Go Backend)
Step 1: User Clicks "Pay"

    Next.js: User reviews their cart and clicks "Pay with eSewa". Next.js sends a POST /api/orders request to your Go backend with the cart details.

Step 2: Order Creation & Idempotency Key Generation

    Go Backend:
        Generates a unique transaction_uuid (e.g., ord_9f8a7b6c). This is your Idempotency Key.
        Saves a new record in the database: { uuid: "ord_9f8a7b6c", amount: 1000, status: "PENDING_PAYMENT" }.
        Why? If the network drops and the user clicks "Pay" again, the frontend might retry. The backend will see the same intent and handle it gracefully, or the transaction_uuid ensures eSewa doesn't process duplicates.

Step 3: Outbound Signature Generation

Go Backend:
Builds the payload string required by eSewa: total_amount=1000,transaction_uuid=ord_9f8a7b6c,product_code=EPAYTEST...
Uses crypto/hmac and crypto/sha256 with your shared secret key to generate the request signature.
Returns a JSON response to Next.js containing the eSewa Payment URL and all the signed form data.

🟡 Phase 2: The Payment (Next.js → eSewa)
Step 4: Browser Redirect

    Next.js: Receives the data from Go and immediately redirects the user's browser (window.location.href) to the eSewa payment page, appending the signed payload as a POST form or URL parameters.

Step 5: User Authenticates and Pays

    eSewa: Displays the payment page. The user enters their eSewa credentials (or uses the sandbox test credentials). eSewa deducts the money and marks the transaction as COMPLETE on their end.

🔴 Phase 3: Verification & State Transition (eSewa → Go Backend → Next.js)
Step 6: The Redirect Back

    eSewa: Redirects the user's browser back to your success_url.
    Crucial Detail: Your success_url should point to your Go Backend (e.g., https://your-api.com/payment/esewa/success?data=BASE64_STRING), not directly to your Next.js frontend. This ensures your backend verifies the payment before the user ever sees a "Success" screen.

Step 7: Raw HMAC Verification (The Core Exercise)

    Go Backend:
        Extracts the data query parameter and Base64-decodes it into a JSON object.
        Extracts the signed_field_names (e.g., "transaction_code,status,total_amount...").
        Reconstructs the exact key=value,key=value string from the decoded JSON.
        Runs hmac.New(sha256.New, []byte(secret)) on that string.
        Base64-encodes the result and compares it to the signature in the JSON using hmac.Equal.
        Why? To mathematically prove the redirect came from eSewa and wasn't forged by a user manually typing a URL.

Step 8: Idempotency & State Machine Check

    Go Backend: Queries the database for transaction_uuid.
        Scenario A: Status is already PAID. → Action: Abort processing. Return a redirect to the frontend success page. (Prevents double-crediting the user if eSewa redirected twice).
        Scenario B: Status is PENDING_PAYMENT. → Action: Proceed to Step 9.
        Scenario C: Status is FAILED or doesn't exist. → Action: Reject the request.

Step 9: Server-to-Server Verification (Distributed Trust)

    Go Backend: Because the redirect in Step 6 happened via the user's browser (which is an untrusted network), your Go backend makes a direct, background HTTP GET request to eSewa's status API:
    GET https://rc.esewa.com.np/api/epay/transaction/status/?product_code=...&total_amount=...&transaction_uuid=...
    eSewa: Responds with { "status": "COMPLETE" }.
    Why? This is the ultimate source of truth. Even if a hacker somehow spoofed the browser redirect, they cannot fake this server-to-server check.

Step 10: State Machine Transition

    Go Backend:
        Updates the database: UPDATE orders SET status = 'PAID', esewa_transaction_code = 'XYZ123' WHERE transaction_uuid = 'ord_9f8a7b6c'.
        Triggers any post-payment logic (e.g., sending an email, unlocking a digital product).

Step 11: Final Redirect to Frontend

    Go Backend: Responds to the browser with an HTTP 302 Redirect to your Next.js success page: https://your-frontend.com/order/success?uuid=ord_9f8a9f8a7b6c.

🟢 Phase 4: Confirmation (Next.js)
Step 12: Frontend Fetches Final State

    Next.js: The success page loads. It immediately calls your Go backend: GET /api/orders/ord_9f8a7b6c.
    Go Backend: Returns { "status": "PAID", "amount": 1000 }.
    Next.js: Displays the "Payment Successful!" receipt to the user.

🛡️ Summary of Security Guarantees in this Workflow:

    Client can't fake payment: The frontend never handles the secret key or the final state change. Only the Go backend does.
    Network can't replay payment: The transaction_uuid (idempotency key) ensures that even if the exact same webhook/redirect is received 10 times, the database state machine only allows the PENDING → PAID transition once.
    Hacker can't spoof the redirect: The raw hmac.Equal check guarantees the payload wasn't tampered with in transit.
    eSewa can't lie (and neither can the browser): The server-to-server GET request in Step 9 is the cryptographic handshake that seals the deal.

This is exactly how professional, high-integrity payment systems are built.
