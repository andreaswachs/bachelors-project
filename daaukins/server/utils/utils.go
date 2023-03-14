package utils

import (
	"encoding/binary"
	"fmt"
	"os"
	"regexp"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	ErrorMemAvailableNotFoundInMeminfo = fmt.Errorf("could not find MemAvailable in /proc/meminfo")
	meminfoFile                        = "/proc/meminfo"
	MemAvailableRegex                  = regexp.MustCompile(`(?m)^MemAvailable:\s+(\d+) kB$`)
)

func RandomName() string {
	return uuid.New().String()
}

// GetAvailableMemory returns the amount of free memory in megabytes.
// This implementation assumes that the OS is Linux, with the precense of the
// /proc/meminfo file.
func GetAvailableMemory() (uint64, error) {
	contents, err := os.ReadFile(meminfoFile)
	if err != nil {
		return 0, err
	}

	var memFree uint64

	// Find the MemAvailable line
	matches := MemAvailableRegex.FindSubmatch(contents)
	log.Debug().Msgf("matches: %v", matches)
	if len(matches) == 2 {
		// Convert the MemAvailable value to an integer
		mem, n := binary.Uvarint(matches[1])
		if n != len(matches[1]) {
			return 0, fmt.Errorf("could not parse MemAvailable value: %s", matches[1])
		}

		memFree = mem
	}

	if memFree == 0 {
		return 0, fmt.Errorf("could not find MemAvailable in /proc/meminfo")
	}

	// Convert from kilobytes to megabytes
	memFree /= 1024

	return memFree, nil
}

// SetMeminfoFile sets the path to the meminfo file to be used by the
// GetAvailableMemory function.
// This should only be used for testing purposes.
func SetMeminfoFile(file string) {
	meminfoFile = file
}
