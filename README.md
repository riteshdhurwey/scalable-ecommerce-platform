# Scalable E-Commerce Platform

A backend-focused e-commerce platform built with **Go, PostgreSQL, Redis, and REST APIs**, following a microservices architecture.

The project focuses on clean service boundaries, authentication, database design, concurrency control, inventory management, and building reliable backend systems.

## Tech Stack

| Category | Technology |
|---|---|
| Language | Go |
| HTTP Framework | Gin |
| Database | PostgreSQL |
| Cache | Redis |
| Authentication | JWT, HTTP Cookies |
| Password Security | bcrypt |
| API Communication | REST / HTTP |
| Testing | Go Testing |
| Load Testing | k6 |
| Containerization | Docker / Docker Compose |
| Version Control | Git / GitHub |

Some infrastructure components and integrations are planned or under development.

## Architecture

The platform is organized into independently developed backend services.

| Service | Responsibility | Status |
|---|---|---|
| User Service | Registration, login, password hashing, JWT authentication | Implemented |
| Catalog Service | Product and category management, product listing and filtering | Core features implemented |
| Cart Service | Shopping cart management and product quantity handling | In Progress |
| Order Service | Order creation, inventory reservation, and order lifecycle | In Progress |
| Payment Service | Payment processing and payment status management | Planned |
| Notification Service | Order and payment notifications | Planned |

For detailed service boundaries, communication patterns, and infrastructure decisions, see [Architecture Documentation](docs/architecture.md).

## Implemented Features

### User Service

- User registration with input validation.
- Email validation and duplicate email/user ID checks.
- Password hashing using bcrypt.
- Password verification during login.
- JWT access token generation.
- HTTP authentication cookies with a 15-minute lifetime.
- Protected endpoints using authentication middleware.
- Layered architecture with handlers, services, repositories, middleware, and models.
- Database migrations for user data.

### Catalog Service

- Product creation and retrieval.
- Product category creation and retrieval.
- Product listing with pagination.
- Product filtering by category.
- Category information in product responses.
- PostgreSQL integration.
- Layered architecture using handlers, services, repositories, DTOs, and models.
- Database migrations for products and categories.

### Cart and Order Services

Currently under development.

The design includes:

- Cart item management.
- Product information retrieval from Catalog Service.
- Quantity validation.
- Order creation and order item snapshots.
- Inventory reservation.
- Order expiration and inventory release.

## Project Structure

```text
scalable-ecommerce-platform/
├── services/
│   ├── user/
│   │   ├── cmd/
│   │   │   └── api/
│   │   ├── database/
│   │   ├── internal/
│   │   │   ├── dto/
│   │   │   ├── handler/
│   │   │   ├── jwt/
│   │   │   ├── middleware/
│   │   │   ├── model/
│   │   │   ├── repository/
│   │   │   ├── routes/
│   │   │   └── service/
│   │   └── migrations/
│   │
│   ├── catalog/
│   │   ├── cmd/
│   │   │   ├── api/
│   │   │   └── seed/
│   │   ├── internal/
│   │   │   ├── database/
│   │   │   ├── dto/
│   │   │   ├── handler/
│   │   │   ├── models/
│   │   │   ├── repository/
│   │   │   └── service/
│   │   └── migrations/
│   │
│   ├── cart/
│   └── order/
│
├── deployments/
│   └── docker-compose.yml
│
├── docs/
│   ├── README.md
│   ├── architecture.md
│   ├── checkout-workflow.md
│   ├── service-workflows.md
│   ├── failure-recovery.md
│   ├── data-ownership.md
│   ├── development-roadmap.md
│   └── testing-and-scalability.md
│
├── .github/
│   └── workflows/
│
├── .gitignore
└── README.md
```

## Prerequisites

- Go
- PostgreSQL
- Git

Redis and Docker are required for the corresponding infrastructure integrations as they are introduced.

## Run Locally

Clone the repository:

```bash
git clone https://github.com/riteshdhurwey/scalable-ecommerce-platform.git
cd scalable-ecommerce-platform
```

### Start User Service

```bash
cd services/user
go mod download
go run ./cmd/api
```

### Start Catalog Service

Open another terminal:

```bash
cd services/catalog
go mod download
go run ./cmd/api
```

Configure the required environment variables and PostgreSQL database connections before starting each service.

Apply the corresponding database migrations before using the APIs.

## Testing

Run tests for each service:

```bash
cd services/user
go test ./...
```

```bash
cd services/catalog
go test ./...
```

The test suite is still being expanded. Package compilation alone should not be considered evidence of comprehensive test coverage.

## Documentation

Detailed architecture and implementation decisions are maintained separately.

| Document | Description |
|---|---|
| [Architecture](docs/architecture.md) | Service boundaries, communication, and infrastructure |
| [Checkout Workflow](docs/checkout-workflow.md) | Order placement, inventory reservation, and payment |
| [Service Workflows](docs/service-workflows.md) | Responsibilities and workflows of individual services |
| [Failure Recovery](docs/failure-recovery.md) | Retries, idempotency, timeouts, and recovery |
| [Data Ownership](docs/data-ownership.md) | Database ownership and persistence rules |
| [Development Roadmap](docs/development-roadmap.md) | Implementation phases and planned features |
| [Testing and Scalability](docs/testing-and-scalability.md) | Testing, concurrency, reliability, and performance |

## Roadmap

- [x] User registration and authentication.
- [x] JWT authentication middleware and protected endpoints.
- [x] Product and category management.
- [x] Product listing, filtering, and pagination.
- [ ] Complete Cart Service.
- [ ] Complete Order Service and checkout workflow.
- [ ] Implement inventory reservation and release.
- [ ] Integrate payment processing.
- [ ] Implement order expiration and recovery workers.
- [ ] Integrate Redis where appropriate.
- [ ] Introduce asynchronous messaging with RabbitMQ.
- [ ] Add Docker-based local development.
- [ ] Expand unit, integration, and end-to-end tests.
- [ ] Add structured logging, metrics, and tracing.
- [ ] Add CI workflows and performance benchmarks.

## Project Goals

The project aims to demonstrate practical backend engineering skills:

- Go backend development and REST API design.
- Microservices architecture and service boundaries.
- PostgreSQL schema design and query optimization.
- Authentication and security fundamentals.
- Inventory management and concurrency control.
- Caching and asynchronous processing.
- Distributed systems fundamentals.
- Automated testing, observability, and reliability.

---

**Development Status:** Actively under development.
