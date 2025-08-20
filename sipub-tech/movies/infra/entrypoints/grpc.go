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
		repo: repo,
		listeningPort: listeningPort,
		controller: &controllers.GRPCMovieController{},
	}
}

type GRPCEntrypoint struct {
	repo ports.MovieQueryRepository
	controller *controllers.GRPCMovieController
	listeningPort int
}


func (entrypoint *GRPCEntrypoint) Serve(ctx context.Context) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", entrypoint.listeningPort))
	if err != nil {
		panic("movie grpc entrypoint failed to listen to desired port")
	}

	s := grpc.NewServer()

	pb.RegisterMovieServiceServer(s, entrypoint.controller)
	log.Printf("Listening on port %d\n", entrypoint.listeningPort)
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	
}



