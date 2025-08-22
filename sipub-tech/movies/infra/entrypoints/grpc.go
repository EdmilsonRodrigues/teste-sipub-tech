package entrypoints

import (
	"fmt"
	"net"
	"context"
	"log"

	"google.golang.org/grpc"
	pb "github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/grpc/movies"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/ports"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/infra/controllers"
)

func NewGRPCEntrypoint(repo ports.MovieQueryRepository, listeningPort int) *GRPCEntrypoint {
	return &GRPCEntrypoint{
		listeningPort: listeningPort,
		server: newGRPCServer(repo, &controllers.GRPCMovieController{}),
	}
}

type GRPCEntrypoint struct {
	server *gRPCServer
	listeningPort int
}


func (entrypoint *GRPCEntrypoint) Serve(ctx context.Context) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", entrypoint.listeningPort))
	if err != nil {
		panic("movie grpc entrypoint failed to listen to desired port")
	}

	s := grpc.NewServer()

	pb.RegisterMovieServiceServer(s, entrypoint.server)
	log.Printf("Listening on port %d\n", entrypoint.listeningPort)
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	
}

func newGRPCServer(repo ports.MovieQueryRepository, controller *controllers.GRPCMovieController) *gRPCServer {
	return &gRPCServer{
		repo: repo,
		controller: controller,
	}
}

type gRPCServer struct {
	pb.MovieServiceServer

	repo ports.MovieQueryRepository
	controller *controllers.GRPCMovieController	
}

func (server *gRPCServer) GetMovie(ctx context.Context, req *pb.GetMovieRequest) (*pb.Movie, error) {
	ctx = context.WithValue(ctx, controllers.RepoKey, server.repo)
	return server.controller.GetMovie(ctx, req)
}

func (server *gRPCServer) GetMovies(ctx context.Context, req *pb.GetMoviesRequest) (*pb.Movies, error) {
	ctx = context.WithValue(ctx, controllers.RepoKey, server.repo)
	return server.controller.GetMovies(ctx, req)
} 

