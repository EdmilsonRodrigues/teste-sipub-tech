package ports

import (
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/dtos"
)

type MovieGetter interface {
	GetMovie(id dtos.MovieID) (*dtos.MovieResponseDTO, error)
}

type MoviesGetter interface {
	GetMovies() (*[]dtos.MovieResponseDTO, error)
}

type MovieSaver interface {
	SaveMovie(movie dtos.CreateMovieDTO) error
}

type MovieDeleter interface {
	DeleteMovie(id dtos.MovieID) error
}


