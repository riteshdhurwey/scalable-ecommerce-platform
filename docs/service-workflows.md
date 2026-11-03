# Service Workflows

This document describes the responsibility and main workflow of each service.

## 1. User Service

**Owns:** user identity, password hashes, profile data.

### Registration
1. Validate required fields.
2. Normalize and validate the email or username.
3. Hash the password with a suitable password-hashing algorithm.
4. Enforce uniqueness through a database constraint.
5. Persist the account.
6. Return a safe response without password or hash fields.

### Authentication
1. Validate credentials.
2. Apply rate limiting and generic failure responses.
3. Issue a session or short-lived access token according to the chosen authentication design.
4. Validate authorization on protected endpoints.

Do not store plaintext passwords. Avoid logging credentials or tokens.

## 2. Catalog & Inventory Service

**Owns:** product records, categories, prices, available stock, reservations.

### Product browsing
1. Receive a product query.
2. Read product details from the database or cache.
3. Return product information and an availability indication.
4. Invalidate or refresh relevant cache entries when product data changes.

### Stock reservation
1. Receive order ID and requested product quantities.
2. Check and reserve all required quantities atomically.
3. If any item cannot be reserved, roll back the reservation transaction.
4. Persist a reservation ID and expiry time.
5. Return the reservation outcome and authoritative price snapshot.

For multi-item orders, ensure a failed reservation does not leave only some items reserved unless partial reservation is explicitly supported.

## 3. Cart Service

**Owns:** active cart contents and quantities.

### Add/update/remove item
1. Authenticate the user.
2. Validate product ID and quantity.
3. Optionally check that the product exists.
4. Add, update, or remove the cart item.
5. Return the current cart.

The cart is not the source of truth for inventory or checkout prices. Revalidate both during checkout. Redis can store active carts, but define expiry and persistence expectations before relying on it as the only copy.

## 4. Order Service

**Owns:** order status, order items, quantity, price snapshots, totals, order history.

### Create order
1. Validate the idempotency key.
2. Retrieve cart contents.
3. Validate products and prices with Catalog/Inventory.
4. Request inventory reservation.
5. Persist the order and its state transition.
6. Initiate payment or publish the next workflow event.
7. Return the order ID and current status.

### Read order
1. Authenticate the caller.
2. Confirm the caller owns the order or has administrative permission.
3. Return order details and current status.

Do not expose another user's order merely because its ID is known.

## 5. Payment Service

**Owns:** payment attempts, provider references, payment status, webhook event IDs.

### Payment initiation
1. Validate the order reference and amount against trusted order data.
2. Create a payment attempt with a stable idempotency key.
3. Call the provider with a timeout.
4. Persist the provider reference and result.
5. Handle pending outcomes through webhook/reconciliation workflows.

### Webhook
1. Verify the provider signature.
2. Validate the event and provider reference.
3. Deduplicate the provider event ID.
4. Apply a valid payment-state transition.
5. Publish a payment outcome using the outbox pattern.

Never assume the browser redirect is proof of payment success.

## 6. Notification Service

**Owns:** notification delivery attempts and statuses.

### Send notification
1. Consume an event such as `OrderConfirmed`.
2. Check whether the event/notification has already been processed.
3. Render the appropriate message from trusted data.
4. Send through an email or SMS provider.
5. Record success or failure.
6. Retry transient failures with bounded backoff.
7. Move repeatedly failing messages to a dead-letter queue for inspection.

Notification failure must not reverse a successful payment or confirmed order.

## 7. API Gateway

1. Accept external requests.
2. Apply request-size and rate limits.
3. Route requests to the appropriate service.
4. Forward request and trace IDs.
5. Apply TLS configuration at the appropriate edge.
6. Return consistent errors where practical.

Authentication and authorization must still be enforced by trusted services; gateway checks alone should not be the only protection for sensitive operations.
