# Data Ownership and Storage

## 1. Ownership rules

Each service is the authoritative owner of its own data. Other services use APIs or events instead of querying its tables directly.

| Service | Core entities | Storage suggestion |
|---|---|---|
| User | User, Profile, Credentials metadata | PostgreSQL |
| Catalog/Inventory | Product, Category, Stock, Reservation | PostgreSQL |
| Cart | Cart, CartItem | Redis initially; persistence if required |
| Order | Order, OrderItem, OrderStateHistory, OutboxEvent | PostgreSQL |
| Payment | PaymentAttempt, ProviderEvent, OutboxEvent | PostgreSQL |
| Notification | Notification, DeliveryAttempt, ProcessedEvent | PostgreSQL |

A single PostgreSQL instance is acceptable for local development. Separate databases or schemas can establish ownership boundaries. Use separate credentials or permissions where practical.

## 2. Important data rules

- Store money as integer minor units, such as paise, with an explicit currency.
- Save a price snapshot and product description/name on order items so historical orders remain meaningful after catalog changes.
- Use database constraints for uniqueness and valid relationships within a service.
- Use migrations instead of manual production schema edits.
- Define retention policies for payment events, logs, carts, and idempotency records.
- Keep secrets out of database seed files and version control.

## 3. Inventory consistency

Inventory reservation must be atomic. Avoid reading stock in one query and later decrementing it without concurrency control.

For an order with several products, the reservation should succeed for all requested items or fail as a unit, unless partial reservation is a deliberate product feature.

## 4. Order and payment boundaries

The Order Service owns order state; the Payment Service owns payment attempts and provider outcomes. They coordinate through APIs/events and explicit state transitions.

The order should retain a reference to its payment attempt, but it should not mutate payment tables. The payment service should not directly update order tables.

## 5. Redis guidance

Use Redis for data whose consistency and durability requirements match its configured persistence and recovery behavior.

Possible initial uses:
- Active carts
- Product cache
- Rate limiting
- Short-lived coordination data

Do not make Redis the only authoritative inventory store. If carts must survive Redis data loss, persist them or configure and test a suitable durability/recovery strategy.
