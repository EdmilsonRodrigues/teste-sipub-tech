//go:build integration

package repositories_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
	"github.com/go-faker/faker/v4"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/domain"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/ports"	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/infra/repositories"	
)

const (
	localStackImage = "localstack/localstack"
	localStackPort = "4566"
	awsRegion = "us-east-1"
)

var (
	localStackExposedPorts = []string{"4566/tcp", "4510-4559/tcp"}
)

func TestMovieRepository(t *testing.T) {
	
    ctx := context.Background()

    req := testcontainers.ContainerRequest{
        Image:        localStackImage,
        ExposedPorts: localStackExposedPorts,
        WaitingFor:   wait.ForLog("Ready."),
		Env: map[string]string{
			"SERVICES": "dynamodb",
			"LOCALSTACK_AUTH_TOKEN": os.Getenv("LOCALSTACK_AUTH_TOKEN"),
		},
    }
    localStackC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    testcontainers.CleanupContainer(t, localStackC)
    require.NoError(t, err)

	endpoint, err := localStackC.PortEndpoint(ctx, localStackPort, "http")
	if err != nil {
		t.Error(err)
	}


	repo := repositories.NewMovieRepository(
		repositories.NewRepositoryConfig(awsRegion, endpoint),
	)
	repo.Open()
	t.Run("should be able to create the tables", func(t *testing.T) {
		if err := repo.CreateTables(ctx); err != nil {
			logError(t, "Error when creating tables: %v", err)
		}
	})
	t.Run("should be able to run an entire sequence of actions with movies", func(t *testing.T) {
		test := "Should be able to get an empty list of movies"
		logTest(t, test)
		if noMovies, _, err := repo.GetAll(ctx, "", 10, 0); err != nil {
			logError(t, "Error listing movies: %v", err)
		} else if len(noMovies) != 0 {
			logError(t, "Movies list is not empty at start time.")
		} else {
			logSuccess(t, test)
		}
		

		test = "Should be able to ceate a movie"
		logTest(t, test)
		movie := domain.Movie{
			Title: faker.Sentence(),
			Year: faker.YearString(),
		}
		if err := repo.Save(ctx, movie); err != nil {
			logError(t, "Error saving movie %+v: %v", movie, err)
		} else {
			logSuccess(t, test)
		}

		test = "Should be able to get a list with movies and see the new movie"
		logTest(t, test)
		var movieId int
		var gottenMovie domain.Movie
		if oneMovie, _, err := repo.GetAll(ctx, "", 10, 0); err != nil {
			logError(t, "Error listing movies: %v", err)
		} else if len(oneMovie) != 1 {
			logError(t, "Movies list has %d movies, when it should have 1.", len(oneMovie))
		} else {
			year, title, id := oneMovie[0].Year, oneMovie[0].Title, oneMovie[0].ID
			if year != movie.Year || title != movie.Title {
				logError(t, "Movie got should have year %s and title %q, but got %s and %q",
					year, movie.Year, title, movie.Title)
			} else {
				logSuccess(t, test)
				movieId = id
				gottenMovie = oneMovie[0]
			}
		}

		test = "Should be able to get the movie by its id"
		logTest(t, test)
		if gottenMovie2, err := repo.GetOne(ctx, movieId); err != nil {
			logError(t, "Error getting movie: %v", err)
		} else if !reflect.DeepEqual(gottenMovie2, gottenMovie) {
			logError(t, "Movie got with GetOne %+v different than the one with GetAll %+v ",
			gottenMovie2, gottenMovie)
		} else {
			logSuccess(t, test)
		}

		test = "Should be able to create a new movie, and get both movies when listing"
		logTest(t, test)
		secondMovie := domain.Movie{
			Title: faker.Sentence(),
			Year: faker.YearString(),
		}
		for secondMovie.Year == movie.Year {
			secondMovie.Year = faker.YearString()
		}
		if err := repo.Save(ctx, secondMovie); err != nil {
			logError(t, "Error saving movie %+v: %v", movie, err)
		}

		if twoMovies, _, err := repo.GetAll(ctx, "", 10, 0); err != nil {
			logError(t, "Error listing movies: %v", err)
		} else if len(twoMovies) != 2 {
			logError(t, "Movies list has %d movies, when it should have 1.", len(twoMovies))
		} else {
			logSuccess(t, test)
		}

		test = "Should be able to get only one movie, if limiting to one"
		logTest(t, test)
		if listFirstMovie, cursor, err := repo.GetAll(ctx, "", 1, 0); err != nil {
			logError(t, "Error listing movies: %v", err)
		} else if len(listFirstMovie) != 1 {
			logError(t, "Should've limited to only one movie, but %d movies were fetched", len(listFirstMovie))
		} else if cursor == 0 {
			logError(t, "Should've returned a cursor for pagination reasons, but cursor was not set")
		} else {
			logSuccess(t, test)
		}

		test = "Should be able to get only first movie by its year"
		logTest(t, test)
		if listFirstMovie, _, err := repo.GetAll(ctx, movie.Year, 2, 0); err != nil {
			logError(t, "Error listing movies: %v", err)
		} else if len(listFirstMovie) != 1 {
			logError(t, "Should've filtered only first movie by its year, but %d movies were fetched", len(listFirstMovie))
		} else {
			logSuccess(t, test)
		}

		test = "Should be able to get only second movie by skipping first movie."
		logTest(t, test)
		if listSecondMovie, _, err := repo.GetAll(ctx, "", 2, gottenMovie.ID); err != nil {
			logError(t, "Error listing movies: %v", err)
		} else if len(listSecondMovie) != 1 {
			logError(t, 
				"Should've filtered only second movie by skipping the first one and all before it, but %d movies were fetched",
				len(listSecondMovie),
			)
		} else {
			logSuccess(t, test)
		}

		test = "Should be able to delete a movie, and then not be able to fetch it again"
		logTest(t, test)
		if err := repo.Delete(ctx, movieId); err != nil {
			logError(t, "Error deleting movie: %v", err)
		}
		if _, err := repo.GetOne(ctx, movieId); err == nil {
			logError(t, "Repository did not return an error after fetching a deleted movie.")
		} else if err != ports.ErrMovieNotFound {
			logError(t, "Repository did not return ErrMovieNotFound after fething a deleted movie. Returned %v instead", err)
		} else {
			logSuccess(t, test)
		}
	})
}

func logTest(t testing.TB, msg string, args ...any) {
	log := fmt.Sprintf("[TEST]: %s", msg)
	t.Logf(log, args...)
}

func logSuccess(t testing.TB, msg string, args ...any) {
	log := fmt.Sprintf("[SUCCESS]: %s", msg)
	t.Logf(log, args...)	
}

func logError(t testing.TB, msg string, args ...any) {
	log := fmt.Sprintf("[ERROR]: %s", msg)
	t.Errorf(log, args...)	
}
