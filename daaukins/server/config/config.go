// This package relates to the configuration of the running server instance

package config

import (
	"os"

	"github.com/andreaswachs/bachelors-project/daaukins/server/utils"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var (
	config   Config
	serverId string
)

type MinionConfig struct {
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
	ServerMode  Mode           `yaml:"server_mode"`
	ServicePort int            `yaml:"service_port"`
	Docker      DockerConfig   `yaml:"docker"`
	Followers   []MinionConfig `yaml:"followers"`
}

type InitializeConfigOptions struct {
	ConfigFile string
}

func Initialize(options *InitializeConfigOptions) {
	configFilename := "server.yaml"

	if options != nil && options.ConfigFile != "" {
		configFilename = options.ConfigFile
	}

	// Load the configuration from the config file
	configBuffer, err := load(configFilename)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to load config file")
	}

	config = configBuffer

	// Set the server id
	serverId = utils.RandomShortName()

	log.Info().Msgf("Loaded config: %+v", config)
}

func GetMinions() []MinionConfig {
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

func load(file string) (Config, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return Config{}, err
	}

	return parse(content)
}

func parse(input []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(input, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}
