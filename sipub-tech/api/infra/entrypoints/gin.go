package entrypoints

import (
	"github.com/gin-gonic/gin"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/ports"
	infraPorts "github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/ports"
)


func NewGinEntrypoint(
	executorMovieService   ports.MovieExecutorService,
	queryMovieService      ports.MovieQueryService,
	movieController   infraPorts.MovieController,
) *GinEntrypoint {
	return &GinEntrypoint{
		executorMovieService: executorMovieService,
		queryMovieService: queryMovieService,

		movieController: movieController,
	}
}

type GinEntrypoint struct {
	executorMovieService       ports.MovieExecutorService
	queryMovieService          ports.MovieQueryService

	movieController     infraPorts.MovieController

	engine                     *gin.Engine
}

func (entrypoint *GinEntrypoint) Setup() {
	router := gin.Default()

	entrypoint.engine = router

	entrypoint.addMovieHandlers()
}

func (entrypoint *GinEntrypoint) Serve() {
	if err := entrypoint.engine.Run(); err != nil { // listen and serve on 0.0.0.0:8080
		panic("gin engine failed to run")
	}
}

func (entrypoint *GinEntrypoint) GetEngine() *gin.Engine {
	return entrypoint.engine
}
