// This package relates to the configuration of the running server instance

package config

import (
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Mode string

var (
	config Config
)

const (
	ModeMinion Mode = "minion"
	ModeLeader Mode = "leader"
)

type MinionConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type Config struct {
	ServerMode Mode           `yaml:"server_mode"`
	Minions    []MinionConfig `yaml:"minions"`
}

func Initialize() {
	// Load the configuration from the config file
	configBuffer, err := load("server.yaml")
	if err != nil {
		log.Panic().Err(err).Msg("Failed to load config file")
	}

	config = configBuffer

	log.Info().Msgf("Loaded config: %+v", config)
}

func GetMinions() []MinionConfig {
	return config.Minions
}

func GetServerMode() Mode {
	return config.ServerMode
}

func load(file string) (Config, error) {
	content, err := ioutil.ReadFile(file)
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
