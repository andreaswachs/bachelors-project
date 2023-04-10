package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	"github.com/andreaswachs/bachelors-project/daaukins/server/labs"
	"github.com/andreaswachs/bachelors-project/daaukins/server/service"
	"github.com/rs/zerolog/log"
)

var (
	configFilename = flag.String("config", "server.yaml", "path to config file")
)

func main() {
	flag.Parse()

	config.Initialize(&config.InitializeConfigOptions{
		ConfigFile: *configFilename,
	})

	service.Initialize()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Info().Msg("shutting down server")
	if err := labs.RemoveAll(); err != nil {
		log.Error().Err(err).Msg("failed to remove all labs. Run `make clean-docker` to cleanup manually")
	}

	service.Stop()
}
