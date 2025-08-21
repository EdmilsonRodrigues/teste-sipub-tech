package dtos

type DataItem map[string]any

type MovieId int

type CreateMovieDTO struct {
	Title string `json:"title"`
	Year string  `json:"year"`
}


type MovieResponseDTO struct {
	ID    int     `json:"id"`
	Title string  `json:"title"`
	Year  string  `json:"year"`	
}

func (dto *MovieResponseDTO) ToDataItem() DataItem {
	return DataItem{
		"id":   dto.ID,
		"title": dto.Title,
		"year":  dto.Year,
	}
}


type MoviesQueryDTO struct {
	Year    string  `json:"year"`
	Cursor  int     `json:"id"`
	Limit   int     `json:"limit"`
}

type MoviesResponseDTO struct {
	Movies  []*MovieResponseDTO  `json:"movies"`
	Cursor  int                 `json:"cursor"`
}

