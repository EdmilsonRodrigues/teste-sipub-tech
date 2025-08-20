//go:build integration

package entrypoints_test


import (
    "context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/messaging/constants"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/domain"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/infra/entrypoints"
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


	repository := &MockMovieExecuteRepository{}
	entrypoint := entrypoints.NewMessagingEntrypoint(repository, connectionUrl)
	client := entrypoint.GetClient()

	t.Run("should be able to create and delete a repository.", func (t *testing.T) {
		createDto := dtos.CreateMovieDTO{
			Title: faker.Sentence(),
			Year: faker.YearString(),
		}
		id := rand.Int31()

		entrypoint.Serve(ctx)
		defer entrypoint.Close()

		expected := domain.Movie{
			Title: createDto.Title,
			Year: createDto.Year,
		}

		test := "Should be able to send a dtos.CreateMovieDTO to message queue"
		logTest(t, test)
		_, sendCreateMovie := client.CreateProducer(constants.MovieCreatorQueueName, nil, nil)
		if err := sendCreateMovie(ctx, createDto); err != nil {
			logError(t, "Error sending dtos.CreateMovieDTO %+v to messaging queue.", createDto)
		} else {
			logSuccess(t, test)
		}

		test = "Should be able to send the id to message queue"
		logTest(t, test)
		_, sendDeleteMovie := client.CreateProducer(constants.MovieDeleterQueueName, nil, nil)
		if err := sendDeleteMovie(ctx, map[string]int{"id": int(id)}); err != nil {
			logError(t, "Error sending id %d to messaging queue.", id)
		} else {
			logSuccess(t, test)
		}

		test = "The message sent should have added the bodies to the repository"
		logTest(t, test)
		time.Sleep(1 * time.Second)
		if !reflect.DeepEqual(expected, repository.MoviePassed) {
			logError(t, "Expected: %+v, Got: %+v", expected, repository.MoviePassed)
		} else if int(id) != repository.IdPassed {
			logError(t, "Expected: %+v, Got: %+v", expected, repository.IdPassed)
		} else {
			logSuccess(t, test)
		}
		
		
	})
}

type MockMovieExecuteRepository struct {
	MoviePassed domain.Movie
	IdPassed int
	ErrorReturned error
}

func (repo *MockMovieExecuteRepository) Save(ctx context.Context, movie domain.Movie) error {
	repo.MoviePassed = movie
	return repo.ErrorReturned
}

func (repo *MockMovieExecuteRepository) Delete(ctx context.Context, id int) error {
	repo.IdPassed = id
	return repo.ErrorReturned
}

func insertAuthInfo(endpoint, authInfo string) string {
	parts := strings.Split(endpoint, "//")
	return parts[0] + "//" + authInfo + "@" + parts[1]
}

func logTest(t *testing.T, msg string, args ...any) {
	formattedMsg := fmt.Sprintf("[TEST]: %s", msg)
	t.Logf(formattedMsg, args...)
}

func logError(t *testing.T, msg string, args ...any) {
	formattedMsg := fmt.Sprintf("[ERROR]: %s", msg)
	t.Errorf(formattedMsg, args...)
}

func logSuccess(t *testing.T, msg string, args ...any) {
	formattedMsg := fmt.Sprintf("[SUCCESS]: %s", msg)
	t.Logf(formattedMsg, args...)
}

