# This is a separate image for migrator
# We use this approach because other services depend on it

# 1. Build it
FROM golang:1.23.1 AS builder
WORKDIR /app
COPY . .
RUN go mod download

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o /app/migrator ./cmd/migrator


# 2. Run it
FROM alpine:latest
COPY --from=builder /app/migrator /app/migrator
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/cities_images/ /app/cities_images/
WORKDIR /app
CMD ["/app/migrator"]
