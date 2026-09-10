Phase 3: Payment Integration (Distributed State & Security)
Payments are not just API calls; they are Distributed State Machines operating over untrusted networks.

    The Abstraction: "Just use the Stripe SDK and let it handle everything."
    First Principle: Moving money requires Two-Phase Commit logic and Cryptographic Verification. You can never trust the client (frontend), and you can't even fully trust the network.
    Prerequisites to learn first:
        HTTP Protocol Deep Dive: Understand Headers, Status Codes (200 vs 201 vs 202), and POST bodies.
        Cryptography Basics: Understand Hashing (SHA256) and HMAC (Hash-based Message Authentication Code).
        Webhooks: Understand that a webhook is just the payment provider sending an HTTP POST request to your server.
    Your Learning Progression:
        The State Machine: Draw your order states on paper. Draft -> Pending Payment -> Paid -> Processing -> Shipped. Never allow an order to jump from Draft to Shipped.
        Manual Webhook Handling: When Stripe/Razorpay sends a webhook to your Gin endpoint, do not use their SDK to verify it. Write the raw Go code using crypto/hmac to verify the signature in the HTTP header yourself. This teaches you how security actually works.
        Idempotency Keys: Learn how to generate a unique ID for every payment attempt. If the network drops and the frontend retries, the payment gateway uses that ID to know it's a duplicate and doesn't charge the card twice.
