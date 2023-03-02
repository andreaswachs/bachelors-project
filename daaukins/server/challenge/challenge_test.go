package challenge

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
)

func TestProvisionChallengeCanProvisionChallenge(t *testing.T) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	configuration := &ProvisionChallengeOptions{
		Image: "alpine",
		DNSServers: []string{
			"8.8.8.8",
		}}

	_, err = Provision(client, configuration)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartChallengeCanStartProvisionedChallenge(t *testing.T) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	configuration := &ProvisionChallengeOptions{
		Image: "alpine",
		DNSServers: []string{
			"8.8.8.8",
		}}

	challenge, err := Provision(client, configuration)
	if err != nil {
		t.Fatal(err)
	}

	defer challenge.Remove()

	err = challenge.Start()
	if err != nil {
		t.Fatal(err)
	}

	// Check to see if the docker container is running
	_, err = client.InspectContainerWithOptions(
		docker.InspectContainerOptions{
			ID: challenge.container.ID,
		})
	if err != nil {
		t.Fatal(err)
	}

}
