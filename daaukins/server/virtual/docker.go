// Package to manage virtualization assets

package virtual

import (
	docker "github.com/fsouza/go-dockerclient"
	"github.com/rs/zerolog/log"
)

var (
	client *docker.Client
)

func Initialize() error {
	var err error
	client, err = docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	return nil
}

func DockerClient() *docker.Client {
	if client == nil {
		// The user should have initialized the docker client before using it
		log.Panic().Msg("Docker client is not initialized")
	}
	return client
}
