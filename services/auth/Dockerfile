FROM golang:1.26.2-alpine AS builder
WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o auth-server ./services/auth/cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/auth-server .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations/auth ./migrations
COPY --from=builder /app/services/auth/internal/api ./api

CMD goose -dir ./migrations postgres "$DB_DSN" up && ./auth-server