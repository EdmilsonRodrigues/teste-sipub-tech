package ports

import (
	"context"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/dtos"
)

type MovieGetter interface {
	GetMovie(ctx context.Context, id dtos.MovieID) (*dtos.MovieResponseDTO, error)
}

type MoviesGetter interface {
	GetMovies(ctx context.Context, query dtos.GetMoviesDTO) (movies *[]dtos.MovieResponseDTO, newCursor int, err error)
}

type MovieSaver interface {
	SaveMovie(ctx context.Context, movie dtos.CreateMovieDTO) error
}

type MovieDeleter interface {
	DeleteMovie(ctx context.Context, id dtos.MovieID) error
}


