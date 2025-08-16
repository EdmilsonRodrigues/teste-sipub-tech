package dtos

import (
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/domain"
)

func NewMovieID(id int) MovieID {
	return MovieID(id)
}

type MovieID int

type CreateMovieDTO struct {
	Title string `json:"title"`
	Year string `json:"year"`
}

func (dto *CreateMovieDTO) ToDomain() domain.Movie {
	return domain.Movie{
		Title: dto.Title,
		Year: dto.Year,
	}
}

func NewMovieResponseDTOFromDomain(movie domain.Movie) *MovieResponseDTO {
	dto := mapDomainToResponse(movie)
	return &dto
}

func MoviesToResponseDTOs(movies []domain.Movie) *[]MovieResponseDTO {
	dtos := make([]MovieResponseDTO, len(movies))
	for index, movie := range(movies) {
		dtos[index] = mapDomainToResponse(movie)
	}
	return &dtos
}

type MovieResponseDTO struct {
	ID MovieID `json:"id"`
	Title string `json:"title"`
	Year string `json:"year"`
}

func mapDomainToResponse(movie domain.Movie) MovieResponseDTO {
	return MovieResponseDTO{
		ID: MovieID(movie.ID),
		Title: movie.Title,
		Year: movie.Year,
	}
}

