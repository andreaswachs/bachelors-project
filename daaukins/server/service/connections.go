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

// ConnectMinions attepmts to connect to all minions in the config
func ConnectMinions() ([]*minion, []*minion) {
	connectedMinionBuffer := make([]*minion, 0)
	disconnectedMinionBuffer := make([]*minion, 0)

	connectedMinionBufferLock := sync.Mutex{}
	disconnectedMinionBufferLock := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, minionConfig := range config.GetMinions() {
		wg.Add(1)

		go func(m *config.MinionConfig) {
			defer wg.Done()
			minionBuffer := &minion{
				config: *m,
			}

			// TODO: mTLS
			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", m.Address, m.Port), grpc.WithInsecure())
			if err != nil {
				disconnectedMinionBufferLock.Lock()
				defer disconnectedMinionBufferLock.Unlock()

				disconnectedMinionBuffer = append(disconnectedMinionBuffer, minionBuffer)
				log.Error().Err(err).Msgf("Failed to connect to follower %s:%d", m.Address, m.Port)
				return
			}

			serviceClient := service.NewServiceClient(conn)
			minionBuffer.client = serviceClient

			response, err := serviceClient.GetServerMode(context.Background(), &service.GetServerModeRequest{})
			if err != nil {
				log.Error().Err(err).Msgf("Failed to get server mode from follower %s:%d", m.Address, m.Port)
			}

			if response.Mode == config.ModeLeader.String() {
				log.Error().Msgf("Server %s:%d is a leader, but it should be a follower. The server will not be used.", m.Address, m.Port)
				return
			}

			connectedMinionBufferLock.Lock()
			defer connectedMinionBufferLock.Unlock()

			connectedMinionBuffer = append(connectedMinionBuffer, minionBuffer)
			log.Info().Msgf("Connected to follower %s:%d", m.Address, m.Port)
		}(&minionConfig)
	}

	wg.Wait()

	return connectedMinionBuffer, disconnectedMinionBuffer
}
