package labs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/andreaswachs/bachelors-project/daaukins/server/store"
	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
	docker "github.com/fsouza/go-dockerclient"
)

type preppedYamlConfig uint8

const (
	goodYaml preppedYamlConfig = iota
	badYaml
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func TestLoad(t *testing.T) {
	// Prep config file
	fileName, err := prepYamlConfigFile(goodYaml, t)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fileName)

	// Load the lab
	lab, err := load(fileName)
	if err != nil {
		t.Fatal(err)
	}

	// Check the lab name
	if lab.Name != "test" {
		t.Fatalf("Expected lab name to be 'test', got '%s'", lab.Name)
	}

	// Check the challenge name
	if lab.Challenges[0].Name != "test" {
		t.Fatalf("Expected challenge name to be 'test', got '%s'", lab.Challenges[0].Name)
	}

	// Check the challenge ID
	if lab.Challenges[0].Challenge != "test" {
		t.Fatalf("Expected challenge ID to be 'test', got '%s'", lab.Challenges[0].Challenge)
	}

	// Check the challenge DNS
	if lab.Challenges[0].Dns[0] != "test" {
		t.Fatalf("Expected challenge DNS to be 'test', got '%s'", lab.Challenges[0].Dns[0])
	}
}

func TestLoadGivenBadPath(t *testing.T) {
	// Load the lab
	_, err := load("bad/path")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestLoadGivenBadYaml(t *testing.T) {
	fileName, err := prepYamlConfigFile(badYaml, t)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fileName)

	// Load the lab
	_, err = load(fileName)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// TODO: refactor the store and the lab configuration to enable use of mocking such that we can test the lab provisioning
// func TestProvision(t *testing.T) {
// 	// Prep config file
// 	fileName, err := prepYamlConfigFile(goodYaml, t)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.Remove(fileName)
//
// 	// Provision the lab
// 	lab, err := Provision(fileName)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	// Check the lab name
// 	if lab.GetName() != "test" {
// 		t.Fatalf("Expected lab name to be 'test', got '%s'", lab.GetName())
// 	}
//
// 	// Check the amount of challenges
// 	if len(lab.GetChallenges()) != 1 {
// 		t.Fatalf("Expected 1 challenge, got %d", len(lab.GetChallenges()))
// 	}
//
// 	if lab.GetChallenges()[0].GetDNS()[0] != "test" {
// 		t.Fatalf("Expected challenge name to be 'test', got '%s'", lab.GetChallenges()[0].GetDNS()[0])
// 	}
//
// }

func TestStart(t *testing.T) {
	// Loads the test store
	if err := store.Load(getTestResource("store.yaml")); err != nil {
		t.Fatal(err)
	}
	// Load the test lab
	lab, err := Provision(getTestResource("lab.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	// Start the lab
	// defer lab.Remove()

	err = lab.Start()
	if err != nil {
		t.Fatal(err)
	}

	if err = verifyContainerRunning(lab.dhcpService.container.ID); err != nil {
		t.Fatal(err)
	}

	if err = verifyContainerRunning(lab.dnsService.container.ID); err != nil {
		t.Fatal(err)
	}

	for _, challenge := range lab.challenges {
		if err = verifyContainerRunning(challenge.GetContainerID()); err != nil {
			t.Fatal(err)
		}
	}
}

func prepYamlConfigFile(yamlSetting preppedYamlConfig, t *testing.T) (string, error) {
	var yamlConfig string

	if yamlSetting == goodYaml {
		yamlConfig = `name: "test"
challenges:
  - name: "test"
    challenge: "test"
    dns:
      - "test"`
	} else {
		yamlConfig = `name: "test"
	aaaaaaaaaaaaaaaaaaaaah`
	}

	// Write to a temporary file
	file, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := file.Write([]byte(yamlConfig)); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	return file.Name(), nil
}

func getTestResource(name string) string {
	return filepath.Join(basepath, "..", "test_resources", name)
}

func verifyContainerRunning(ID string) error {
	container, err := virtual.DockerClient().InspectContainerWithOptions(
		docker.InspectContainerOptions{
			ID: ID,
		})

	if err != nil {
		return err
	}

	if !container.State.Running {
		return fmt.Errorf("container %s is not running", ID)
	}

	return nil
}
