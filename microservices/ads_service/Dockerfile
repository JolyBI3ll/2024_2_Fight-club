# 1. Build it

FROM golang:1.23.1 AS builder
WORKDIR /app
# Копируем go.mod и go.sum
COPY go.mod go.sum ./ 
RUN go mod download

COPY .env ./
# This microservice uses other modules so we can't just copy only it
# Therefore we need to copy the whole fucking project
# I have wasted 3 hours of my life on this
# COPY ./microservices/ads_service/ ./microservices/ads_service/
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o /microservices/ads_service/cmd/ads_service ./microservices/ads_service/cmd/main.go


# 2. Run it
FROM alpine:latest
# WORKDIR /microservices/ads_service
COPY --from=builder ./microservices/ads_service/cmd/ads_service /app/ads_service
COPY --from=builder ./app/.env /.env
CMD ["/app/ads_service"]