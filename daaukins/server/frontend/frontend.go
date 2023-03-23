package frontend

import (
	"fmt"

	"github.com/andreaswachs/bachelors-project/daaukins/server/utils"
	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/rs/zerolog/log"
)

var (
	ErrorPortEmpty = fmt.Errorf("port is empty")
	ErrorIpEmpty   = fmt.Errorf("ip is empty")
	ErrorNetwork   = fmt.Errorf("network is empty")
)

type T interface {
	Start() error
	Stop() error
	GetContainer() *docker.Container
}

type frontend struct {
	container *docker.Container
}

type ProvisionFrontendOptions struct {
	Network string
	DNS     []string
	Ip      string
	Port    int
}

func Provision(options *ProvisionFrontendOptions) (*frontend, error) {
	if err := validateProvisionFrontendOptions(options); err != nil {
		return nil, err
	}

	guestPort := fmt.Sprintf("%d/tcp", options.Port)

	container, err := virtual.DockerClient().CreateContainer(docker.CreateContainerOptions{
		Name: utils.RandomName(),
		Config: &docker.Config{
			Image: "lscr.io/linuxserver/webtop:ubuntu-xfce",
			Labels: map[string]string{
				"daaukins": "true",
			},
			Memory: 2 * 1024 * 1024 * 1024, // 2 GB
		},
		HostConfig: &docker.HostConfig{
			NetworkMode: options.Network,
			DNS:         options.DNS,
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port(guestPort): {
					{
						HostIP:   options.Ip,
						HostPort: fmt.Sprint(options.Port),
					},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return &frontend{
		container: container,
	}, nil
}

func (f *frontend) Start() error {
	if f.container == nil {
		return fmt.Errorf("container is nil")
	}
	return virtual.DockerClient().StartContainer(f.container.ID, &docker.HostConfig{})
}

func (f *frontend) Stop() error {
	force := false

	err := virtual.DockerClient().StopContainer(f.container.ID, 0)
	if err != nil {
		log.Error().Err(err).Msg("failed to stop frontend container")
		force = true
	}

	return virtual.DockerClient().RemoveContainer(docker.RemoveContainerOptions{
		ID:    f.container.ID,
		Force: force,
	})
}

func (f *frontend) GetContainer() *docker.Container {
	return f.container
}

// validateProvisionFrontendOptions validates the options for provisioning a frontend.
// It does not check the IP as that can be empty
func validateProvisionFrontendOptions(options *ProvisionFrontendOptions) error {
	if options.Network == "" {
		return ErrorNetwork
	}

	if options.Port == 0 {
		return ErrorPortEmpty
	}

	return nil
}
