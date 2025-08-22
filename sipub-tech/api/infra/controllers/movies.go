package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	
	"github.com/gin-gonic/gin"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/ports"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/errors"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/middlewares"
	infraPorts "github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/ports"
	infraDtos "github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/dtos"
)

func NewMovieController() *MovieController {
	return &MovieController{}
}

type MovieController struct {
	infraPorts.MovieController
}


func (controller *MovieController) GetMovieHandler(usecase ports.GetMovieCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, ok := controller.getService(ctx)
		if !ok {
			return
		}

		svc, ok := service.(ports.MovieOneGetterService)
		if !ok {
			controller.internalServerError(ctx, "service unavailable.", "Service malformed.")
			return
		}

		id, ok := controller.getId(ctx)
		if !ok {
			return
		}

		movie, err := usecase.GetMovie(ctx, svc, dtos.MovieId(id))
		if err != nil {
			if strings.Contains(err.Error(), ports.ErrMovieNotFound.Error())  {
				ctx.JSON(http.StatusNotFound, errors.MovieNotFoundErrorResponse)
				ctx.Abort()
				return
			}
			controller.internalServerError(ctx, "path broken.", fmt.Sprintf("Failed to fetch movie with id %d: %v", id, err))
			return
		}

		ctx.JSON(http.StatusOK, infraDtos.NewJSONResponse(&movie))
	}
}


func (controller *MovieController) GetMoviesHandler(usecase ports.GetMoviesCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exists := controller.getService(ctx)
		if !exists {
			return
		}

		svc, ok := service.(ports.MovieAllGetterService)
		if !ok {
			controller.internalServerError(ctx, "service unavailable.", "Service malformed.")
			return
		}

		dto, exists := ctx.Get(middlewares.DtoKey)
		if !exists {
			controller.internalServerError(ctx, "path broken.", "DTO Parser middleware did not set context.")
			return
		}

		query, ok := dto.(*dtos.MoviesQueryDTO)
		if !ok {
			controller.internalServerError(ctx, "path broken.", "DTO Parser middleware set malformed context.")
			return
		}

		movies, err := usecase.GetMovies(ctx, svc, *query)
		if err != nil {
			controller.internalServerError(ctx, "path broken.", fmt.Sprintf("Failed to fetch movies with query %+v: %v", query, err))
			return
		}

		dtos := make([]infraDtos.JsonData, len(movies.Movies))
		for index, movie := range movies.Movies {
			dtos[index] = movie
		} 

		ctx.JSON(http.StatusOK, infraDtos.NewPaginatedResponse(dtos, query.Limit, movies.Cursor))
	}
}


func (controller *MovieController) SaveMovieHandler(usecase ports.SaveMovieCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exists := controller.getService(ctx)
		if !exists {
			return
		}

		svc, ok := service.(ports.MovieSaverService)
		if !ok {
			controller.internalServerError(ctx, "service unavailable.", "Service malformed.")
			return
		}

		var dto dtos.CreateMovieDTO

		if err := ctx.ShouldBind(&dto); err != nil {
			controller.unprocessableEntityError(ctx, "body malformed.", fmt.Sprintf("Body could not be marshalled: %v", err))
			return
		}

		if err := usecase.SaveMovie(ctx, svc, dto); err != nil {
			controller.internalServerError(ctx, "path broken", fmt.Sprintf("Could not create movie with body %+v: %v", dto, err))
			return
		}

		ctx.JSON(http.StatusCreated, http.NoBody)
	}
}


func (controller *MovieController) DeleteMovieHandler(usecase ports.DeleteMovieCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exists := controller.getService(ctx)
		if !exists {
			return
		}

		svc, ok := service.(ports.MovieDeleterService)
		if !ok {
			controller.internalServerError(ctx, "service unavailable.", "Service malformed.")
			return
		}

		id, ok := controller.getId(ctx)
		if !ok {
			return
		}

		if err := usecase.DeleteMovie(ctx, svc, dtos.MovieId(id)); err != nil {
			controller.internalServerError(ctx, "path broken", fmt.Sprintf("Could not delete movie with id %d: %v", id, err))
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}


func (controller *MovieController) internalServerError(ctx *gin.Context, msg, logging string ) {
	log.Println(logging)
	ctx.JSON(http.StatusInternalServerError, errors.InternalServerError(msg))
	ctx.Abort()
}

func (controller *MovieController) unprocessableEntityError(ctx *gin.Context, msg, logging string ) {
	log.Println(logging)
	ctx.JSON(http.StatusUnprocessableEntity, errors.UnprocessableEntity(msg))
	ctx.Abort()	
}

func (controller *MovieController) getService(ctx *gin.Context) (any, bool) {
	service, exists := ctx.Get(ports.ServiceKey)
	if !exists {
		controller.internalServerError(ctx, "service unavailable.", "Service not set to context.")
	}
	return service, exists
}

func (controller *MovieController) getId(ctx *gin.Context) (int, bool) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		controller.unprocessableEntityError(ctx, "malformed id.", fmt.Sprintf("Could not parse id: %v", err))
		return 0, false
	}
	return id, true
}
