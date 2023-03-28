package config

import (
	"os"
)

const (
	rootConfigDir = "/dkn"
)

func DknBasePath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	return dir + rootConfigDir
}

func ensureFolderPathExists() error {
	return os.MkdirAll(DknBasePath(), os.ModePerm)
}
