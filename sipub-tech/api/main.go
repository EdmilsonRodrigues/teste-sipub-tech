package main

import (
	"os"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/entrypoints"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/services"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/controllers"
)

func main() {
	rabbitmqConnectionURL := os.Getenv("API_GATEWAY_RABBITMQ_CONNECTION_URL")
	executorService := services.NewMovieMessagingService(rabbitmqConnectionURL)
	defer executorService.Close()

	gRPCConnectionUrl := os.Getenv("API_GATEWAY_GRPC_CONNECTION_URL")
	queryService := services.NewMovieGRPCService(gRPCConnectionUrl)
	defer queryService.Close()
	
	movieController := controllers.NewMovieController()

	server := entrypoints.NewGinEntrypoint(
		executorService,
		queryService,
		movieController,
	)
	server.Setup()

	server.Serve()
}
