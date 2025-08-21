package ports

import (
	"context"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
)

type GetMovieCase interface {
	GetMovie(ctx context.Context, service MovieOneGetterService, id dtos.MovieId) (movie dtos.MovieResponseDTO, err error)
}

type GetMoviesCase interface {
	GetMovies(ctx context.Context, service MovieAllGetterService, query dtos.MoviesQueryDTO) (movies dtos.MoviesResponseDTO, err error)
}

type SaveMovieCase interface {
	SaveMovie(ctx context.Context, service MovieSaverService, movie dtos.CreateMovieDTO) error
}

type DeleteMovieCase interface {
	DeleteMovie(ctx context.Context, service MovieDeleterService, id dtos.MovieId) error
}


