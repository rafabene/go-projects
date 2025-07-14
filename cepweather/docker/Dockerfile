# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . /app/
RUN go mod download
RUN go build -o cepweather ./cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/cepweather .
COPY ./configs/.env ./configs/.env
EXPOSE 8080
CMD ["/app/cepweather"]
