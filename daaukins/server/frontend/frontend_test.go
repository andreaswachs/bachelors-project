package frontend

import (
	"testing"

	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
	docker "github.com/fsouza/go-dockerclient"
)

func TestFrontendCanProvision(t *testing.T) {
	fe, err := Provision(&ProvisionFrontendOptions{
		Network: "bridge",
		Ip:      "",
		Port:    3000,
	})

	if err != nil {
		t.Error(err)
	}

	if fe == nil {
		t.Error("frontend is nil")
	}

	if fe.container == nil {
		t.Error("frontend container is nil")
	}

	defer virtual.DockerClient().RemoveContainer(docker.RemoveContainerOptions{
		ID: fe.container.ID,
	})

	if fe.container.ID == "" {
		t.Error("frontend container ID is empty")
	}

	if fe.container.Name == "" {
		t.Error("frontend container name is empty")
	}
}

func TestFrontendPortEmpty(t *testing.T) {
	_, err := Provision(&ProvisionFrontendOptions{
		Network: "bridge",
		Ip:      "",
	})

	if err == nil {
		t.Error("expected error")
	}
}

func TestFrontendStart(t *testing.T) {
	fe, err := Provision(&ProvisionFrontendOptions{
		Network: "bridge",
		Port:    3000,
	})

	if err != nil {
		t.Error(err)
	}

	if err := fe.Start(); err != nil {
		t.Error(err)
	}

	defer virtual.DockerClient().StopContainer(fe.container.ID, 2)

	_, err = virtual.DockerClient().InspectContainerWithOptions(docker.InspectContainerOptions{
		ID: fe.container.ID,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestFrontendStop(t *testing.T) {
	fe, err := Provision(&ProvisionFrontendOptions{
		Network: "bridge",
		Port:    3000,
	})

	if err != nil {
		t.Error(err)
	}

	if err := fe.Start(); err != nil {
		t.Error(err)
	}

	// In case stop fails, we still want to remove the container
	// defer virtual.DockerClient().StopContainer(fe.container.ID, 2)

	if err := fe.Stop(); err != nil {
		t.Error(err)
	}

	container, err := virtual.DockerClient().InspectContainerWithOptions(docker.InspectContainerOptions{
		ID: fe.container.ID,
	})

	if err != nil {
		t.Error(err)
	}

	if container.State.Running {
		t.Errorf("expected container not running, got: %+v", container)
	}
}
