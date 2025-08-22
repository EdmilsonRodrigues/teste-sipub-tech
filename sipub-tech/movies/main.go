package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/infra/repositories"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/infra/entrypoints"
)

func main() {
	awsRegion := os.Getenv("MOVIE_SERVICE_AWS_REGION")
	dynamoDBEndpoint := os.Getenv("MOVIE_SERVICE_DYNAMO_DB_ENDPOINT")
	grpcListeningPort := os.Getenv("MOVIE_SERVICE_GRPC_LISTENING_PORT")
	rabbitmqConnectionURL := os.Getenv("MOVIE_SERVICE_RABBITMQ_CONNECTION_URL")
	listeningPort, err := strconv.Atoi(grpcListeningPort)
	if err != nil {
		panic("malformed grpc listening port configuration")
	}
	
	ctx := context.Background()
	repo := repositories.NewMovieRepository(repositories.NewRepositoryConfig(awsRegion, dynamoDBEndpoint))
	repo.Open()
	if err := repo.CreateTables(ctx); err != nil {
		panic(fmt.Sprintf("Failed to create tables: %v", err))
	}

	grpcEntrypoint := entrypoints.NewGRPCEntrypoint(repo, listeningPort)
	messagingEntrypoint := entrypoints.NewMessagingEntrypoint(repo, rabbitmqConnectionURL)

	messagingEntrypoint.Serve(ctx)
	defer messagingEntrypoint.Close()

	grpcEntrypoint.Serve(ctx)
}
