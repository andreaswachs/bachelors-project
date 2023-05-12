// This package relates to the configuration of the running server instance

package config

import (
	"fmt"
	"os"

	"github.com/andreaswachs/bachelors-project/daaukins/server/utils"
	"gopkg.in/yaml.v3"
)

var (
	config      Config
	serverId    string
	confOptsObj Configure = &defaultConfigOptions{}
)

type FollowerConfig struct {
	Name    string `yaml:"name"` // TODO: use the name for something useful
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type ContainerConfig struct {
	Image string `yaml:"image"`
}
type DockerConfig struct {
	Frontend ContainerConfig `yaml:"frontend"`
	Proxy    ContainerConfig `yaml:"proxy"`
	Dhcp     ContainerConfig `yaml:"dhcp"`
	Dns      ContainerConfig `yaml:"dns"`
}

type Config struct {
	ServerMode  Mode             `yaml:"server_mode"`
	ServicePort int              `yaml:"service_port"`
	Docker      DockerConfig     `yaml:"docker"`
	Followers   []FollowerConfig `yaml:"followers"`
}

type Configure interface {
	GetConfigFilename() string
	UsingDockerCompose() bool
	NewServerID() string
}

type defaultConfigOptions struct{}

func (d *defaultConfigOptions) GetConfigFilename() string {
	filename := os.Getenv("DAAUKINS_SERVER_CONFIG")
	if filename == "" {
		return "server.yaml"
	}

	return filename
}

func (d *defaultConfigOptions) UsingDockerCompose() bool {
	return os.Getenv("DAAUKINS_USING_DOCKER_COMPOSE") != ""
}

func (d *defaultConfigOptions) NewServerID() string {
	return utils.RandomShortName()
}

type InitializeConfigOptions struct {
	ConfigFile string
}

// InitializeWith initializes the config with the given configurer.This uses the dependency injection pattern, in that users can provide their own configurer to load the config from a different source.
func InitializeWith(configurer Configure) error {
	confOptsObj = configurer
	return Initialize()
}

func Initialize() error {
	// Load the configuration from the config file
	configBuffer, err := load(confOptsObj.GetConfigFilename())
	if err != nil {
		return err
	}

	config = configBuffer

	// Set the server id
	serverId = confOptsObj.NewServerID()

	return nil
}

func GetFollowers() []FollowerConfig {
	return config.Followers
}

func GetServerMode() Mode {
	return config.ServerMode
}

func GetServerID() string {
	return serverId
}

func GetServicePort() int {
	return config.ServicePort
}

func GetDockerConfig() DockerConfig {
	return config.Docker
}

func IsUsingDockerCompose() bool {
	return confOptsObj.UsingDockerCompose()
}

func load(file string) (Config, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return Config{}, err
	}

	return parse(content)
}

func parse(input []byte) (Config, error) {
	if len(input) == 0 {
		return Config{}, fmt.Errorf("config file is empty")
	}

	var config Config
	err := yaml.Unmarshal(input, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}
