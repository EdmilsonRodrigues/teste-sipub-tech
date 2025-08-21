package entrypoints

import (
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/usecases"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/middlewares"
)

func (entrypoint *GinEntrypoint) addMovieHandlers() {
	queryGroup := entrypoint.engine.Group(
		"/movies",
		middlewares.AddMovieQueryService(entrypoint.queryMovieService),
	)
	queryGroup.GET(
		"/",
		middlewares.ParseQueryParameters(),
		entrypoint.movieController.GetMoviesHandler(usecases.NewGetMoviesCase()),
	)
	queryGroup.GET(
		"/:id",
		entrypoint.movieController.GetMovieHandler(usecases.NewGetMovieCase()),
	)

	executorGroup := entrypoint.engine.Group(
		"/movies",
		middlewares.AddMovieExecutorService(entrypoint.executorMovieService),
	)

	executorGroup.POST(
		"/",
		entrypoint.movieController.SaveMovieHandler(usecases.NewSaveMovieCase()),
	)
	executorGroup.DELETE(
		"/:id",
		entrypoint.movieController.DeleteMovieHandler(usecases.NewDeleteMovieCase()),
	)
	
}
