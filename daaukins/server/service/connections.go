package service

import (
	"context"
	"fmt"
	"time"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	service "github.com/andreaswachs/daaukins-service"
	"github.com/rs/zerolog/log"
	grpc "google.golang.org/grpc"
)

type connectionStatus uint8

const (
	connected connectionStatus = iota
	disconnected
)

type connAttempt struct {
	follower *follower
	result   connectionStatus
}

// ConnectFollowers attepmts to connect to all follower in the config
func ConnectFollowers() ([]*follower, []*follower) {
	connectedFollowerBuffer := make([]*follower, 0)
	disconnectedFollowerBuffer := make([]*follower, 0)

	attemps := make(chan connAttempt)

	for _, follower := range disconnectedFollowers {
		log.Debug().
			Str("address", follower.config.Address).
			Int("port", follower.config.Port).
			Msg("Attempting to connect to follower")

		go connect(follower, attemps)
	}

	for i := 0; i < len(disconnectedFollowers); i++ {
		attempt := <-attemps
		switch attempt.result {
		case connected:
			connectedFollowerBuffer = append(connectedFollowerBuffer, attempt.follower)
		case disconnected:
			disconnectedFollowerBuffer = append(disconnectedFollowerBuffer, attempt.follower)
		}
	}

	return connectedFollowerBuffer, disconnectedFollowerBuffer
}

func connect(f *follower, comm chan<- connAttempt) {
	log.Debug().
		Str("address", f.config.Address).
		Int("port", f.config.Port).
		Msg("Begin attempt to connect to follower")

	ctx, cancel := shortTimeoutContext()
	defer cancel()

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", f.config.Address, f.config.Port), grpc.WithInsecure())
	if err != nil {
		comm <- connAttempt{
			follower: f,
			result:   disconnected,
		}

		log.Error().Err(err).Msgf("Failed to connect to follower %s:%d", f.config.Address, f.config.Port)
		return
	}

	serviceClient := service.NewServiceClient(conn)
	f.client = serviceClient

	response, err := serviceClient.GetServerMode(context.Background(), &service.GetServerModeRequest{})
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get server mode from follower %s:%d", f.config.Address, f.config.Port)
		comm <- connAttempt{
			follower: f,
			result:   disconnected,
		}
		return
	}

	if response.Mode == config.ModeLeader.String() {
		log.Error().Msgf("Server %s:%d is a leader, but it should be a follower. The server will not be used.", f.config.Address, f.config.Port)
		comm <- connAttempt{
			follower: f,
			result:   disconnected,
		}
		return
	}

	f.serverId = response.ServerId

	log.Info().Msgf("Connected to follower %s:%d", f.config.Address, f.config.Port)
	comm <- connAttempt{
		follower: f,
		result:   connected,
	}
}

func shortTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}
