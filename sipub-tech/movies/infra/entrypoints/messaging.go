package entrypoints

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/messaging/constants"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/messaging/rabbitmq"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/ports"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/infra/controllers"
)

func NewMessagingEntrypoint(repo ports.MovieExecuteRepository, connectionUrl string) *MessagingEntrypoint {
	nodeId := uuid.New().String()
	return &MessagingEntrypoint{
		repo: repo,
		client: rabbitmq.NewRabbitMqServer(connectionUrl, nodeId),
		controller: &controllers.MessagingMovieController{},
	}
}

func NewMessagingEntrypointFromClient(repo ports.MovieExecuteRepository, client *rabbitmq.RabbitMqServer)  *MessagingEntrypoint {
	return &MessagingEntrypoint{
		repo: repo,
		client: client,
		controller: &controllers.MessagingMovieController{},
	}
}

type MessagingEntrypoint struct {
	repo ports.MovieExecuteRepository
	client *rabbitmq.RabbitMqServer
	controller *controllers.MessagingMovieController
}

func (entrypoint *MessagingEntrypoint) Serve(ctx context.Context) {
	entrypoint.client.Open()

	entrypoint.client.RegisterConsumer(constants.MovieCreatorQueueName, nil, nil, func(ctx context.Context, body any) error {
		ctx = context.WithValue(ctx, controllers.RepoKey, entrypoint.repo)
		dto, err := entrypoint.parseCreateDtoMap(body)
		if err != nil {
			return fmt.Errorf("couldn't parse body %+v to CreateMovieDTO", body)
		}

		return entrypoint.controller.SaveMovie(ctx, *dto)
	})

	entrypoint.client.RegisterConsumer(constants.MovieDeleterQueueName, nil, nil, func(ctx context.Context, body any) error {
		ctx = context.WithValue(ctx, controllers.RepoKey, entrypoint.repo)
		idMap, ok := body.(map[string]any)
		if !ok {
			return fmt.Errorf("couldn't parse body %+v to a map with an id", body)
		}

		return entrypoint.controller.DeleteMovie(ctx, dtos.MovieID(idMap["id"].(float64)))
	})

	entrypoint.client.Listen(ctx)
}

func (entrypoint *MessagingEntrypoint) Close() {
	entrypoint.client.Close()
}


func (entrypoint *MessagingEntrypoint) GetClient() *rabbitmq.RabbitMqServer {
	return entrypoint.client
}

func (entrypoint *MessagingEntrypoint) parseCreateDtoMap(rawDto any) (*dtos.CreateMovieDTO, error) {
	dtoMap, ok := rawDto.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("couldn't parse body %+v to a map", rawDto)
	}

	title, ok := dtoMap["title"]
	if !ok {
		return nil, fmt.Errorf("raw DTO did not have title key")
	}
	year, ok := dtoMap["year"]
	if !ok {
		return nil, fmt.Errorf("raw DTO did not have year key")
	}
	return &dtos.CreateMovieDTO{
		Title: title.(string),
		Year: year.(string),
	}, nil
}
