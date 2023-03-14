package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreaswachs/bachelors-project/daaukins/server/service"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
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

// Credit: https://github.com/islishude/grpc-mtls-example
func middlefunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// get client tls info
	if p, ok := peer.FromContext(ctx); ok {
		if mtls, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			for _, item := range mtls.State.PeerCertificates {
				log.Info().Msgf("request certificate subject:", item.Subject)
			}
		}
	}
	return handler(ctx, req)
}

// Credit: https://github.com/islishude/grpc-mtls-example
// TODO: Change paths to certs
func LoadKeyPair() credentials.TransportCredentials {
	certificate, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Panic().Msgf("failed to load server certification: " + err.Error())
	}

	data, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		log.Panic().Msgf("failed to load CA file: " + err.Error()
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(data) {
		log.Panic().Msgf("can't add ca cert")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    capool,
	}
	return credentials.NewTLS(tlsConfig)
}
