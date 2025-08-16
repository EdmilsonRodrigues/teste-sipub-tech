package dtos_test

import (
	"reflect"
	"testing"
	"testing/quick"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/domain"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/dtos"
)


func TestMoviesToResponseDTOs(t *testing.T) {
	t.Run("should parse correctly an array with 10 random movies.", func(t *testing.T) {
		assertion := func(ids [10]int, titles, years [10]string) bool {
			movies := make([]domain.Movie, 10)
			expected := make([]dtos.MovieResponseDTO, 10)

			for index := range(10) {
				movies[index] = domain.Movie{ID: ids[index], Title: titles[index], Year: years[index]}
				expected[index] = dtos.MovieResponseDTO{ID: dtos.MovieID(ids[index]), Title: titles[index], Year: years[index]}
			}

			result := *dtos.MoviesToResponseDTOs(movies)
			if !reflect.DeepEqual(expected, result) {
				t.Logf("Expected: %+v different than Result: %+v", expected, result)
				return false
			}
			return true
		}
		if err := quick.Check(assertion, nil); err != nil {
			t.Errorf("Error found testing assertion: %v", err)
		}
	})
}
