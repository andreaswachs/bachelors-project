package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
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
		minionBuffer := &minion{
			config: minionConfig,
		}

		wg.Add(1)
		go func(m *minion) {
			defer wg.Done()

			// TODO: mTLS
			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", minionConfig.Address, minionConfig.Port), grpc.WithInsecure())
			if err != nil {
				disconnectedMinionBufferLock.Lock()
				defer disconnectedMinionBufferLock.Unlock()

				disconnectedMinionBuffer = append(disconnectedMinionBuffer, minionBuffer)
				log.Error().Err(err).Msgf("Failed to connect to minion %s:%d", minionConfig.Address, minionConfig.Port)
				return
			}

			serviceClient := NewServiceClient(conn)
			minionBuffer.client = serviceClient

			response, err := serviceClient.GetServerMode(context.Background(), &GetServerModeRequest{})
			if err != nil {
				log.Error().Err(err).Msgf("Failed to get server mode from minion %s:%d", minionConfig.Address, minionConfig.Port)
			}

			if response.Mode == config.ModeLeader.String() {
				log.Error().Msgf("Minion %s:%d is a leader, but it should be a follower", minionConfig.Address, minionConfig.Port)
				return
			}

			connectedMinionBufferLock.Lock()
			defer connectedMinionBufferLock.Unlock()

			connectedMinionBuffer = append(connectedMinionBuffer, minionBuffer)
			log.Info().Msgf("Connected to minion %s:%d", minionConfig.Address, minionConfig.Port)
		}(minionBuffer)

	}

	wg.Wait()

	return connectedMinionBuffer, disconnectedMinionBuffer
}
