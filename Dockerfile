# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git (for go modules) + ca-certs
RUN apk add --no-cache git ca-certificates tzdata

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build (static-ish)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /app/bin/darulabror-api ./cmd/echo-server

# Runtime stage
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && adduser -D -H -u 10001 appuser

COPY --from=builder /app/bin/darulabror-api /app/darulabror-api

USER appuser

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/darulabror-api"]