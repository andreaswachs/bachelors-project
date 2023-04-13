package frontend

import (
	"os"
	"testing"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
	docker "github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v3"
)

// type Configure interface {
// 	GetConfigFilename() string
// 	UsingDockerCompose() bool
// 	NewServerID() string
// }

type testConfig struct {
	filename string
}

// Implement the Configure interface
func (t *testConfig) GetConfigFilename() string {
	return t.filename
}

func (t *testConfig) UsingDockerCompose() bool {
	return false
}

func (t *testConfig) NewServerID() string {
	return "test"
}

func init() {
	// Initialize the docker client
	if err := virtual.Initialize(); err != nil {
		panic(err)
	}

	// Create a new config file
	file, err := os.CreateTemp("", "server.yaml")
	if err != nil {
		panic(err)
	}

	// Create the config object
	confObj := &config.Config{
		ServerMode:  config.ModeLeader,
		ServicePort: 8080,
		Docker: config.DockerConfig{
			Frontend: config.ContainerConfig{
				Image: "alpine",
			},
			Proxy: config.ContainerConfig{
				Image: "alpine",
			},
			Dhcp: config.ContainerConfig{
				Image: "alpine",
			},
			Dns: config.ContainerConfig{
				Image: "alpine",
			},
		},
		Followers: []config.FollowerConfig{
			{
				Name:    "follower1",
				Address: "localhost",
				Port:    8080,
			},
		},
	}
	asText, err := yaml.Marshal(confObj)
	if err != nil {
		panic(err)
	}

	// Write the config to the file
	if _, err := file.Write(asText); err != nil {
		panic(err)
	}

	defer file.Close()
	defer os.Remove(file.Name())

	if err := config.InitializeWith(&testConfig{filename: file.Name()}); err != nil {
		panic(err)
	}
}

func TestFrontendCanProvision(t *testing.T) {
	fe, err := Provision(&ProvisionFrontendOptions{
		Network:   "bridge",
		Ip:        "",
		ProxyPort: 3000,
	})

	if err != nil {
		t.Fatal(err)
	}

	if fe == nil {
		t.Fatal("frontend is nil")
	}

	if fe.container == nil {
		t.Fatal("frontend container is nil")
	}

	defer virtual.DockerClient().RemoveContainer(docker.RemoveContainerOptions{
		ID: fe.container.ID,
	})

	if fe.container.ID == "" {
		t.Fatal("frontend container ID is empty")
	}

	if fe.container.Name == "" {
		t.Fatal("frontend container name is empty")
	}
}

func TestFrontendPortEmpty(t *testing.T) {
	_, err := Provision(&ProvisionFrontendOptions{
		Network: "bridge",
		Ip:      "",
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFrontendStart(t *testing.T) {
	fe, err := Provision(&ProvisionFrontendOptions{
		Network:   "bridge",
		ProxyPort: 3000,
	})

	if err != nil {
		t.Fatal(err)
	}

	if err := fe.Start(); err != nil {
		t.Fatal(err)
	}

	defer virtual.DockerClient().StopContainer(fe.container.ID, 2)

	_, err = virtual.DockerClient().InspectContainerWithOptions(docker.InspectContainerOptions{
		ID: fe.container.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestFrontendStop(t *testing.T) {
	fe, err := Provision(&ProvisionFrontendOptions{
		Network:   "bridge",
		ProxyPort: 3000,
	})

	if err != nil {
		t.Fatal(err)
	}

	if err := fe.Start(); err != nil {
		t.Fatal(err)
	}

	defer virtual.DockerClient().StopContainer(fe.container.ID, 2)

	if err := fe.Stop(); err != nil {
		t.Fatal(err)
	}

	_, err = virtual.DockerClient().InspectContainerWithOptions(docker.InspectContainerOptions{
		ID: fe.container.ID,
	})

	// The container should not exist
	if err == nil {
		t.Fatal(err)
	}

}
