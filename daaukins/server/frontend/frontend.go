package frontend

import (
	"fmt"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	"github.com/andreaswachs/bachelors-project/daaukins/server/utils"
	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
	"github.com/andreaswachs/sizes"
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
	GetProxyPort() int
}

type frontend struct {
	container *docker.Container
	proxyPort int
}

type ProvisionFrontendOptions struct {
	Network   string
	DNS       []string
	Ip        string
	ProxyPort int
	LabID     string
}

func Provision(options *ProvisionFrontendOptions) (*frontend, error) {
	if err := validateProvisionFrontendOptions(options); err != nil {
		return nil, err
	}

	container, err := virtual.DockerClient().CreateContainer(docker.CreateContainerOptions{
		Name: fmt.Sprintf("daaukins-frontend-%s", utils.RandomName()),
		Config: &docker.Config{
			Image: config.GetDockerConfig().Frontend.Image,
			Labels: map[string]string{
				"daaukins":     "frontend",
				"daaukins.lab": options.LabID,
			},
			Memory: sizes.Gigabytes[int64](2),
			Tty:    true,
		},
		HostConfig: &docker.HostConfig{
			NetworkMode: options.Network,
			DNS:         options.DNS,
		},
	})
	if err != nil {
		return nil, err
	}

	return &frontend{
		container: container,
		proxyPort: options.ProxyPort,
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

func (f *frontend) GetProxyPort() int {
	return f.proxyPort
}

// validateProvisionFrontendOptions validates the options for provisioning a frontend.
// It does not check the IP as that can be empty
func validateProvisionFrontendOptions(options *ProvisionFrontendOptions) error {
	if options.Network == "" {
		return ErrorNetwork
	}

	if options.ProxyPort == 0 {
		return ErrorPortEmpty
	}

	return nil
}
