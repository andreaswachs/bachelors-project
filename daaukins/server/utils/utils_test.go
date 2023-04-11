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

func TestDeleteFiles(t *testing.T) {
	// Create a mock file
	file, err := os.CreateTemp("", "file")
	if err != nil {
		t.Error(err)
	}

	// Delete the mock file
	err = DeleteFiles([]string{file.Name()})
	if err != nil {
		t.Error(err)
	}

	// Check that the file was deleted
	_, err = os.Stat(file.Name())
	if err == nil {
		t.Error("file was not deleted")
	}
}

func TestDeleteFilesNotExists(t *testing.T) {
	// Delete a file that does not exist
	err := DeleteFiles([]string{"/does/not/exist"})
	if err == nil {
		t.Error("DeleteFiles did not return an error")
	}
}

func TestGetAvailableMemoryMissingMeminfoFile(t *testing.T) {
	// Set the meminfo file to be used by the GetAvailableMemory function
	meminfoFile = "/does/not/exist"

	// Get the available memory
	_, err := GetAvailableMemory()
	if err == nil {
		t.Error("GetAvailableMemory did not return an error")
	}
}

func TestGetAvailableMemoryMissingMemAvailable(t *testing.T) {
	// Create a mock meminfo file
	meminfo, err := os.CreateTemp("", "meminfo")
	if err != nil {
		t.Error(err)
	}

	// Write the mock meminfo file
	_, err = meminfo.Write([]byte(`MemTotal:       16270828 kB
MemFree:         3745292 kB
Buffers:            3896 kB
Cached:          5039020 kB
SwapCached:            0 kB
Active:          5644404 kB
Inactive:        3942080 kB
Active(anon):    3743676 kB`))
	if err != nil {
		t.Error(err)
	}

	// Set the meminfo file to be used by the GetAvailableMemory function
	meminfoFile = meminfo.Name()

	// Get the available memory
	_, err = GetAvailableMemory()
	if err == nil {
		t.Error("GetAvailableMemory did not return an error")
	}

	// Close and remove the mock meminfo file
	os.Remove(meminfo.Name())
	meminfo.Close()
}

func TestRandomName(t *testing.T) {
	// Get a random name
	firstName := RandomName()

	// Get another random name
	secondName := RandomName()

	// Check that the names are different
	if firstName == secondName {
		t.Error("RandomName returned the same name twice")
	}
}

func TestRandomShortName(t *testing.T) {
	// Get a random short name
	firstName := RandomShortName()

	// Get another random short name
	secondName := RandomShortName()

	// Check that the names are different
	if firstName == secondName {
		t.Error("GetRandomShortName returned the same name twice")
	}
}

func createMockMeminfoFile(t *testing.T, memAvailableMb int) (func(), error) {
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
