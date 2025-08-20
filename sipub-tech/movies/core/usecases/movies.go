package usecases

import (
	"context"
	"fmt"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/ports"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/domain"
)

func NewGetMovieCase(repo ports.MovieOneGetterRepository) *GetMovieCase {
	return &GetMovieCase{
		repo: repo,
	}
}

type GetMovieCase struct {
	repo ports.MovieOneGetterRepository
}

func (ucase *GetMovieCase) GetMovie(ctx context.Context, id dtos.MovieID) (*dtos.MovieResponseDTO, error) {
	movie, err := ucase.repo.GetOne(ctx, int(id))
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

func (ucase *GetMoviesCase) GetMovies(
	ctx context.Context, query dtos.GetMoviesDTO,
) (movies *[]dtos.MovieResponseDTO, newCursor int, err error) {
	var fetchedMovies []domain.Movie

	fetchedMovies, newCursor, err = ucase.repo.GetAll(ctx, query.Year, query.Limit, query.Cursor)
	if err != nil {
		err = fmt.Errorf("error getting movies %w", err)
		return
	}

	movies = dtos.MoviesToResponseDTOs(fetchedMovies)
	return
}

func NewSaveMovieCase(repo ports.MovieSaverRepository) *SaveMovieCase {
	return &SaveMovieCase{
		repo: repo,
	}
}

type SaveMovieCase struct {
	repo ports.MovieSaverRepository
}

func (ucase *SaveMovieCase) SaveMovie(ctx context.Context, movie dtos.CreateMovieDTO) error {
	if err := ucase.repo.Save(ctx, movie.ToDomain()); err != nil {
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

func (ucase *DeleteMovieCase) DeleteMovie(ctx context.Context, id dtos.MovieID) error {
	if err := ucase.repo.Delete(ctx, int(id)); err != nil {
		return fmt.Errorf("error deleting movie %w", err)
	}
	return nil
}

