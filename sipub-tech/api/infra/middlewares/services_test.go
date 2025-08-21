package middlewares_test

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"	

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/ports"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/middlewares"
)

func TestAddMovieQueryService(t *testing.T) {
	ctx := &gin.Context{}

	queryService := &FakeQueryService{}

	f := middlewares.AddMovieQueryService(queryService)
	f(ctx)

	service, exists := ctx.Get(ports.ServiceKey)
	if !exists {
		t.Errorf("service not set")
	}

	assert.Equal(t, queryService, service)
}

func TestAddMovieExecutorService(t *testing.T) {
	ctx := &gin.Context{}

	executorService := &FakeExecutorService{}

	f := middlewares.AddMovieExecutorService(executorService)
	f(ctx)

	service, exists := ctx.Get(ports.ServiceKey)
	if !exists {
		t.Errorf("service not set")
	}

	assert.Equal(t, executorService, service)
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

 

