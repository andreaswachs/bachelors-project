package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	"github.com/andreaswachs/bachelors-project/daaukins/server/labs"
	"github.com/andreaswachs/bachelors-project/daaukins/server/service"
	"github.com/andreaswachs/bachelors-project/daaukins/server/store"
	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
	"github.com/rs/zerolog/log"
)

func main() {
	flag.Parse()

	if err := config.Initialize(); err != nil {
		return
	}

	if err := virtual.Initialize(); err != nil {
		log.Error().
			Err(err).
			Msg("failed to initialize virtualization")
		return
	}

	if err := store.Initialize(); err != nil {
		log.Error().
			Err(err).
			Msg("failed to initialize store")
		return
	}

	service.Initialize()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Info().Msg("shutting down server")
	if err := labs.RemoveAll(); err != nil {
		log.Error().
			Err(err).
			Msg("failed to remove all labs. Run `make clean-docker` to cleanup manually")
	}

	service.Stop()
}
