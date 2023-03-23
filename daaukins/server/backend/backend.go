// The backend is the brains of the server.
package backend

import (
	"time"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	"github.com/andreaswachs/bachelors-project/daaukins/server/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

var (
	stdTimeout = 5 * time.Second

	connections map[string]*grpc.ClientConn
)

func init() {
	connections = make(map[string]*grpc.ClientConn)
}

func ConnectMinions() {
	for _, minion := range config.GetMinions() {
		conn, err := grpc.Dial(minion.Address,
			grpc.WithTransportCredentials(), // TODO: add mTLS
			grpc.WithBlock())
		if err != nil {
			log.Error().Err(err).Msgf("failed to connect to minion %s", minion.Address)
		}

		// NASA safety rule
		for safety := 0; safety < 1000000000; safety++ {
			newName := utils.RandomShortName()

			if _, ok := connections[newName]; !ok {
				connections[newName] = conn
				break
			}
		}

	}
}
