package controllers_test

import (
	"context"
	"fmt"
	"testing"
	"testing/quick"

	pb "github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/grpc/movies"
	pb_exceptions "github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/grpc/exceptions"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/infra/controllers"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/domain"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/ports"
)

func TestGRPCMovieController(t *testing.T) {
	controller := controllers.GRPCMovieController{}
	ctx := context.Background()

	const MaxId = int32(^uint32(0) >> 1)
	const MinId = 1
	t.Run("when executing GetMovie", func (t *testing.T) {
		t.Run("should return movie pb.Movie when getting existing movie", func (t *testing.T) {
			assertion := func(movie domain.Movie) bool {
				if movie.ID < MinId || movie.ID > int(MaxId) {
					return true
				}
				repo := &StubMovieOneGetter{
					movieReturned: movie,
				}
				ctx := context.WithValue(ctx, controllers.RepoKey, repo)

				result, err := controller.GetMovie(ctx, &pb.GetMovieRequest{Id: int32(movie.ID)})
				if err != nil {
					t.Logf("Error found when getting movie %v", err)
					return false
				}

				expected := pb.Movie{
					Id: int32(movie.ID),
					Title: movie.Title,
					Year: movie.Year,
				}
				
				if expected.String() != result.String() {
					t.Logf("Expected: %q different than Result: %q", expected.String(), result.String())
					return false
				}				
				
				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("failed assertion: %v", err)
			}
		})

		t.Run("should return pb_exceptions.ErrMovieNotFound when getting a non existing movie", func (t *testing.T) {
			assertion := func(id dtos.MovieID) bool {
				repo := &StubMovieOneGetter{}
				ctx := context.WithValue(ctx, controllers.RepoKey, repo)

				_, err := controller.GetMovie(ctx, &pb.GetMovieRequest{Id: int32(id)})

				if err == nil {
					t.Logf("No error return when getting not existent movie.")
					return false
				}

				if err != pb_exceptions.ErrMovieNotFound {
					t.Logf("Custom error returned intead of ErrMovieNotFound")
					return false
				}
				
				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("failed assertion: %v", err)
			}
		})

		t.Run("should return custom error when receiving an error from the repository", func (t *testing.T) {
			assertion := func(id int, errorMessage string) bool {
				err := fmt.Errorf("random error: %s", errorMessage)
				repo := &StubMovieOneGetter{
					movieReturned: domain.Movie{ID: id},
					errorReturned: err,
				}
				ctx := context.WithValue(ctx, controllers.RepoKey, repo)

				_, receivedErr := controller.GetMovie(ctx, &pb.GetMovieRequest{Id: int32(id)})
				if receivedErr == nil {
					t.Logf("No error return when getting not existent movie.")
					return false
				}

				if receivedErr == err {
					t.Logf("Custom error not returned, same error returned instead.")
					return false
				}
				
				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("failed assertion: %v", err)
			}
		})

		t.Run("should return error if repository not set in context.", func(t *testing.T) {
			assertion := func(id dtos.MovieID) bool {
				_, err := controller.GetMovie(ctx, &pb.GetMovieRequest{Id: int32(id)})

				if err == nil {
					t.Logf("No error return when checking for repository.")
					return false
				}

				if err != controllers.ErrUnsetRespository {
					t.Logf("Did not return ErrUnsetRepository when repository was unset.")
					return false
				}
				
				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("failed assertion: %v", err)
			}

		})

	})

	t.Run("when executing GetMovies", func(t *testing.T) {		
		t.Run("should return array of movie response dtos", func (t *testing.T) {
			assertion := func(movies []domain.Movie) bool {
				repo := &StubMovieAllGetter{
					moviesReturned: movies,
				}
				ctx := context.WithValue(ctx, controllers.RepoKey, repo)

				results, err := controller.GetMovies(ctx, &pb.GetMoviesRequest{})
				if err != nil {
					t.Logf("Error found when getting movies %v", err)
					return false
				}

				
				for index, movie := range(movies) {
					expected := &pb.Movie{
						Id: int32(movie.ID),
						Title: movie.Title,
						Year: movie.Year,	
					}
					
					if expected.String() != (*results.Movies[index]).String() {
						t.Logf("Expected %q but got %q", expected.String(), (*results.Movies[index]).String())
						return false
					}
				}

				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("failed assertion: %v", err)
			}
		})

		t.Run("should return custom error when receiving an error from the repository", func (t *testing.T) {
			assertion := func(id int, errorMessage string) bool {
				err := fmt.Errorf("random error: %s", errorMessage)
				repo := &StubMovieAllGetter{
					errorReturned: err,
				}
				ctx := context.WithValue(ctx, controllers.RepoKey, repo)

				_, receivedErr := controller.GetMovies(ctx, &pb.GetMoviesRequest{})
				if receivedErr == nil {
					t.Logf("No error return when getting not existent movie.")
					return false
				}

				if receivedErr == err {
					t.Logf("Custom error not returned, same error returned instead.")
					return false
				}
				
				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("failed assertion: %v", err)
			}
		})

		t.Run("should return error if repository not set in context.", func(t *testing.T) {
			assertion := func(id dtos.MovieID) bool {
				_, err := controller.GetMovies(ctx, &pb.GetMoviesRequest{})

				if err == nil {
					t.Logf("No error return when checking for repository.")
					return false
				}

				if err != controllers.ErrUnsetRespository {
					t.Logf("Did not return ErrUnsetRepository when repository was unset.")
					return false
				}
				
				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("failed assertion: %v", err)
			}

		})
	})
}


type StubMovieOneGetter struct {
	movieReturned domain.Movie
	errorReturned error
}

func (repo *StubMovieOneGetter) GetOne(ctx context.Context, id int) (movie domain.Movie, err error) {
	if id != repo.movieReturned.ID {
		return repo.movieReturned, ports.ErrMovieNotFound
	}
	return repo.movieReturned, repo.errorReturned
}


type StubMovieAllGetter struct {
	moviesReturned []domain.Movie
	errorReturned error
}

func (repo *StubMovieAllGetter) GetAll(
	ctx context.Context, year string, limit int, lastMovieId int,
) (movies []domain.Movie, cursor int, err error) {
	return repo.moviesReturned, 0, repo.errorReturned
}


func TestMessagingController(t *testing.T) {
	controller := &controllers.MessagingMovieController{}
	ctx := context.Background()

	t.Run("when executing SaveMovie", func(t *testing.T) {
		t.Run("should pass the movie to the usecase and return its error", func(t *testing.T) {
			assertion := func(movie dtos.CreateMovieDTO, errString string) bool {
				var err error
				if errString != "" {
					err = fmt.Errorf("an error: %s", errString)
				}

				repo := &MockMovieSaver{errorReturned: err}
				ctx := context.WithValue(ctx, controllers.RepoKey, repo)

				resultErr := controller.SaveMovie(ctx, movie)

				if  (err == nil) != (resultErr == nil) {
					t.Logf("Expected %v error found %v", err, resultErr)
					return false
				}
				
				return true
				
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("failed assertion: %v", err)
			}
		})
		
		t.Run("should return error if repository not set in context.", func(t *testing.T) {
			assertion := func(movie dtos.CreateMovieDTO) bool {
				err := controller.SaveMovie(ctx, movie)

				if err == nil {
					t.Logf("No error return when checking for repository.")
					return false
				}

				if err != controllers.ErrUnsetRespository {
					t.Logf("Did not return ErrUnsetRepository when repository was unset.")
					return false
				}
				
				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("failed assertion: %v", err)
			}

		})
	})

	t.Run("when executing DeleteMovie", func(t *testing.T) {
		t.Run("should pass the movie id to the usecase and return its error", func(t *testing.T) {
			assertion := func(id dtos.MovieID, errString string) bool {
				var err error
				if errString != "" {
					err = fmt.Errorf("an error: %s", errString)
				}

				repo := &MockMovieDeleter{errorReturned: err}
				ctx := context.WithValue(ctx, controllers.RepoKey, repo)

				resultErr := controller.DeleteMovie(ctx, id)

				if  (err == nil) != (resultErr == nil) {
					t.Logf("Expected %v error found %v", err, resultErr)
					return false
				}
				
				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("failed assertion: %v", err)
			}
		})

		t.Run("should return error if repository not set in context.", func(t *testing.T) {
			assertion := func(id dtos.MovieID) bool {
				err := controller.DeleteMovie(ctx, id)

				if err == nil {
					t.Logf("No error return when checking for repository.")
					return false
				}

				if err != controllers.ErrUnsetRespository {
					t.Logf("Did not return ErrUnsetRepository when repository was unset.")
					return false
				}
				
				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("failed assertion: %v", err)
			}

		})
	})
}

type MockMovieSaver struct {
	moviePassed domain.Movie
	errorReturned error
}

func (repo *MockMovieSaver) Save(ctx context.Context, movie domain.Movie) error {
	repo.moviePassed = movie
	return repo.errorReturned
}


type MockMovieDeleter struct {
	idPassed int
	errorReturned error
}

func (repo *MockMovieDeleter) Delete(ctx context.Context, id int) error {
	repo.idPassed = id
	return repo.errorReturned
}


