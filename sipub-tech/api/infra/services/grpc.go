package services

import (
	"context"
	"fmt"

	pb "github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/grpc/movies"
	pb_exceptions "github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/grpc/exceptions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/ports"
)


func NewMovieGRPCService(serverUrl string) *MovieGRPCService {
	conn, err := grpc.NewClient(serverUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("could not connect with grpc server.")
	}

	movieClient := pb.NewMovieServiceClient(conn)

	return &MovieGRPCService{
		client: movieClient,
		conn: conn,
	}
}

type MovieGRPCService struct {
	ports.MovieQueryService

	client pb.MovieServiceClient
	conn *grpc.ClientConn
}

func (service *MovieGRPCService) Close() {
	if err := service.conn.Close(); err != nil {
		panic("grpc connection failed to close")
	}
}

func (service *MovieGRPCService) GetOne(ctx context.Context, id dtos.MovieId) (dtos.MovieResponseDTO, error) {
	response, err := service.client.GetMovie(
		ctx,
		&pb.GetMovieRequest{Id: int32(id)},
	)
	if err != nil {
		empty := dtos.MovieResponseDTO{} 
		if err == pb_exceptions.ErrMovieNotFound {
			return empty, ports.ErrMovieNotFound
		}

		return empty, fmt.Errorf("failed getting movie: %w", err)
	}
	
	return service.parseMovieResponse(response), nil
}

func (service *MovieGRPCService) GetAll(ctx context.Context, query dtos.MoviesQueryDTO) (dtos.MoviesResponseDTO, error) {
	var movies dtos.MoviesResponseDTO

	response, err := service.client.GetMovies(
		ctx,
		&pb.GetMoviesRequest{
			Year:   query.Year,
			Limit:  int32(query.Limit),
			Cursor: int32(query.Cursor),
		},
	)
	if err != nil {
		return movies, fmt.Errorf("failed getting movies: %w", err)
	}

	movies = dtos.MoviesResponseDTO{
		Movies: service.parseMovieResponseArray(response.Movies),
		Cursor: int(response.Cursor),
	}
	
	return movies, nil
}

func (service *MovieGRPCService) parseMovieResponse(movie *pb.Movie) dtos.MovieResponseDTO {
	return dtos.MovieResponseDTO{
		ID: int(movie.Id),
		Title: movie.Title,
		Year: movie.Year,
	}
}

func (service *MovieGRPCService) parseMovieResponseArray(movies []*pb.Movie) []*dtos.MovieResponseDTO {
	response := make([]*dtos.MovieResponseDTO, len(movies))
	for index, movie := range movies {
		dto := service.parseMovieResponse(movie)
		response[index] = &dto
	}
	return response
}

