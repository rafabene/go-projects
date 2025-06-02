# syntax=docker/dockerfile:1

# Builder stage
FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/ordersystem
RUN CGO_ENABLED=0 GOOS=linux go build -o ordersystem .

# Final image
FROM alpine:latest

RUN apk add --no-cache bash

COPY --from=builder /app /app

ADD --chmod=770 https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /app/wait-for-it.sh

EXPOSE 8080 50051 8000
WORKDIR /app
CMD ["./ordersystem"]