# Failure Recovery and Reliability

## 1. Design objective

A temporary network failure, process restart, duplicate event, or timeout should not silently corrupt business state. Recovery should be observable and repeatable.

## 2. Retry policy

Retry only errors likely to be transient, such as selected network failures, rate limits, or server errors.

- Use bounded exponential backoff with jitter.
- Set maximum attempts or a retry deadline.
- Do not retry invalid input or permanent business failures.
- Make the operation idempotent before enabling automatic retries.
- Record attempt count and last error.
- Route exhausted messages to a dead-letter queue or recovery process.

## 3. Idempotency

Use idempotency keys for operations that may be repeated, especially checkout and payment initiation.

A typical request lifecycle:
1. Receive the key and request payload.
2. Look up an existing result for that key.
3. If a matching completed or pending operation exists, return its result.
4. If the same key is used with a different payload, reject it.
5. Otherwise, persist the key and operation safely with database uniqueness constraints.

Idempotency retention must cover the expected retry window. It prevents duplicate logical operations; it does not replace database constraints or provider-side payment idempotency.

## 4. Transactional outbox

Without an outbox, a service can commit a database change and crash before publishing its event.

Recommended pattern:
1. In one local database transaction, update business state and insert an outbox record.
2. A publisher reads unpublished outbox records.
3. Publish each event to the broker.
4. Mark the record as published after broker acknowledgement.
5. Retry records after failures.
6. Make consumers idempotent because a publish acknowledgement can be lost.

## 5. Consumer deduplication

Each event should have a unique event ID.

A consumer should atomically record the processed event ID alongside its local business change when possible. If the event ID already exists, it should avoid applying the same logical change again.

Do not assume a broker provides exactly-once business processing across a database and external provider.

## 6. Crash recovery

| Crash point | Recovery |
|---|---|
| Before database commit | Transaction rolls back; client or worker can retry |
| After commit, before publish | Outbox publisher finds the pending event |
| After publish, before marking outbox row | Event may publish again; consumer deduplicates |
| After consuming, before local commit | Broker redelivery causes a retry |
| After local commit, before acknowledging message | Redelivery is safe through deduplication |
| During payment request | Reconcile using the provider's idempotency key/reference; outcome may be unknown |
| During inventory reservation | Reservation state and expiry worker determine whether to commit or release |

## 7. Timeouts and circuit breaking

- Use explicit request and database timeouts.
- Avoid unbounded goroutines or unbounded work queues.
- Limit concurrent requests to downstream dependencies.
- Consider circuit breaking only when its behavior and recovery policy are understood.
- Do not treat timeout as proof that a remote operation did not happen.

## 8. Inventory and payment compensation

Examples of compensating actions:
- Release an inventory reservation when payment definitively fails.
- Cancel an order if inventory reservation fails before payment begins.
- Reconcile an unknown payment outcome before releasing stock or initiating another charge.
- Escalate inconsistent states to a retryable reconciliation worker or manual review.

Compensation is a new business operation, not a rollback of a distributed transaction.

## 9. Operational visibility

Log:
- timestamp, service, severity
- request ID and trace ID
- order ID or event ID where appropriate
- operation name, outcome, duration
- retry count and safe error classification

Never log passwords, access tokens, card data, or sensitive provider secrets.

Monitor:
- request rate, error rate, p95/p99 latency
- database connection pool usage
- queue depth and oldest-message age
- retry and dead-letter counts
- expired inventory reservations
- outbox backlog
- payment reconciliation backlog
