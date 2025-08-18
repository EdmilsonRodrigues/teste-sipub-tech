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

type ContextKey string

const (
	RepoKey ContextKey = "repository"
)

var (
	ErrUnsetRespository = fmt.Errorf("repository not set in context or incomplete")
)

type GRPCMovieController struct {
	pb.UnimplementedMovieServiceServer
}

func (controller *GRPCMovieController) GetMovie(ctx context.Context, req *pb.GetMovieRequest) (*pb.Movie, error) {
	repo, ok := ctx.Value(RepoKey).(ports.MovieOneGetterRepository)
	if !ok {
		return nil, ErrUnsetRespository
	}
	usecase := usecases.NewGetMovieCase(repo)

	movie, err := usecase.GetMovie(dtos.MovieID(req.Id))
	if err == ports.ErrMovieNotFound {
		return nil, pb_exceptions.ErrMovieNotFound
	}

	return controller.responseDtoToPbMovie(movie), err	
}
	
func (controller *GRPCMovieController) GetMovies(ctx context.Context, req *pb.GetMoviesRequest) (*pb.Movies, error) {
	repo, ok := ctx.Value(RepoKey).(ports.MovieAllGetterRepository)
	if !ok {
		return nil, ErrUnsetRespository
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

type MessagingMovieController struct {

}

func (controler *MessagingMovieController) SaveMovie(ctx context.Context, movie dtos.CreateMovieDTO) error {
	repo, ok := ctx.Value(RepoKey).(ports.MovieSaverRepository)
	if !ok {
		return ErrUnsetRespository
	}
	usecase := usecases.NewSaveMovieCase(repo)

	return usecase.SaveMovie(movie)
}


func (controller *MessagingMovieController) DeleteMovie(ctx context.Context, id dtos.MovieID) error {
	repo, ok := ctx.Value(RepoKey).(ports.MovieDeleterRepository)
	if !ok {
		return ErrUnsetRespository
	}
	usecase := usecases.NewDeleteMovieCase(repo)

	return usecase.DeleteMovie(id)
}



