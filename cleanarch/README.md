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

- Go 1.23+
- Docker e Docker Compose

### Rodando a aplicação

Agora, utilize o Docker Compose para buildar e rodar toda a stack (aplicação e banco de dados):

```sh
docker-compose up --build
```

As migrations serão executadas automáticamente.

A aplicação estará disponível em:

- REST: http://localhost:8000
- GraphQL Playground: http://localhost:8080/graphql
- gRPC: consulte protofiles em `internal/infra/grpc/protofiles/`

## Exemplos de uso

### GraphQL

Acesse <http://localhost:8080/graphql> no browser

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

Em outra sessão do terminal, em qualquer path.

#### Criar pedido

```sh
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"id": "2", "Price": 200, "Tax": 20}'
```

#### Listar pedidos

```sh
curl http://localhost:8000/order
```

### gRPC (usando evans previamente instalado)

Em outro terminal, à partir da pasta `cleanarch`(raiz deste projeto).

#### Iniciar o evans

```sh
~/go/bin/evans internal/infra/grpc/protofiles/order.proto
```

#### Criar pedido

```sh
call CreateOrder
{
  "id": "3",
  "price": 300,
  "tax": 30
}
```

#### Listar pedidos

```sh
call ListOrders
{}
```

#### Sair

```sh
exit
```

## Licença

Apache-2.0
