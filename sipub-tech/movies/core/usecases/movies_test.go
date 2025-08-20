package usecases_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/domain"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/ports"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/usecases"
)


func TestGetMovieCase(t *testing.T) {
	t.Run("should return movie response dto when getting existing movie", func (t *testing.T) {
		assertion := func(movie domain.Movie) bool {

			repo := &StubMovieOneGetter{
				movieReturned: movie,
			}
			ucase := usecases.NewGetMovieCase(repo)

			result, err := ucase.GetMovie(context.Background(), dtos.NewMovieID(movie.ID))
			if err != nil {
				t.Logf("Error found when getting movie %v", err)
				return false
			}

			expected := dtos.MovieResponseDTO{
				ID: dtos.MovieID(movie.ID),
				Title: movie.Title,
				Year: movie.Year,
			}

			if !reflect.DeepEqual(expected, *result) {
				t.Logf("Expected: %+v different than Result: %+v", expected, *result)
				return false
			}
			
			
			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("Failed assertion: %v", err)
		}
	})

	t.Run("should return NotFoundError when getting a non existing movie", func (t *testing.T) {
		assertion := func(id int) bool {
			repo := &StubMovieOneGetter{}
			ucase := usecases.NewGetMovieCase(repo)

			_, err := ucase.GetMovie(context.Background(), dtos.NewMovieID(id))
			if err == nil {
				t.Logf("No error return when getting not existent movie.")
				return false
			}

			if err != ports.ErrMovieNotFound {
				t.Logf("Custom error returned intead of ErrMovieNotFound")
				return false
			}
			
			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("Failed assertion: %v", err)
		}
	})

	t.Run("should return custom error when receiving an error from the repository", func (t *testing.T) {
		assertion := func(id int, errorMessage string) bool {
			err := fmt.Errorf("random error: %s", errorMessage)
			repo := &StubMovieOneGetter{
				movieReturned: domain.Movie{ID: id},
				errorReturned: err,
			}
			ucase := usecases.NewGetMovieCase(repo)

			_, receivedErr := ucase.GetMovie(context.Background(), dtos.NewMovieID(id))
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
			t.Errorf("Failed assertion: %v", err)
		}
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


func TestGetMoviesCase(t *testing.T) {
	t.Run("should return array of movie response dtos", func (t *testing.T) {
		assertion := func(movies []domain.Movie, query dtos.GetMoviesDTO) bool {
			repo := &StubMovieAllGetter{
				moviesReturned: movies,
			}
			ucase := usecases.NewGetMoviesCase(repo)

			results, _, err := ucase.GetMovies(context.Background(), query)
			if err != nil {
				t.Logf("Error found when getting movies %v", err)
				return false
			}

			expected := make([]dtos.MovieResponseDTO, len(movies))
			for index, movie := range(movies) {
				expected[index] = dtos.MovieResponseDTO{
					ID: dtos.MovieID(movie.ID),
					Title: movie.Title,
					Year: movie.Year,	
				}
			}

			if !reflect.DeepEqual(expected, *results) {
				t.Logf("Results: %+v different from Expected: %+v", *results, expected)
				return false
			}

			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("Failed assertion: %v", err)
		}
	})

	t.Run("should return custom error when receiving an error from the repository", func (t *testing.T) {
		assertion := func(errorMessage string, query dtos.GetMoviesDTO) bool {
			err := fmt.Errorf("random error: %s", errorMessage)
			repo := &StubMovieAllGetter{
				errorReturned: err,
			}
			ucase := usecases.NewGetMoviesCase(repo)

			_, _, receivedErr := ucase.GetMovies(context.Background(), query)
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
			t.Errorf("Failed assertion: %v", err)
		}
	})
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


func TestSaveMovieCase(t *testing.T) {
	t.Run("should pass domain.Movie to repository", func (t *testing.T) {
		assertion := func(movie dtos.CreateMovieDTO) bool {
			repo := &MockMovieSaver{}
			ucase := usecases.NewSaveMovieCase(repo)

			if err := ucase.SaveMovie(context.Background(), movie); err != nil {
				t.Logf("Error found when saving movie %v", err)
				return false
			}

			expected := domain.Movie{
				Title: movie.Title,
				Year: movie.Year,
			}

			result := repo.moviePassed

			if !reflect.DeepEqual(expected, result) {
				t.Logf("Movie passed: %+v different from Expected: %+v", result, expected)
				return false
			}

			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("Failed assertion: %v", err)
		}
	})

	t.Run("should return custom error when receiving an error from the repository", func (t *testing.T) {
		assertion := func(errorMessage string, movie dtos.CreateMovieDTO) bool {
			err := fmt.Errorf("random error: %s", errorMessage)
			repo := &MockMovieSaver{
				errorReturned: err,
			}
			ucase := usecases.NewSaveMovieCase(repo)

			receivedErr := ucase.SaveMovie(context.Background(), movie)
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
			t.Errorf("Failed assertion: %v", err)
		}
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



func TestDeleteMovieCase(t *testing.T) {
	t.Run("should pass the id to the repository", func (t *testing.T) {
		assertion := func(id dtos.MovieID) bool {
			repo := &MockMovieDeleter{}
			ucase := usecases.NewDeleteMovieCase(repo)

			if err := ucase.DeleteMovie(context.Background(), id); err != nil {
				t.Logf("Error found when deleting movie %v", err)
				return false
			}

			expected := int(id)

			result := repo.idPassed

			if expected != result {
				t.Logf("Movie passed: %+v different from Expected: %+v", result, expected)
				return false
			}

			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("Failed assertion: %v", err)
		}
	})

	t.Run("should return custom error when receiving an error from the repository", func (t *testing.T) {
		assertion := func(errorMessage string, id dtos.MovieID) bool {
			err := fmt.Errorf("random error: %s", errorMessage)
			repo := &MockMovieDeleter{
				errorReturned: err,
			}
			ucase := usecases.NewDeleteMovieCase(repo)

			receivedErr := ucase.DeleteMovie(context.Background(), id)
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
			t.Errorf("Failed assertion: %v", err)
		}
	})
}


type MockMovieDeleter struct {
	idPassed int
	errorReturned error
}

func (repo *MockMovieDeleter) Delete(ctx context.Context, id int) error {
	repo.idPassed = id
	return repo.errorReturned
}



