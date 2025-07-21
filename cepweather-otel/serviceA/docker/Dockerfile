FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . /app/
RUN go mod download
RUN go test ./tests/
RUN go build -o servicea ./cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/servicea .
EXPOSE 8080
CMD ["/app/servicea"]
