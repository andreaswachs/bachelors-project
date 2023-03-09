// Package to manage virtualization assets

package virtual

import docker "github.com/fsouza/go-dockerclient"

var (
	client *docker.Client
)

func init() {
	var err error
	client, err = docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
}

func DockerClient() *docker.Client {
	return client
}
