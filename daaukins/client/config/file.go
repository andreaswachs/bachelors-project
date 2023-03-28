// File to work with the dkn client configuration file

package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrFileNotFound = fmt.Errorf("config file not found. Please run 'dkn config init' to create a new config file")
)

type Server struct {
	// The server address
	Address string `yaml:"address"`
	// The server port
	Port string `yaml:"port"`
}

type Config struct {
	// The server configuration
	Server Server `yaml:"server"`
}

func CreateConfigFile(force bool) (bool, error) {
	_, err := os.Stat(DknConfigFile())
	if os.IsNotExist(err) || force {
		defaultConfig := &Config{
			Server: Server{
				Address: "localhost",
				Port:    "8080",
			},
		}

		bytes, err := yaml.Marshal(defaultConfig)
		if err != nil {
			return false, err
		}

		if err = ensureFolderPathExists(); err != nil {
			return false, err
		}

		return true, os.WriteFile(DknConfigFile(), bytes, os.ModePerm)
	}

	return false, err
}

func DknConfigFile() string {
	return DknBasePath() + "/config.yaml"
}

func loadFile(path string) (*Config, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println("Config file not found", err)
		return nil, ErrFileNotFound
	}
	if err != nil {
		log.Println("Could not stat config file", err)
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		log.Println("Could not read config file", err)
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(content, config); err != nil {
		log.Println("Could not unmarshal config file", err)
		return nil, err
	}

	return config, nil
}
