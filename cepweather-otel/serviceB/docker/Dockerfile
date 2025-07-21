FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . /app/
RUN go mod download
RUN go test ./tests/
RUN go build -o serviceb ./cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/serviceb .
COPY ./configs/.env ./configs/.env
EXPOSE 8081
CMD ["/app/serviceb"]
