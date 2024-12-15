# 1. Build it
FROM golang:1.23.1 AS builder
WORKDIR /app
COPY . .
RUN go mod download

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o /app/backend ./cmd/webapp


# 2. Run it
FROM alpine:latest
COPY --from=builder /app/backend /app/backend
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/ssl/ /app/ssl/

WORKDIR /app
CMD ["/app/backend"]