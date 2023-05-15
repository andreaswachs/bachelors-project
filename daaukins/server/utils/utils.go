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
	memAvailableRegex                  = regexp.MustCompile(`(?m)^MemAvailable\:\s*(\d+)\s+kB\s*$`)
)

func RandomName() string {
	return uuid.New().String()
}

func RandomShortName() string {
	return uuid.New().String()[:8]
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

	// using an int for memFree means that available system memory may not exceed
	// roughly ~2.1TB, as that would produce wrong calculations (int32 max value = 2^31-1),
	// the unit is kb and thus the total max value translates to 2.1TB
	var memFree int

	// Find the MemAvailable line
	matches := memAvailableRegex.FindSubmatch(contents)
	// Matches are the full match and the first capture group
	if len(matches) != 2 {
		return 0, ErrorMemAvailableNotFoundInMeminfo
	}

	// Convert the MemAvailable value to an integer
	// We can ignore the error because the regex only matches numbers
	memFree, _ = strconv.Atoi(string(matches[1]))

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
