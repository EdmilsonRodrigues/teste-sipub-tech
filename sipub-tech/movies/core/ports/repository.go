package ports

import (
	"fmt"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/domain"
)

type MovieRepository interface {
	MovieOneGetterRepository
	MovieAllGetterRepository
	MovieSaverRepository
	MovieDeleterRepository
}

var (
	MovieNotFoundError = fmt.Errorf("Movie not found in the repository.")
)

type MovieOneGetterRepository interface {
	GetOne(id int) (domain.Movie, error)
}

type MovieAllGetterRepository interface {
	GetAll() ([]domain.Movie, error)
}

type MovieSaverRepository interface {
	Save(domain.Movie) error
}

type MovieDeleterRepository interface {
	Delete(id int) error
}
