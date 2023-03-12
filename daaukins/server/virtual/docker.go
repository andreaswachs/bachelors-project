// Package to manage virtualization assets

package virtual

import (
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	client *docker.Client
)

func init() {
	var err error
	client, err = docker.NewClientFromEnv()
	if err != nil {
		log.Panic(err)
	}
}

func DockerClient() *docker.Client {
	return client
}
