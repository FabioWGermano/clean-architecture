# clean-architecture

Projeto contendo o exercício Clean-Architecture do curso Go-Experts.
Nesse projeto buscamos informações de cep em dois serviços diferentes em multithreading e devolvemos o mais rapido como resultado

## Tecnologias Utilizadas

- Visual Studio
- Golang

## Instalação

- Baixar e configurar o docker
- Baixar o projeto https://github.com/fabiowgermano/clean-architecture
- Executar o comando `docker-compose up -d` para criar a imagem do banco de dados no docker
- Executar o comando `make migrate` pra criar a tabela e alguns dados de testes

- Acessar a pasta ordersystem com o comando `cd cmd/ordersystem` e executar o comando `go run main.go wire_gen.go` para iniciar a aplicação

- Para acessar a aplicação graphQL - http://localhost:8080/
    - Listar as Orders que foram cadastradas no banco de dados:
        query orders {
            listOrders{
                id,
                Price,
                Tax,
                FinalPrice
            }
        }
    - Adicionar nova Order:
        mutation createOrder {
            createOrder(input: {id: "1", Price:2.0, Tax:3.0}) {id}
        }

- Executar a listagem de orders por serviço `curl --location --request GET 'http://localhost:8000/order'`
- Executar a criação de uma order por serviço `curl --location 'http://localhost:8000/order' --header 'Content-Type: application/json' --data '{"id": "4","price": 4,"tax": 4}'`


Executar com o comando "go run cmd/main.go"