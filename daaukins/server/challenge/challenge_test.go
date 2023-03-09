package challenge

import (
	"testing"

	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
	docker "github.com/fsouza/go-dockerclient"
)

func TestProvisionChallengeCanProvisionChallenge(t *testing.T) {
	configuration := &ProvisionChallengeOptions{
		Image: "alpine",
		DNSServers: []string{
			"8.8.8.8",
		}}

	_, err := Provision(configuration)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartChallengeCanStartProvisionedChallenge(t *testing.T) {
	configuration := &ProvisionChallengeOptions{
		Image: "alpine",
		DNSServers: []string{
			"8.8.8.8",
		}}

	challenge, err := Provision(configuration)
	if err != nil {
		t.Fatal(err)
	}

	defer challenge.Remove()

	err = challenge.Start()
	if err != nil {
		t.Fatal(err)
	}

	// Check to see if the docker container is running
	_, err = virtual.DockerClient().InspectContainerWithOptions(
		docker.InspectContainerOptions{
			ID: challenge.container.ID,
		})
	if err != nil {
		t.Fatal(err)
	}

}
