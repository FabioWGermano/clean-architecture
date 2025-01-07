package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fabiowgermano/clean-architecture/configs"
	"github.com/fabiowgermano/clean-architecture/internal/event/handler"
	"github.com/fabiowgermano/clean-architecture/internal/infra/graph"
	"github.com/fabiowgermano/clean-architecture/internal/infra/grpc/pb"
	"github.com/fabiowgermano/clean-architecture/internal/infra/grpc/service"
	"github.com/fabiowgermano/clean-architecture/internal/infra/web/webserver"
	"github.com/fabiowgermano/clean-architecture/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// ----- CONFIGS
	configs, err := configs.LoadConfig(".")
	print(configs)

	if err != nil {
		panic(err)
	}

	// ----- DATABASE
	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// ----- RABBITMQ
	rabbitMQChannel := getRabbitMQChannel(configs.RMQUrl)
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUsecase := NewListOrderUseCase(db, eventDispatcher)

	// ----- WEBSERVER
	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", "POST", webOrderHandler.Create)
	webserver.AddHandler("/order", "GET", webOrderHandler.List)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	// ----- GRPC_SERVER
	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUsecase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	// ----- GRAPHQL_SERVER
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrderUsecase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel(end string) *amqp.Channel {
	conn, err := amqp.Dial(end)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
