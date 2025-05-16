# Clean Architecture Order System

Este projeto é um sistema de pedidos (Order System) desenvolvido em Go, seguindo os princípios da Clean Architecture. Ele utiliza GraphQL, REST, gRPC, eventos e injeção de dependências com Google Wire.

## Funcionalidades

- Cadastro de pedidos (Order)
- Listagem de pedidos
- API GraphQL (gqlgen)
- API REST
- API gRPC
- Eventos de domínio (OrderCreated, OrderListed)
- Injeção de dependências com Google Wire

## Estrutura do Projeto

```
cmd/ordersystem/         # Entrypoint da aplicação (main.go, wire.go, wire_gen.go)
internal/entity/        # Entidades de domínio
internal/event/         # Eventos de domínio
internal/event/handler/ # Handlers de eventos
internal/infra/database/# Infraestrutura de banco de dados
internal/infra/web/     # Handlers HTTP/REST
internal/infra/web/webserver/ # Webserver
internal/infra/graph/   # GraphQL (schema, resolvers, models)
internal/infra/grpc/    # gRPC (pb, protofiles, service)
internal/usecase/       # Casos de uso
pkg/events/             # Event dispatcher e interfaces
api/                    # Exemplos de requisições HTTP
configs/                # Configurações

```

## Como rodar o projeto

### Pré-requisitos

- Go 1.20+
- Docker e Docker Compose

### Subindo o banco de dados

```sh
docker-compose up -d
```

### Rodando as migrations

Elas serão executadas automaticamente

### Rodando a aplicação

```sh
go build
./ordersystem
```

A aplicação estará disponível em:

- REST: http://localhost:8080
- GraphQL Playground: http://localhost:8080/graphql
- gRPC: consulte protofiles em `internal/infra/grpc/protofiles/`

## Exemplos de uso

### GraphQL

#### Criar pedido

```graphql
mutation {
  createOrder(input: { id: "1", Price: 100, Tax: 10 }) {
    id
    Price
    Tax
    FinalPrice
  }
}
```

#### Listar pedidos

```graphql
query {
  getOrders {
    id
    Price
    Tax
    FinalPrice
  }
}
```

### REST

#### Criar pedido

```sh
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"id": "1", "Price": 100, "Tax": 10}'
```

#### Listar pedidos

```sh
curl http://localhost:8000/order
```

### gRPC (usando evans)

#### Iniciar o evans

```sh
evans internal/infra/grpc/protofiles/order.proto
```

#### Criar pedido

```sh
call CreateOrder
{
  "id": "1",
  "price": 100,
  "tax": 10
}
```

#### Listar pedidos

```sh
call ListOrders
{}
```

## Licença

Apache-2.0

