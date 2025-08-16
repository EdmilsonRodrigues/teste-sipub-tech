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
	repoKey = "repository"
)

type GRPCMovieController struct {
	pb.UnimplementedMovieServiceServer

	
}

func (controller *GRPCMovieController) GetMovie(ctx context.Context, req *pb.GetMovieRequest) (*pb.Movie, error) {
	repo, ok := ctx.Value(repoKey).(ports.MovieOneGetterRepository)
	if !ok {
		return nil, fmt.Errorf("Repository not set in context or incomplete.")
	}
	usecase := usecases.NewGetMovieCase(repo)

	movie, err := usecase.GetMovie(dtos.MovieID(req.Id))
	if err == ports.MovieNotFoundError {
		return nil, pb_exceptions.MovieNotFoundException
	}
	
}
	
func (controller *GRPCMovieController) GetMovies(ctx context.Context, req *pb.GetMoviesRequest) (*pb.Movies, error) {
	
}
