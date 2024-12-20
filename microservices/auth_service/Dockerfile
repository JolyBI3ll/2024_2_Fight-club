# 1. Build it

FROM golang:1.23.1 AS builder
WORKDIR /app
# Копируем go.mod и go.sum
COPY go.mod go.sum ./ 
RUN go mod download

# This microservice uses other modules so we can't just copy only it
# Therefore we need to copy the whole fucking project
# I have wasted 3 hours of my life on this
# COPY ./microservices/auth_service/ ./microservices/auth_service/
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o /microservices/auth_service/cmd/auth_service ./microservices/auth_service/cmd/main.go


# 2. Run it
FROM alpine:latest
# WORKDIR /microservices/auth_service
COPY --from=builder ./microservices/auth_service/cmd/auth_service /app/auth_service
CMD ["/app/auth_service"]