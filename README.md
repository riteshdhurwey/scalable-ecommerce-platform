# Scalable E-Commerce Platform

A backend-focused e-commerce platform built with **Go, PostgreSQL, Redis, and REST APIs**, designed around microservices architecture. The project focuses on clean service boundaries, authentication, database design, concurrency, and building reliable backend systems.

## Tech Stack

- **Language:** Go
- **HTTP Framework:** Gin
- **Database:** PostgreSQL
- **Cache and Queueing:** Redis (planned integration across additional services)
- **Authentication:** JWT and HTTP cookies
- **Password Security:** bcrypt
- **Testing:** Go testing framework
- **Load Testing:** k6
- **Containerization:** Docker (planned/in progress)

## Architecture

The platform is being developed as a collection of independently organized backend services.

| Service | Responsibility | Status |
|---|---|---|
| User Service | User registration, login, password hashing, JWT authentication | Implemented |
| Catalog Service | Product and category management | Core features implemented |
| Cart Service | Shopping cart management | Planned |
| Order Service | Order creation and lifecycle management | Planned |
| Payment Service | Payment workflow and status management | Planned |
| Notification Service | Order and payment notifications | Planned |

## Implemented Features

### User Service

- Create user accounts with input validation.
- Validate email addresses and check for duplicate email addresses and user IDs.
- Hash passwords using bcrypt.
- Verify passwords during login.
- Generate JWT access tokens.
- Set authentication cookies with a 15-minute lifetime.
- Organize code into handler, service, repository, middleware, and model layers.
- Define database migrations for user data.

### Catalog Service

- Create and retrieve product records.
- Create and retrieve product categories.
- Retrieve products by category.
- Implement pagination for product listings.
- Use PostgreSQL for persistent storage.
- Organize code into handlers, services, repositories, DTOs, and models.
- Define database migrations for categories and products.

## Project Structure

```text
scalable-ecommerce-platform/
├── services/
│   ├── user/
│   │   ├── cmd/api/
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
│   └── catalog/
│       ├── cmd/
│       │   ├── api/
│       │   └── seed/
│       ├── internal/
│       │   ├── database/
│       │   ├── dto/
│       │   ├── handler/
│       │   ├── models/
│       │   ├── repository/
│       │   └── service/
│       └── migrations/
└── README.md
```

## Prerequisites

- Go installed
- PostgreSQL running
- Git installed
- Redis and Docker when required by the additional services being developed

## Run Locally

Clone the repository:

```bash
git clone https://github.com/riteshdhurwey/scalable-ecommerce-platform.git
cd scalable-ecommerce-platform
```

Run the User Service:

```bash
cd services/user
go mod download
go run ./cmd/api
```

Run the Catalog Service in a separate terminal:

```bash
cd services/catalog
go mod download
go run ./cmd/api
```

Configure the database connection and other required settings according to each service's configuration before starting it. Apply the relevant database migrations before using endpoints that require the corresponding tables.

## Testing

Run the available Go tests for each service:

```bash
cd services/user
go test ./...
```

```bash
cd services/catalog
go test ./...
```

These commands currently verify package compilation. Automated unit and integration test coverage is still being developed.

## Roadmap

- [ ] Complete JWT authentication middleware and protected endpoints.
- [ ] Add product filtering, update/delete operations, and inventory validation.
- [ ] Implement cart management and persistence.
- [ ] Implement order creation and order lifecycle management.
- [ ] Integrate payment processing and notifications.
- [ ] Integrate Redis where appropriate for caching and asynchronous processing.
- [ ] Add Docker-based local development.
- [ ] Add unit and integration tests.
- [ ] Add structured logging, metrics, tracing, and health checks.
- [ ] Add CI workflows and performance benchmarks.

## Goals

The goal is to develop a production-oriented backend that demonstrates practical skills in:

- Go backend development and REST API design
- Microservices architecture and service boundaries
- PostgreSQL schema design and query optimization
- Authentication and security fundamentals
- Caching, concurrency, and distributed systems
- Automated testing, observability, and reliability
