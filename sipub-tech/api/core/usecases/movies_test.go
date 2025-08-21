package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"testing/quick"
	
	"github.com/stretchr/testify/assert"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/ports"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/usecases"
)


func TestGetMovieCase(t *testing.T) {
	usecase := usecases.NewGetMovieCase()
	t.Run("should return movie from service if no error occurs", func(t *testing.T) {
		assertion := func(movie dtos.MovieResponseDTO, id dtos.MovieId) bool {
			service := &MockMovieOneGetterService{
				MovieToReturn: movie,
			}
			response, err := usecase.GetMovie(context.Background(), service, id)
			if err != nil {
				t.Logf("Error returned when should not have error: %v", err)
			}

			if !assert.Equal(t, movie, response) {
				return false
			}
			if !assert.Equal(t, id, service.IdPassed) {
				return false
			}
			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("assertion failed: %v", err)
		}		
	})

	t.Run("should return ports.ErrMovieNotFound if service returns ports.ErrMovieNotFound", func(t *testing.T) {
		assertion := func(id dtos.MovieId) bool {
			service := &MockMovieOneGetterService{
				ReturnedError: ports.ErrMovieNotFound,
			}
			_, err := usecase.GetMovie(context.Background(), service, id)
			if err == nil {
				t.Logf("Error not returned when should have error")
				return false
			}

			if !assert.Equal(t, err, ports.ErrMovieNotFound) {
				return false
			}

			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("assertion failed: %v", err)
		}
	})

	t.Run("should return custom error if service returns another error", func(t *testing.T) {
		assertion := func(id dtos.MovieId, msg string) bool {
			err := fmt.Errorf("error found: %s", msg)
			service := &MockMovieOneGetterService{
				ReturnedError: err,
			}
			_, returnedErr := usecase.GetMovie(context.Background(), service, id)
			if returnedErr == nil {
				t.Logf("Error not returned when should have error")
				return false
			}

			if !assert.NotEqual(t, err, returnedErr) {
				return false
			}

			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("assertion failed: %v", err)
		}
	})
}


type MockMovieOneGetterService struct {
	ports.MovieOneGetterService

	MovieToReturn dtos.MovieResponseDTO 
	ReturnedError error
	IdPassed      dtos.MovieId
}

func (svc *MockMovieOneGetterService) GetOne(ctx context.Context, id dtos.MovieId) (dtos.MovieResponseDTO, error) {
	svc.IdPassed = id
	return svc.MovieToReturn, svc.ReturnedError
}


func TestGetMoviesCase(t *testing.T) {
	usecase := usecases.NewGetMoviesCase()
	t.Run("should return movies returned by service if no error ocurred.", func(t *testing.T) {
		assertion := func(movies dtos.MoviesResponseDTO, query dtos.MoviesQueryDTO) bool {
			service := &MockMovieAllGetterService{
				MoviesToReturn: movies,
			}
			response, err := usecase.GetMovies(context.Background(), service, query)
			if err != nil {
				t.Logf("Error returned when should not have error: %v", err)
			}

			if !assert.Equal(t, movies, response) {
				return false
			}
			if !assert.Equal(t, query, service.QueryPassed) {
				return false
			}
			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("assertion failed: %v", err)
		}		
	})

	t.Run("should return custom error if error was returned.", func(t *testing.T) {
		assertion := func(query dtos.MoviesQueryDTO, msg string) bool {
			err := fmt.Errorf("error found: %s", msg)
			service := &MockMovieAllGetterService{
				ReturnedError: err,
			}
			_, returnedErr := usecase.GetMovies(context.Background(), service, query)
			if returnedErr == nil {
				t.Logf("Error not returned when should have error")
				return false
			}

			if !assert.NotEqual(t, err, returnedErr) {
				return false
			}

			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("assertion failed: %v", err)
		}
	})
}


type MockMovieAllGetterService struct {
	ports.MovieAllGetterService

	MoviesToReturn  dtos.MoviesResponseDTO
	ReturnedError   error
	QueryPassed     dtos.MoviesQueryDTO
}

func (svc *MockMovieAllGetterService) GetAll(ctx context.Context, query dtos.MoviesQueryDTO) (dtos.MoviesResponseDTO, error) {
	svc.QueryPassed = query
	return svc.MoviesToReturn, svc.ReturnedError
}

func TestSaveMovieCase(t *testing.T) {
	usecase := usecases.NewSaveMovieCase()
	t.Run("should pass movie to service when called.", func(t *testing.T) {
		assertion := func(movie dtos.CreateMovieDTO) bool {
			service := &MockMovieSaverService{}
			err := usecase.SaveMovie(context.Background(), service, movie)
			if err != nil {
				t.Logf("Error returned when should not have error: %v", err)
			}

			if !assert.Equal(t, movie, service.MoviePassed) {
				return false
			}
			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("assertion failed: %v", err)
		}
	})

	t.Run("should return custom error if error was returned by service", func(t *testing.T) {
		assertion := func(movie dtos.CreateMovieDTO, msg string) bool {
			err := fmt.Errorf("error found: %s", msg)
			service := &MockMovieSaverService{
				ReturnedError: err,
			}
			returnedErr := usecase.SaveMovie(context.Background(), service, movie)
			if returnedErr == nil {
				t.Logf("Error not returned when should have error")
				return false
			}

			if !assert.NotEqual(t, err, returnedErr) {
				return false
			}

			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("assertion failed: %v", err)
		}
	})
}


type MockMovieSaverService struct {
	ports.MovieSaverService

	ReturnedError error
	MoviePassed   dtos.CreateMovieDTO
}

func (svc *MockMovieSaverService) Save(ctx context.Context, movie dtos.CreateMovieDTO) error {
	svc.MoviePassed = movie
	return svc.ReturnedError
}


func TestDeleteMovieCase(t *testing.T) {
	usecase := usecases.NewDeleteMovieCase()
	t.Run("should pass movie id to service when called.", func(t *testing.T) {
		assertion := func(id dtos.MovieId) bool {
			service := &MockMovieDeleterService{}
			err := usecase.DeleteMovie(context.Background(), service, id)
			if err != nil {
				t.Logf("Error returned when should not have error: %v", err)
			}

			if !assert.Equal(t, id, service.IdPassed) {
				return false
			}
			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("assertion failed: %v", err)
		}
	})

	t.Run("should return custom error if error was returned by service", func(t *testing.T) {
		assertion := func(id dtos.MovieId, msg string) bool {
			err := fmt.Errorf("error found: %s", msg)
			service := &MockMovieDeleterService{
				ReturnedError: err,
			}
			returnedErr := usecase.DeleteMovie(context.Background(), service, id)
			if returnedErr == nil {
				t.Logf("Error not returned when should have error")
				return false
			}

			if !assert.NotEqual(t, err, returnedErr) {
				return false
			}

			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("assertion failed: %v", err)
		}
	})	
}


type MockMovieDeleterService struct {
	ports.MovieDeleterService

	ReturnedError error
	IdPassed      dtos.MovieId
}

func (svc *MockMovieDeleterService) Delete(ctx context.Context, id dtos.MovieId) error {
	svc.IdPassed = id
	return svc.ReturnedError
}

