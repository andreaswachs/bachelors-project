package challenge

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/google/uuid"
)

type ProvisionChallengeOptions struct {
	Image       string
	DNSServers  []string
	DNSSettings []string
}

type Challenge struct {
	container                *docker.Container
	client                   *docker.Client
	containerConfiguration   *docker.CreateContainerOptions
	provisionedConfiguration *ProvisionChallengeOptions
	ip                       string
	dnsSettings              []string
}

// Provisions creates and prepares the configuration for the challenge.
// It does not start the challenge container itself. You need to call Start() for that.
func Provision(client *docker.Client, conf *ProvisionChallengeOptions) (*Challenge, error) {
	if err := validateProvisionChallengeOptions(conf); err != nil {
		return nil, err
	}

	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	containerConfiguration := &docker.CreateContainerOptions{
		Name: fmt.Sprintf("daaukins-%s", uuid.New().String()),
		Config: &docker.Config{
			Image: conf.Image,
		},
		HostConfig: &docker.HostConfig{
			DNS: conf.DNSServers,
		},
	}

	return &Challenge{
		client:                   client,
		provisionedConfiguration: conf,
		containerConfiguration:   containerConfiguration,
		dnsSettings:              conf.DNSSettings,
	}, nil
}

// Start starts the challenge container
func (c *Challenge) Start() error {
	if c.container != nil {
		return fmt.Errorf("challenge is already started")
	}

	if err := handleErr(c); err != nil {
		return err
	}

	container, err := c.client.CreateContainer(*c.containerConfiguration)
	if err != nil {
		return err
	}

	c.container = container
	return nil
}

func (c *Challenge) Remove() error {
	if err := handleErr(c); err != nil {
		return err
	}

	err := c.client.RemoveContainer(docker.RemoveContainerOptions{
		ID: c.container.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Challenge) GetContainerID() string {
	if err := handleErr(c); err != nil {
		return ""
	}

	return c.container.ID
}

func (c *Challenge) GetIP() string {
	if err := handleErr(c); err != nil {
		return ""
	}

	return c.ip
}

func (c *Challenge) SetIP(ip string) {
	c.ip = ip
}

func (c *Challenge) GetDNS() []string {
	if c.dnsSettings == nil {
		return []string{}
	}

	return c.dnsSettings
}

func handleErr(c *Challenge) error {
	if c == nil {
		return fmt.Errorf("challenge is nil")
	}

	if c.client == nil {
		return fmt.Errorf("client is nil")
	}

	return nil
}

func validateProvisionChallengeOptions(conf *ProvisionChallengeOptions) error {
	if conf.Image == "" {
		return fmt.Errorf("image is empty")
	}

	if len(conf.DNSServers) == 0 {
		return fmt.Errorf("dns servers is empty")
	}

	return nil
}
