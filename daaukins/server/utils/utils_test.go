package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestGetAvailableMemoryHasExactMemory(t *testing.T) {
	cleanup, err := createMockMeminfoFile(t, 49)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	// Get the available memory
	availableMemory, err := GetAvailableMemory()
	if err != nil {
		t.Error(err)
	}

	// Check that the available memory is less than 50 MB
	if availableMemory >= 50 {
		t.Error("availableMemory is greater than 50 MB")
	}

	// Check that the available memory is 49
	if availableMemory != 49 {
		t.Error("availableMemory is not 49 MB")
	}
}

func TestGetAvailableMemoryHasZeroMemory(t *testing.T) {
	cleanup, err := createMockMeminfoFile(t, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	// Get the available memory
	availableMemory, err := GetAvailableMemory()
	if err != nil {
		t.Error(err)
	}

	if availableMemory != 0 {
		t.Errorf("availableMemory is not 0 MB, was %d MB", availableMemory)
	}
}

func createMockMeminfoFile(t *testing.T, memAvailableMb uint64) (func(), error) {
	meminfoContents := fmt.Sprintf(`MemTotal:       16270828 kB
MemFree:         3745292 kB
MemAvailable:    %d kB
Buffers:            3896 kB
Cached:          5039020 kB
SwapCached:            0 kB
Active:          5644404 kB
Inactive:        3942080 kB
Active(anon):    3743676 kB`, memAvailableMb*1024)

	// Create a mock meminfo file
	meminfo, err := os.CreateTemp("", "meminfo")
	if err != nil {
		t.Error(err)
	}

	// Write the mock meminfo file
	_, err = meminfo.Write([]byte(meminfoContents))
	if err != nil {
		t.Error(err)
	}

	// Set the meminfo file to be used by the GetAvailableMemory function
	meminfoFile = meminfo.Name()

	return func() {
		os.Remove(meminfo.Name())
		meminfo.Close()
	}, nil
}
