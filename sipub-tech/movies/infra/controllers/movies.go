package controllers

import (
	"context"
	"fmt"
	
	pb "github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/grpc/movies"
	pb_exceptions "github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/grpc/exceptions"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/ports"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/usecases"
)

const (
	RepoKey = "repository"
)

var (
	UnsetRespositoryError = fmt.Errorf("Repository not set in context or incomplete.")
)

type GRPCMovieController struct {
	pb.UnimplementedMovieServiceServer
}

func (controller *GRPCMovieController) GetMovie(ctx context.Context, req *pb.GetMovieRequest) (*pb.Movie, error) {
	repo, ok := ctx.Value(RepoKey).(ports.MovieOneGetterRepository)
	if !ok {
		return nil, UnsetRespositoryError
	}
	usecase := usecases.NewGetMovieCase(repo)

	movie, err := usecase.GetMovie(dtos.MovieID(req.Id))
	if err == ports.MovieNotFoundError {
		return nil, pb_exceptions.MovieNotFoundException
	}

	return controller.responseDtoToPbMovie(movie), err	
}
	
func (controller *GRPCMovieController) GetMovies(ctx context.Context, req *pb.GetMoviesRequest) (*pb.Movies, error) {
	repo, ok := ctx.Value(RepoKey).(ports.MovieAllGetterRepository)
	if !ok {
		return nil, UnsetRespositoryError
	}
	usecase := usecases.NewGetMoviesCase(repo)

	movies, err := usecase.GetMovies()
	if err != nil {
		return nil, err
	}

	parsedMovies := make([]*pb.Movie, len(*movies))

	for index, movie := range(*movies) {
		parsedMovies[index] = controller.responseDtoToPbMovie(&movie)
	}
	return &pb.Movies{Movies: parsedMovies}, err
}

func (controller *GRPCMovieController) responseDtoToPbMovie(movie *dtos.MovieResponseDTO) *pb.Movie {
	return &pb.Movie{
		Id: int32(movie.ID),
		Title: movie.Title,
		Year: movie.Year,
	}
}




