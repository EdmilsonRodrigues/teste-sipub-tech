package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/domain"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/infra/repositories"
)

func main() {
	awsRegion := os.Getenv("AWS_REGION")
	dynamoDBEndpoint := os.Getenv("DYNAMO_DB_ENDPOINT")
	jsonPath := os.Getenv("JSON_PATH")

	ctx := context.Background()
	repo := repositories.NewMovieRepository(repositories.NewRepositoryConfig(awsRegion, dynamoDBEndpoint))
	repo.Open()

	if err := repo.CreateTables(ctx); err != nil {
		panic(fmt.Sprintf("Failed to create tables: %v", err))
	}

	ready := make(chan bool)
	var movies []domain.Movie

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		panic(fmt.Sprintf("error reading file: %v", err))
	}

	if err := json.Unmarshal(data, &movies); err != nil {
		panic(fmt.Sprintf("failed to unmarshall movies data, here's why: %v", err))
	}

	save := func(movie domain.Movie) {
		if err := repo.SaveWithId(ctx, movie); err != nil {
			panic(fmt.Sprintf("failed to save movie %+v to db: %v", movie, err))
		}
		ready <- true
	}

	for _, movie := range movies {
		go save(movie)
	}

	for range len(movies) {
		<-ready
	}
}
