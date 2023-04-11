package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	service "github.com/andreaswachs/daaukins-service"
	"github.com/rs/zerolog/log"
	grpc "google.golang.org/grpc"
)

// ConnectFollowers attepmts to connect to all follower in the config
func ConnectFollowers() ([]*follower, []*follower) {
	connectedFollowerBuffer := make([]*follower, 0)
	disconnectedFollowerBuffer := make([]*follower, 0)

	connectedFollowerBufferLock := sync.Mutex{}
	disconnectedFollowerBufferLock := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, followerConfig := range config.GetFollowers() {
		wg.Add(1)

		go func(m *config.FollowerConfig) {
			defer wg.Done()
			followerBuffer := &follower{
				config: *m,
			}

			// TODO: mTLS
			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", m.Address, m.Port), grpc.WithInsecure())
			if err != nil {
				disconnectedFollowerBufferLock.Lock()
				defer disconnectedFollowerBufferLock.Unlock()

				disconnectedFollowerBuffer = append(disconnectedFollowerBuffer, followerBuffer)
				log.Error().Err(err).Msgf("Failed to connect to follower %s:%d", m.Address, m.Port)
				return
			}

			serviceClient := service.NewServiceClient(conn)
			followerBuffer.client = serviceClient

			response, err := serviceClient.GetServerMode(context.Background(), &service.GetServerModeRequest{})
			if err != nil {
				log.Error().Err(err).Msgf("Failed to get server mode from follower %s:%d", m.Address, m.Port)
			}

			if response.Mode == config.ModeLeader.String() {
				log.Error().Msgf("Server %s:%d is a leader, but it should be a follower. The server will not be used.", m.Address, m.Port)
				return
			}

			connectedFollowerBufferLock.Lock()
			defer connectedFollowerBufferLock.Unlock()

			connectedFollowerBuffer = append(connectedFollowerBuffer, followerBuffer)
			log.Info().Msgf("Connected to follower %s:%d", m.Address, m.Port)
		}(&followerConfig)
	}

	wg.Wait()

	return connectedFollowerBuffer, disconnectedFollowerBuffer
}
