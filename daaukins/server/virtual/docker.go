// Package to manage virtualization assets

package virtual

import (
	docker "github.com/fsouza/go-dockerclient"
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
	return client
}
