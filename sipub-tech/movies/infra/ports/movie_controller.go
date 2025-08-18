package ports

import (
	"context"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/dtos"
)

type MovieSaverController interface {
	SaveMovie(ctx context.Context, movie dtos.CreateMovieDTO)
}

type MovieDeleterController interface {
	DeleteMovie(ctx context.Context, id dtos.MovieID)
}


