FROM golang:1.26.5-alpine AS builder

WORKDIR /app

# Root module
COPY go.mod go.sum ./

# Cart module
COPY cart/go.mod cart/go.sum ./cart/

# Download Cart dependencies
WORKDIR /app/cart
RUN go mod download

# Copy entire project
WORKDIR /app
COPY . .

# Build Cart
WORKDIR /app/cart
RUN CGO_ENABLED=0 go build -o cart ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cart/cart .

EXPOSE 8082

CMD ["./cart"]