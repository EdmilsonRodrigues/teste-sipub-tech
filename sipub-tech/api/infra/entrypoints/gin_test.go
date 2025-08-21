package entrypoints_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
	"testing"

	"github.com/stretchr/testify/assert"	
	"github.com/gin-gonic/gin"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/ports"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/entrypoints"
)

func TestGinEntrypoint(t *testing.T) {
	executorService := &FakeExecutorService{}
	queryService := &FakeQueryService{}

	movieController := &MockMovieController{}

	entrypoint := entrypoints.NewGinEntrypoint(
		executorService,
		queryService,
		movieController,
	)
	entrypoint.Setup()
	engine := entrypoint.GetEngine()
	
	t.Run("should call GetMoviesHandler when hit a GET to /movies/", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/movies/", nil)

		engine.ServeHTTP(w, req)

		assert.Equal(t, queryService, movieController.GetMoviesService)
	})

	t.Run("should correctly parse query to a dto and set in context when hit a GET to /movies/", func(t *testing.T) {
		limit, cursor := rand.Int32(), rand.Int32()

		year := rand.IntN(time.Now().Year() - 1880) + 1880
		t.Logf("year: %d", year)

		path := fmt.Sprintf("/movies/?year=%d&limit=%d&cursor=%d", year, limit, cursor)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", path, nil)

		engine.ServeHTTP(w, req)

		expected := dtos.MoviesQueryDTO{
			Year: strconv.Itoa(year),
			Limit: int(limit),
			Cursor: int(cursor),
		}

		assert.Equal(t, expected, movieController.GetMoviesDTO)
	})

	t.Run("should call GetMovieHandler when hit a GET to /movies/:id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/movies/75", nil)

		engine.ServeHTTP(w, req)

		assert.Equal(t, queryService, movieController.GetMovieService)
	})

	t.Run("should call SaveMovieHandler when hit a POST to /movies", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/movies/", nil)

		engine.ServeHTTP(w, req)

		assert.Equal(t, executorService, movieController.SaveMovieService)
	})

	t.Run("should call DeleteMovieHandler when hit a GET to /movies/:id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/movies/75", nil)

		engine.ServeHTTP(w, req)

		assert.Equal(t, executorService, movieController.DeleteMovieService)
	})
}

type MockMovieController struct {
	GetMovieService     any
	GetMovieError       error

	GetMoviesService    any
	GetMoviesError      error
	GetMoviesDTO        dtos.MoviesQueryDTO

	SaveMovieService    any
	SaveMovieError      error

	DeleteMovieService  any
	DeleteMovieError    any
}

func (controller *MockMovieController) GetMovieHandler(usecase ports.GetMovieCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exists := ctx.Get(ports.ServiceKey)
		if !exists {
			controller.GetMovieError = fmt.Errorf("query service not set")
		}

		controller.GetMovieService = service
		ctx.JSON(204, http.NoBody)
	}
}

func (controller *MockMovieController) GetMoviesHandler(usecase ports.GetMoviesCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exists := ctx.Get(ports.ServiceKey)
		if !exists {
			controller.GetMoviesError = fmt.Errorf("query service not set")
		}

		dto, exists := ctx.Get("dto")
		if !exists {
			controller.GetMoviesError = fmt.Errorf("query dto not passed to context")
		}
		parsedDto, ok := dto.(*dtos.MoviesQueryDTO)
		if !ok {
			controller.GetMoviesError = fmt.Errorf("malformed query dto")
		}
		controller.GetMoviesDTO = *parsedDto

		controller.GetMoviesService = service
		ctx.JSON(204, http.NoBody)
	}
}

func (controller *MockMovieController) SaveMovieHandler(usecase ports.SaveMovieCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exists := ctx.Get(ports.ServiceKey)
		if !exists {
			controller.SaveMovieError = fmt.Errorf("executor service not set")
		}

		controller.SaveMovieService = service
		ctx.JSON(204, http.NoBody)
	}
}

func (controller *MockMovieController) DeleteMovieHandler(usecase ports.DeleteMovieCase)  gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exists := ctx.Get(ports.ServiceKey)
		if !exists {
			controller.DeleteMovieError = fmt.Errorf("executor service not set")
		}

		controller.DeleteMovieService = service
		ctx.JSON(204, http.NoBody)
	}
}

type FakeQueryService struct {
	ports.MovieQueryService
}

func (service *FakeQueryService) GetOne(ctx context.Context, id dtos.MovieId) (dtos.MovieResponseDTO, error) {
	return dtos.MovieResponseDTO{}, nil
}
	
func (service *FakeQueryService) GetAll(ctx context.Context, query dtos.MoviesQueryDTO) (dtos.MoviesResponseDTO, error) {
	return dtos.MoviesResponseDTO{}, nil
}

 
type FakeExecutorService struct {
	ports.MovieExecutorService
}

func (service *FakeExecutorService) Save(ctx context.Context, movie dtos.CreateMovieDTO) error {
	return nil
}
	
func (service *FakeExecutorService) Delete(ctx context.Context, id dtos.MovieId) error {
	return nil
}

 
