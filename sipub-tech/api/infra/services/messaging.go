package services

import (
	"context"
	"fmt"
	
	"github.com/google/uuid"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/messaging/constants"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/messaging/rabbitmq"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/ports"
)

func NewMovieMessagingService(connectionUrl string) *MovieMessagingService {
	nodeId := uuid.New().String()

	client := rabbitmq.NewRabbitMqServer(connectionUrl, nodeId)
	client.Open()

	_, saver   := client.CreateProducer(constants.MovieCreatorQueueName, nil, nil)
	_, deleter := client.CreateProducer(constants.MovieDeleterQueueName, nil, nil)
	
	return &MovieMessagingService{
		client: client,
		save:   saver,
		delete: deleter,		
	}
}

type IdBody struct {
	Id dtos.MovieId  `json:"id"`
}

type MovieMessagingService struct {
	ports.MovieExecutorService

	client *rabbitmq.RabbitMqServer
	save   rabbitmq.ProducerFunction
	delete rabbitmq.ProducerFunction
}

func (service *MovieMessagingService) Close() {
	service.client.Close()
}

func (service *MovieMessagingService) GetClient() *rabbitmq.RabbitMqServer {
	return service.client
}

func (service *MovieMessagingService) Save(ctx context.Context, movie dtos.CreateMovieDTO) error {
	if err := service.save(ctx, movie); err != nil {
		return fmt.Errorf("failed saving movie %+v: %w", movie, err)
	}
	return nil
}

func (service *MovieMessagingService) Delete(ctx context.Context, id dtos.MovieId) error {
	dto := IdBody{Id: id}
	if err := service.delete(ctx, dto); err != nil {
		return fmt.Errorf("failed deleting movie with id %d: %w", id, err)
	}
	return nil
}
