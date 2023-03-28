// Package that handles configuration of the client

package config

import (
	"fmt"
)

var (
	config *Config
)

func Load() error {
	configBuffer, err := loadFile(DknConfigFile())
	if err != nil {
		fmt.Println("Could not load configuration file. Please run 'dkn config init' to create a new config file")
		return err
	}

	config = configBuffer
	return nil
}

func Initialize(force bool) error {
	wasWritten, err := CreateConfigFile(force)
	if err != nil {
		return err
	}

	if !wasWritten {
		fmt.Println("Config file already exists. Use --force to overwrite.")
		return fmt.Errorf("config file already exists")
	}

	fmt.Println("Config file created at", DknConfigFile())
	fmt.Println("Please edit the file and add your Daaukins server address and port")
	return nil
}

func ServerAddress() string {
	return config.Server.Address
}

func ServerPort() string {
	return config.Server.Port
}
