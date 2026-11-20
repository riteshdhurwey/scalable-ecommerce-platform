FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY catalog/go.mod catalog/go.sum ./catalog/

WORKDIR /app/catalog
RUN go mod download

WORKDIR /app
COPY . .

WORKDIR /app/catalog
RUN CGO_ENABLED=0 go build -o catalog ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/catalog/catalog .

EXPOSE 8080

CMD ["./catalog"]