package ports

import (
	"context"
	"fmt"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/domain"
)

type MovieRepository interface {
	TableCreatorRepository
	MovieOneGetterRepository
	MovieAllGetterRepository
	MovieSaverRepository
	MovieDeleterRepository
}

var (
	ErrMovieNotFound = fmt.Errorf("movie not found in the repository")
)

type TableCreatorRepository interface {
	CreateTables(ctx context.Context) error
}

type MovieOneGetterRepository interface {
	GetOne(ctx context.Context, id int) (movie domain.Movie, err error)
}

type MovieAllGetterRepository interface {
	GetAll(ctx context.Context, year string, limit int, lastMovieId int) (movies []domain.Movie, cursor int, err error)
}

type MovieSaverRepository interface {
	Save(ctx context.Context, movie domain.Movie) error
}

type MovieDeleterRepository interface {
	Delete(ctx context.Context, id int) error
}
