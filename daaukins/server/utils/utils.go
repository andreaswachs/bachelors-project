package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/google/uuid"
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
// This will be faulty if the system RAM exeeds 2TB
func GetAvailableMemory() (int, error) {
	contents, err := os.ReadFile(meminfoFile)
	if err != nil {
		return 0, err
	}

	var memFree int

	// Find the MemAvailable line
	matches := MemAvailableRegex.FindSubmatch(contents)
	// Matches are the full match and the first capture group
	if len(matches) == 2 {
		// Convert the MemAvailable value to an integer
		memFree, err = strconv.Atoi(string(matches[1]))
		if err != nil {
			return 0, err
		}
	}

	if memFree == 0 {
		return 0, fmt.Errorf("could not find MemAvailable in /proc/meminfo")
	}

	// Convert from kilobytes to megabytes
	memFree /= 1024

	return memFree, nil
}

func DeleteFiles(files []string) error {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

// SetMeminfoFile sets the path to the meminfo file to be used by the
// GetAvailableMemory function.
// This should only be used for testing purposes.
func SetMeminfoFile(file string) {
	meminfoFile = file
}
