package frontend

import (
	"fmt"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	"github.com/andreaswachs/bachelors-project/daaukins/server/utils"
	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
	"github.com/andreaswachs/sizes"

	docker "github.com/fsouza/go-dockerclient"
)

type Proxy struct {
	container *docker.Container
}

type ProvisionProxyOptions struct {
	LabID string
}

func ProvisionProxy(options *ProvisionProxyOptions) (*Proxy, error) {
	container, err := virtual.DockerClient().CreateContainer(docker.CreateContainerOptions{
		Name: fmt.Sprintf("daaukins-proxy-%s", utils.RandomName()),
		Config: &docker.Config{
			Image:  config.GetDockerConfig().Proxy.Image,
			Memory: sizes.Megabytes[int64](128),
			Labels: map[string]string{
				"daaukins":         "network-service",
				"daaukins.service": "proxy",
				"daaukins.lab":     options.LabID,
			},
		},
		HostConfig: &docker.HostConfig{},
	})

	if err != nil {
		return nil, err
	}

	return &Proxy{container: container}, nil
}

func (p *Proxy) Start() error {
	return virtual.DockerClient().StartContainer(p.container.ID, nil)
}

func (p *Proxy) Stop() error {
	return virtual.DockerClient().StopContainer(p.container.ID, 0)
}

func (p *Proxy) GetContainer() *docker.Container {
	return p.container
}
