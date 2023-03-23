package service

import (
	"fmt"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	"github.com/rs/zerolog/log"
	grpc "google.golang.org/grpc"
)

// ConnectMinions attepmts to connect to all minions in the config
func ConnectMinions() ([]*minion, []*minion) {
	connectedMinionBuffer := make([]*minion, 0)
	disconnectedMinionBuffer := make([]*minion, 0)

	for _, minionConfig := range config.GetMinions() {
		minionBuffer := &minion{
			config: minionConfig,
		}

		// TODO: mTLS
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", minionConfig.Address, minionConfig.Port), grpc.WithInsecure())
		if err != nil {
			disconnectedMinionBuffer = append(disconnectedMinionBuffer, minionBuffer)
			log.Error().Err(err).Msgf("Failed to connect to minion %s:%d", minionConfig.Address, minionConfig.Port)
		}

		serviceClient := NewServiceClient(conn)
		minionBuffer.client = serviceClient

		connectedMinionBuffer = append(connectedMinionBuffer, minionBuffer)
		log.Info().Msgf("Connected to minion %s:%d", minionConfig.Address, minionConfig.Port)
	}

	return connectedMinionBuffer, disconnectedMinionBuffer
}
