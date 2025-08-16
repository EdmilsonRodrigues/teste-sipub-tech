package usecases_test

import (
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

			result, err := ucase.GetMovie(dtos.NewMovieID(movie.ID))
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

			_, err := ucase.GetMovie(dtos.NewMovieID(id))
			if err == nil {
				t.Logf("No error return when getting not existent movie.")
				return false
			}

			if err != ports.MovieNotFoundError {
				t.Logf("Custom error returned intead of MovieNotFoundError")
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

			_, receivedErr := ucase.GetMovie(dtos.NewMovieID(id))
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

func (repo *StubMovieOneGetter) GetOne(id int) (domain.Movie, error) {
	if id != repo.movieReturned.ID {
		return repo.movieReturned, ports.MovieNotFoundError
	}
	return repo.movieReturned, repo.errorReturned
}


func TestGetMoviesCase(t *testing.T) {
	t.Run("should return array of movie response dtos", func (t *testing.T) {
		assertion := func(movies []domain.Movie) bool {
			repo := &StubMovieAllGetter{
				moviesReturned: movies,
			}
			ucase := usecases.NewGetMoviesCase(repo)

			results, err := ucase.GetMovies()
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
		assertion := func(errorMessage string) bool {
			err := fmt.Errorf("random error: %s", errorMessage)
			repo := &StubMovieAllGetter{
				errorReturned: err,
			}
			ucase := usecases.NewGetMoviesCase(repo)

			_, receivedErr := ucase.GetMovies()
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

func (repo *StubMovieAllGetter) GetAll() ([]domain.Movie, error) {
	return repo.moviesReturned, repo.errorReturned
}


func TestSaveMovieCase(t *testing.T) {
	t.Run("should pass domain.Movie to repository", func (t *testing.T) {
		assertion := func(movie dtos.CreateMovieDTO) bool {
			repo := &MockMovieSaver{}
			ucase := usecases.NewSaveMovieCase(repo)

			if err := ucase.SaveMovie(movie); err != nil {
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

			receivedErr := ucase.SaveMovie(movie)
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

func (repo *MockMovieSaver) Save(movie domain.Movie) error {
	repo.moviePassed = movie
	return repo.errorReturned
}



func TestDeleteMovieCase(t *testing.T) {
	t.Run("should pass the id to the repository", func (t *testing.T) {
		assertion := func(id dtos.MovieID) bool {
			repo := &MockMovieDeleter{}
			ucase := usecases.NewDeleteMovieCase(repo)

			if err := ucase.DeleteMovie(id); err != nil {
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

			receivedErr := ucase.DeleteMovie(id)
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

func (repo *MockMovieDeleter) Delete(id int) error {
	repo.idPassed = id
	return repo.errorReturned
}



