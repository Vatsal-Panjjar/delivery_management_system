# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o delivery_app ./cmd/server

# Run stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/delivery_app .
COPY web ./web
COPY migrations ./migrations

ENV PORT=8080

EXPOSE 8080

CMD ["./delivery_app"]
