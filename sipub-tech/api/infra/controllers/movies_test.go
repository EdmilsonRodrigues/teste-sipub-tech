package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	//"runtime"
	"strconv"
	"testing"
	"testing/quick"
	
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"	

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/ports"
	infraDtos "github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/controllers"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/middlewares"
)

func TestMovieController(t *testing.T) {
	controller := controllers.NewMovieController()
	t.Run("the function returned by controllers.MovieController.GetMovieHandler must", func(t *testing.T) {
		t.Run("Return a success response with the movie gotten by the usecase", func(t *testing.T) {
			assertion := func(id uint16) bool {
				handler := controller.GetMovieHandler(&MockGetMovieCase{})

				req, _ := http.NewRequest("GET", fmt.Sprintf("/movies/%d", int(id)), nil)
				ctx, writer := getContext(req)
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(int(id))})
				ctx.Set(ports.ServiceKey, &FakeQueryService{})

				handler(ctx)

				var body infraDtos.JSONResponse

				if !assert.False(t, ctx.IsAborted()) {
					return false
				}

				if !assert.Equal(t, 200, writer.Status()) {
					return false
				}
				
				if err := json.Unmarshal(writer.Body, &body); err != nil {
					t.Logf("Error unmarshalling body %+v: %v", writer.Body, err)
					return false
				}

				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("Failed checking assertion: %v", err)
			}
		})

		t.Run("return an unprocessable entity response if the passed id is not an integer", func(t *testing.T) {
			assertion := func(id float32) bool {
				handler := controller.GetMovieHandler(&MockGetMovieCase{})

				req, _ := http.NewRequest("GET", fmt.Sprintf("/movies/%f", id), nil)
				ctx, writer := getContext(req)
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: fmt.Sprintf("%f", id)})
				ctx.Set(ports.ServiceKey, &FakeQueryService{})

				handler(ctx)

				var body infraDtos.ErrorResponse

				if !assert.True(t, ctx.IsAborted()) {
					return false
				}

				if !assert.Equal(t, 422, writer.Status()) {
					return false
				}
				
				if err := json.Unmarshal(writer.Body, &body); err != nil {
					t.Logf("Error unmarshalling body %+v: %v", writer.Body, err)
					return false
				}

				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("Failed checking assertion: %v", err)
			}
		})

		t.Run("return a movie not found response if the usecase return a ports.ErrMovieNotFound", func(t *testing.T) {
			assertion := func(id uint16) bool {
				handler := controller.GetMovieHandler(&MockGetMovieCase{ErrorReturned: ports.ErrMovieNotFound})

				req, _ := http.NewRequest("GET", fmt.Sprintf("/movies/%d", id), nil)
				ctx, writer := getContext(req)
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: fmt.Sprintf("%d", id)})
				ctx.Set(ports.ServiceKey, &FakeQueryService{})

				handler(ctx)

				var body infraDtos.ErrorResponse

				if !assert.True(t, ctx.IsAborted()) {
					return false
				}

				if !assert.Equal(t, 404, writer.Status()) {
					return false
				}
				
				if err := json.Unmarshal(writer.Body, &body); err != nil {
					t.Logf("Error unmarshalling body %+v: %v", writer.Body, err)
					return false
				}

				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("Failed checking assertion: %v", err)
			}
		})
	})

	t.Run("the function returned by controllers.MovieController.GetMoviesHandler must", func(t *testing.T) {
		t.Run("return a paginated response", func(t *testing.T) {
			assertion := func() bool {
				handler := controller.GetMoviesHandler(&MockGetMoviesCase{})

				req, _ := http.NewRequest("GET", "/movies", nil)
				ctx, writer := getContext(req)
				ctx.Set(middlewares.DtoKey, &dtos.MoviesQueryDTO{})
				ctx.Set(ports.ServiceKey, &FakeQueryService{})

				handler(ctx)

				var body infraDtos.PaginatedJSONResponse

				if !assert.False(t, ctx.IsAborted()) {
					return false
				}

				if !assert.Equal(t, 200, writer.Status()) {
					return false
				}
				
				if err := json.Unmarshal(writer.Body, &body); err != nil {
					t.Logf("Error unmarshalling body %+v: %v", writer.Body, err)
					return false
				}

				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("Failed checking assertion: %v", err)
			}
		})

	})

	t.Run("the function returned by controllers.MovieController.SaveMovieHandler must", func(t *testing.T) {
		t.Run("return a success response with no body", func(t *testing.T) {
			assertion := func(movie dtos.CreateMovieDTO) bool {
				handler := controller.SaveMovieHandler(&StubSaveMovieCase{})

				movieBytes, err := json.Marshal(movie)
				if err != nil {
					t.Logf("Failed marshalling dto %+v: %v", movie, err)
				}
				req, _ := http.NewRequest("POST", "/movies/", bytes.NewReader(movieBytes))
				ctx, writer := getContext(req)
				ctx.Set(ports.ServiceKey, &FakeExecutorService{})

				handler(ctx)

				var body infraDtos.PaginatedJSONResponse

				if !assert.False(t, ctx.IsAborted()) {
					return false
				}

				if !assert.Equal(t, 201, writer.Status()) {
					return false
				}
				
				if err := json.Unmarshal(writer.Body, &body); err != nil {
					t.Logf("Error unmarshalling body %+v: %v", writer.Body, err)
					return false
				}

				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("Failed checking assertion: %v", err)
			}
		})

		t.Run("return an unprocessable entity response if the body could not be parsed", func(t *testing.T) {
			assertion := func() bool {
				handler := controller.SaveMovieHandler(&StubSaveMovieCase{})

				req, _ := http.NewRequest("POST", "/movies/", nil)
				ctx, writer := getContext(req)
				ctx.Set(ports.ServiceKey, &FakeExecutorService{})

				handler(ctx)

				var body infraDtos.ErrorResponse

				if !assert.True(t, ctx.IsAborted()) {
					return false
				}

				if !assert.Equal(t, 422, writer.Status()) {
					return false
				}
				
				if err := json.Unmarshal(writer.Body, &body); err != nil {
					t.Logf("Error unmarshalling body %+v: %v", writer.Body, err)
					return false
				}

				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("Failed checking assertion: %v", err)
			}
		})

	})

	t.Run("the function returned by controllers.MovieController.DeleteMovieHandler must", func(t *testing.T) {
		t.Run("return a success response with no body", func(t *testing.T) {
			assertion := func(id uint16) bool {
				handler := controller.DeleteMovieHandler(&StubDeleteMovieCase{})

				req, _ := http.NewRequest("DELETE", fmt.Sprintf("/movies/%d", int(id)), nil)
				ctx, writer := getContext(req)
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(int(id))})
				ctx.Set(ports.ServiceKey, &FakeExecutorService{})

				handler(ctx)

				if !assert.False(t, ctx.IsAborted()) {
					return false
				}

				if !assert.Equal(t, 204, writer.Status()) {
					return false
				}
 				
				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("Failed checking assertion: %v", err)
			}
		})

		t.Run("return an unprocessable entity response if the passed id is not an integer", func(t *testing.T) {
			assertion := func(id float32) bool {
				handler := controller.DeleteMovieHandler(&StubDeleteMovieCase{})

				req, _ := http.NewRequest("DELETE", fmt.Sprintf("/movies/%f", id), nil)
				ctx, writer := getContext(req)
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: fmt.Sprintf("%f", id)})
				ctx.Set(ports.ServiceKey, &FakeExecutorService{})

				handler(ctx)
				
				var body infraDtos.ErrorResponse

				if !assert.True(t, ctx.IsAborted()) {
					return false
				}

				if !assert.Equal(t, 422, writer.Status()) {
					return false
				}
				
				if err := json.Unmarshal(writer.Body, &body); err != nil {
					t.Logf("Error unmarshalling body %+v: %v", writer.Body, err)
					return false
				}

				return true
			}
			if err := quick.Check(assertion, nil); err != nil {
				t.Errorf("Failed checking assertion: %v", err)
			}
		})
	})
}

type MockGetMovieCase struct {
	ports.GetMovieCase

	ErrorReturned error
}

func (usecase *MockGetMovieCase) GetMovie(
	ctx context.Context, service ports.MovieOneGetterService, id dtos.MovieId,
) (movie dtos.MovieResponseDTO, err error) {
	err = usecase.ErrorReturned
	return 
}

type MockGetMoviesCase struct {
	ports.GetMoviesCase
}

func (usecase *MockGetMoviesCase) GetMovies(
	ctx context.Context, service ports.MovieAllGetterService, query dtos.MoviesQueryDTO,
) (movies dtos.MoviesResponseDTO, err error) {
	return
}

type StubSaveMovieCase struct {
	ports.SaveMovieCase
}

func (usecase *StubSaveMovieCase) SaveMovie(ctx context.Context, service ports.MovieSaverService, movie dtos.CreateMovieDTO) error {
	return nil
}

type StubDeleteMovieCase struct {
	ports.DeleteMovieCase
}

func (usecase *StubDeleteMovieCase) DeleteMovie(ctx context.Context, service ports.MovieDeleterService, id dtos.MovieId) error {
	return nil
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

 
type FakeExecutorService struct {}

func (service *FakeExecutorService) Save(ctx context.Context, movie dtos.CreateMovieDTO) error {
	return nil
}
	
func (service *FakeExecutorService) Delete(ctx context.Context, id dtos.MovieId) error {
	return nil
}

 


func getContext(req *http.Request) (*gin.Context, *FakeWriter) {
	writer := &FakeWriter{HeadersMapping: make(http.Header)}
	return &gin.Context{
		Request: req,
		Writer: writer,
	}, writer
}

type FakeWriter struct {
	gin.ResponseWriter

	StatusCode     int
	HeadersMapping http.Header
	Body           []byte
}

func (w *FakeWriter) WriteHeader(code int) {
	w.StatusCode = code
}

func (w *FakeWriter) Header() http.Header {
	return w.HeadersMapping
}

func (w *FakeWriter) Write(b []byte) (int, error) {
	w.Body = append(w.Body, b...)
	return len(b), nil
}

func (w *FakeWriter) Status() int {
	return w.StatusCode
}

