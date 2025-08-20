package main

import (
	"context"
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
	
	repo := repositories.NewMovieRepository(repositories.NewRepositoryConfig(awsRegion, dynamoDBEndpoint))
	ctx := context.Background()

	grpcEntrypoint := entrypoints.NewGRPCEntrypoint(repo, listeningPort)
	messagingEntrypoint := entrypoints.NewMessagingEntrypoint(repo, rabbitmqConnectionURL)

	messagingEntrypoint.Serve(ctx)
	defer messagingEntrypoint.Close()

	grpcEntrypoint.Serve(ctx)
}
