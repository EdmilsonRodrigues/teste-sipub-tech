package ports

import (
	"github.com/gin-gonic/gin"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/ports"
)

type MovieController interface {
	GetMovieHandler(usecase ports.GetMovieCase) gin.HandlerFunc
	GetMoviesHandler(usecase ports.GetMoviesCase) gin.HandlerFunc
	SaveMovieHandler(usecase ports.SaveMovieCase) gin.HandlerFunc
	DeleteMovieHandler(usecase ports.DeleteMovieCase) gin.HandlerFunc
}


