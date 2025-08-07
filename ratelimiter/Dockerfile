# Dockerfile para Rate Limiter em Go
FROM golang:1.21-alpine AS builder

# Instalar dependências necessárias
RUN apk add --no-cache git

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de módulo Go
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o servidor ./cmd/servidor

# Estágio final - imagem mínima
FROM alpine:latest

# Instalar certificados CA para HTTPS
RUN apk --no-cache add ca-certificates

# Criar usuário não-root
RUN adduser -D -s /bin/sh appuser

# Definir diretório de trabalho
WORKDIR /root/

# Copiar binário do estágio de build
COPY --from=builder /app/servidor .

# Copiar arquivo de configuração
COPY --from=builder /app/.env .

# Mudar para usuário não-root
USER appuser

# Expor porta 8080
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./servidor"]