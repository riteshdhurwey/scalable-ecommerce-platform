# Testing and Scalability

## 1. Testing layers

### Unit tests
Test domain logic, validation, price calculations, state transitions, retry decisions, and event handling.

### Integration tests
Use real PostgreSQL and Redis instances in an isolated test environment. Test repository behavior, transactions, uniqueness constraints, migrations, and service-to-database interactions.

### End-to-end tests
Exercise the customer journey across services:
1. Create or authenticate a user.
2. Browse products.
3. Add items to the cart.
4. Submit checkout.
5. Simulate payment success or failure.
6. Verify order and inventory states.
7. Verify notification event delivery.

### Failure tests
Test timeouts, duplicate requests, duplicate events, service restarts, unavailable dependencies, and worker crashes.

## 2. Critical concurrency tests

### Last-unit inventory test
- Create a product with one unit.
- Submit many concurrent checkout/reservation attempts.
- Verify no more than one reservation succeeds.
- Verify stock never becomes negative.
- Verify losing attempts do not proceed to payment.

### Duplicate checkout test
- Send the same checkout request concurrently with the same idempotency key.
- Verify it maps to one logical order.
- Verify reusing the key with a different payload is rejected.

### Duplicate payment webhook test
- Deliver the same verified success event repeatedly.
- Verify the payment transition and order confirmation are applied once logically.

### Outbox and consumer recovery test
- Crash or stop a publisher after business commit but before publication.
- Restart it and verify the event is eventually published.
- Simulate duplicate delivery and verify the consumer does not duplicate the business effect.

## 3. Load testing

Measure a baseline before optimizing. Test catalog reads and checkout separately because they have different resource and consistency requirements.

Record:
- requests per second
- p50, p95, and p99 latency
- HTTP error rate and timeouts
- CPU and memory usage
- database connection pool usage
- Redis latency
- broker queue depth and oldest-message age
- order/payment failure and retry counts

Increase concurrency in controlled steps. Keep test data and environment details reproducible. Do not present a benchmark number without documenting the test setup and workload.

## 4. Scaling strategy

- Keep request-handling services stateless where practical.
- Scale service replicas based on measured bottlenecks.
- Use database indexes based on query plans and workload.
- Use caching for appropriate read-heavy data, with explicit invalidation/expiry behavior.
- Bound worker concurrency and queue growth.
- Apply backpressure rather than accepting unlimited work.
- Monitor database connection limits before increasing replicas.
- Consider read replicas, partitioning, or Kubernetes only when measured needs justify the additional complexity.

## 5. Example dashboards

### API
- request rate by route
- error rate by status class
- p95/p99 latency

### Workers and messaging
- queue depth
- oldest message age
- processing duration
- retry and dead-letter counts

### Business workflow
- orders by status
- payment success/failure/pending counts
- expired reservations
- outbox backlog
- reconciliation backlog

## 6. Resume-quality evidence

A strong project report should include:
- architecture diagram
- documented API and state transitions
- automated test commands and results
- a reproducible concurrency test
- load-test methodology and measured results
- failure-recovery scenarios and expected behavior
- known limitations and future improvements

Be precise about what was actually tested. A Dockerized microservices demo is not automatically production-ready.
