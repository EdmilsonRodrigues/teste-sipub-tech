package usecases

import (
	"fmt"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/ports"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/dtos"
)

func NewGetMovieCase(repo ports.MovieOneGetterRepository) *GetMovieCase {
	return &GetMovieCase{
		repo: repo,
	}
}

type GetMovieCase struct {
	repo ports.MovieOneGetterRepository
}

func (ucase *GetMovieCase) GetMovie(id dtos.MovieID) (*dtos.MovieResponseDTO, error) {
	movie, err := ucase.repo.GetOne(int(id))
	if err != nil {
		if err == ports.ErrMovieNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("error getting movie: %w", err)
	}
	return dtos.NewMovieResponseDTOFromDomain(movie), nil
}


func NewGetMoviesCase(repo ports.MovieAllGetterRepository) *GetMoviesCase {
	return &GetMoviesCase{
		repo: repo,
	}
}

type GetMoviesCase struct {
	repo ports.MovieAllGetterRepository
}

func (ucase *GetMoviesCase) GetMovies() (*[]dtos.MovieResponseDTO, error) {
	movies, err := ucase.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting movies %w", err)
	}
	return dtos.MoviesToResponseDTOs(movies), nil
}

func NewSaveMovieCase(repo ports.MovieSaverRepository) *SaveMovieCase {
	return &SaveMovieCase{
		repo: repo,
	}
}

type SaveMovieCase struct {
	repo ports.MovieSaverRepository
}

func (ucase *SaveMovieCase) SaveMovie(movie dtos.CreateMovieDTO) error {
	if err := ucase.repo.Save(movie.ToDomain()); err != nil {
		return fmt.Errorf("error saving movie %w", err)
	}
	return nil
}

func NewDeleteMovieCase(repo ports.MovieDeleterRepository) *DeleteMovieCase {
	return &DeleteMovieCase{
		repo: repo,
	}
}

type DeleteMovieCase struct {
	repo ports.MovieDeleterRepository
}

func (ucase *DeleteMovieCase) DeleteMovie(id dtos.MovieID) error {
	if err := ucase.repo.Delete(int(id)); err != nil {
		return fmt.Errorf("error deleting movie %w", err)
	}
	return nil
}

