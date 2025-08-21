package ports

import (
	"context"
	"fmt"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
)

const (
	ServiceKey = "service"
)

var (
	ErrMovieNotFound = fmt.Errorf("movie not found")
)

type MovieQueryService interface {
	MovieOneGetterService
	MovieAllGetterService
}

type MovieExecutorService interface {
	MovieSaverService
	MovieDeleterService
}

type MovieOneGetterService interface {
	GetOne(ctx context.Context, id dtos.MovieId) (dtos.MovieResponseDTO, error)
}

type MovieAllGetterService interface {
	GetAll(ctx context.Context, query dtos.MoviesQueryDTO) (dtos.MoviesResponseDTO, error)
}

type MovieSaverService interface {
	Save(ctx context.Context, movie dtos.CreateMovieDTO) error
}

type MovieDeleterService interface {
	Delete(ctx context.Context, id dtos.MovieId) error
}

