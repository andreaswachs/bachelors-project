package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreaswachs/bachelors-project/daaukins/server/service"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer(
	// grpc.Creds(LoadKeyPair()),
	// grpc.UnaryInterceptor(middlefunc),
	)

	service.RegisterServiceServer(server, new(service.Server))

	go func() {
		l, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Panic().Err(err).Msgf("failed to listen on port %d", 1905)
		}
		log.Info().Msgf("listening on port %d", 1905)
		if err := server.Serve(l); err != nil {
			log.Panic().Err(err).Msg("failed to serve")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Info().Msg("shutting down server")
	server.GracefulStop()
}
