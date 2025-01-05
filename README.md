# clean-architecture

Projeto contendo o exercício Clean-Architecture do curso Go-Experts.
API criar e listar ordens de serviço, informando preço e taxa, retornando preço final.

## Requisitos
- [x] API deve conter servidor REST, GraphQL e gRPC.
    - [x] Endpoint REST (GET /order)
    - [x] Service ListOrders com GRPC
    - [x] Query ListOrders GraphQL
- [x] Requests para criar e listar orders no arquivo api.http
- [x] Migrações necessárias
- [x] Usar docker e docker-compose e executar em containers
- [x] Documentar projeto e como executá-lo

## Tecnologias Utilizadas

- wire
- viper
- testify
- amqp(rabbitmq)
- net/http, grpc, graphql
- mysql

## Instalação

1. Instale as dependências e suba os containers
``` shell
make install
make up
```

2. Como testar a aplicação: REST API server
- faça uma chamada POST para criar uma nova order via rest-client usando o arquivo [api.http](api/api.http)
- faça uma chamada GET para listar as orders via rest-client usando o arquivo [api.http](api/api.http)

3. Como testar a aplicação: gRPC server
``` shell
## criar nova order
evans --proto internal/infra/grpc/protofiles/order.proto --host localhost --port 50051
=> call CreateOrder
=> 2
=> 10.5
=> 0.5

## listar orders
evans --proto internal/infra/grpc/protofiles/order.proto --host localhost --port 50051
=> call ListOrders
```

4. Testar aplicação: GraphQL server
``` shell
## criar nova order
mutation createOrder {
  createOrder(input: {id:"1", Price: 1.1, Tax: 0.1}) {
    id
    Price
    Tax
  }
}

## listar orders
query queryOrders {
  listOrders {
    id
    Price
    Tax
    FinalPrice
  }
}
```

5. Baixar containers 
``` shell
make down
```