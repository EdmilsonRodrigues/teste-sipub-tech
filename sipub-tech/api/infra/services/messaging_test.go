//go:build integration

package services_test


import (
    "context"
	"fmt"
	"strings"
	"testing"
	"time"
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/messaging/constants"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/services"
)


const (
	rabbitMqImage = "rabbitmq:3"
	rabbitMqConnectionPort = "5672"
	rabbitMqUser = "guest"
	rabbitMqPassword = "guest"
	rabbitMqProtocol = "amqp"
)

var (
	rabbitMqExposedPorts = []string{"4369/tcp", "5672/tcp"}
)

func TestMessagingEntrypoint(t *testing.T) {
	
    ctx := context.Background()

    req := testcontainers.ContainerRequest{
        Image:        rabbitMqImage,
        ExposedPorts: rabbitMqExposedPorts,
        WaitingFor:   wait.ForLog("Ready to start client connection listeners"),
    }
    rabbitmqC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    testcontainers.CleanupContainer(t, rabbitmqC)
    require.NoError(t, err)

	endpoint, err := rabbitmqC.PortEndpoint(ctx, rabbitMqConnectionPort, rabbitMqProtocol)
	if err != nil {
		t.Error(err)
	}

	authInfo := fmt.Sprintf("%s:%s", rabbitMqUser, rabbitMqPassword)
	connectionUrl := insertAuthInfo(endpoint, authInfo)



	t.Run("should be able to send a create and a delete signal with the right bodies.", func (t *testing.T) {
		service := services.NewMovieMessagingService(connectionUrl)
		client := service.GetClient()
		defer service.Close()
		ready := make(chan bool)

		createDto := dtos.CreateMovieDTO{
			Title: faker.Sentence(),
			Year: faker.YearString(),
		}
		expectedCreateBody := map[string]any{
			"title": createDto.Title,
			"year":  createDto.Year,
		}
		client.RegisterConsumer(constants.MovieCreatorQueueName, nil, nil, func(ctx context.Context, body any) error {
			assert.Equal(t, expectedCreateBody, body)
			ready <- true
			return nil
		})

		id := rand.Int32()
		expected := map[string]any{
			"id": float64(id),
		}
		client.RegisterConsumer(constants.MovieDeleterQueueName, nil, nil, func(ctx context.Context, body any) error {
			assert.Equal(t, expected, body)
			ready <- true
			return nil
		})
		client.Listen(ctx)

		service.Save(ctx, createDto)
		service.Delete(ctx, dtos.MovieId(id))

		messagesCount := 2

		for range messagesCount {
			select {
			case <- ready:
				return
			case <- time.After(1* time.Second):
				t.Errorf("Did not reach consumer")
			}
		}
	})
}

func insertAuthInfo(endpoint, authInfo string) string {
	parts := strings.Split(endpoint, "//")
	return parts[0] + "//" + authInfo + "@" + parts[1]
}

