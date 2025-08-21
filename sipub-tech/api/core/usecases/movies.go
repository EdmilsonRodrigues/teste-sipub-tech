package usecases

import (
	"context"
	"fmt"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/ports"
)

func NewGetMovieCase() *GetMovieCase {
	return &GetMovieCase{}
}

type GetMovieCase struct {}

func (ucase *GetMovieCase) GetMovie(
	ctx context.Context, service ports.MovieOneGetterService, id dtos.MovieId,
) (movie dtos.MovieResponseDTO, err error) {
	if movie, err = service.GetOne(ctx, id); err != nil {
		if err == ports.ErrMovieNotFound {
			return movie, err
		}
		return movie, fmt.Errorf("could not get movie with id %d: %w", id, err)
	}
	return
}

func NewGetMoviesCase() *GetMoviesCase {
	return &GetMoviesCase{}
}

type GetMoviesCase struct {}

func (ucase *GetMoviesCase) GetMovies(
	ctx context.Context, service ports.MovieAllGetterService, query dtos.MoviesQueryDTO,
) (movies dtos.MoviesResponseDTO, err error) {
	if movies, err = service.GetAll(ctx, query); err != nil {
		return movies, fmt.Errorf("could not get movies with query %+v: %w", query, err)
	}
	return
}

func NewSaveMovieCase() *SaveMovieCase {
	return &SaveMovieCase{}
}

type SaveMovieCase struct {}

func (ucase *SaveMovieCase) SaveMovie(ctx context.Context, service ports.MovieSaverService, movie dtos.CreateMovieDTO) error {
	if err := service.Save(ctx, movie); err != nil {
		return fmt.Errorf("could not save movie %+v: %w", movie, err)
	}
	return nil
}

func NewDeleteMovieCase() *DeleteMovieCase {
	return &DeleteMovieCase{}
}

type DeleteMovieCase struct {}

func (ucase *DeleteMovieCase) DeleteMovie(ctx context.Context, service ports.MovieDeleterService, id dtos.MovieId) error {
	if err := service.Delete(ctx, id); err != nil {
		return fmt.Errorf("could not delete movie with id %d", id)
	}
	return nil
}
