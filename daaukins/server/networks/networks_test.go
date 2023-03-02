package networks

import (
	"testing"

	"github.com/andreaswachs/bachelors-project/daaukins/server/challenge"
	docker "github.com/fsouza/go-dockerclient"
)

func TestCreateCanCreateNetwork(t *testing.T) {
	// Create a new docker client
	client, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	config, err := newProvisionNetworkOptions()
	if err != nil {
		t.Fatal(err)
	}

	// Provision the new network
	network, err := Provision(client, *config)
	if err != nil {
		t.Fatal(err)
	}

	err = network.Create()
	if err != nil {
		t.Fatal(err)
	}

	// Check if the network is not nil
	if network == nil {
		t.Fatal("network is nil")
	}

	// Check if the network is created
	_, err = client.NetworkInfo(network.GetNetworkID())
	if err != nil {
		t.Fatal("network is not created")
	}

	// Remove the network
	err = client.RemoveNetwork(network.GetNetworkID())
	if err != nil {
		t.Fatal(err)
	}

	// Check if the network is removed
	_, err = client.NetworkInfo(network.GetNetworkID())
	if err == nil {
		t.Fatal("network is not removed. You need to do this manually, sorry")
	}
}

// Testing the Remove function
func TestRemoveCanRemoveCreatedNetwork(t *testing.T) {
	// Create a new docker client
	client, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	config, err := newProvisionNetworkOptions()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new network
	network, err := Provision(client, *config)
	if err != nil {
		t.Fatal(err)
	}

	err = network.Create()
	if err != nil {
		t.Fatal(err)
	}

	// Check if the network is not nil
	if network == nil {
		t.Fatal("network is nil")
	}

	// Remove the network
	err = network.Remove()
	if err != nil {
		t.Fatal(err)
	}

	// Check if the network is removed
	_, err = client.NetworkInfo(network.GetNetworkID())
	if err == nil {
		t.Fatal("network is not removed")
	}
}

func TestConnectCanConnectContainerToNetwork(t *testing.T) {
	// Create a new docker client
	client, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	config, err := newProvisionNetworkOptions()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new network
	network, err := Provision(client, *config)
	if err != nil {
		t.Fatal(err)
	}

	if err = network.Create(); err != nil {
		t.Fatal(err)
	}

	defer client.RemoveNetwork(network.GetNetworkID())

	// Check if the network is not nil
	if network == nil {
		t.Fatal("network is nil")
	}

	challengeConf := &challenge.ProvisionChallengeOptions{
		Image: "alpine",
		DNSServers: []string{
			"8.8.8.8",
		},
	}

	// Create configuration for a new challenge
	challenge, err := challenge.Provision(client, challengeConf)
	if err != nil {
		t.Fatal(err)
	}

	challenge.Start()
	defer challenge.Remove()

	// Check if the challenge is not nil
	if challenge == nil {
		t.Fatal("challenge is nil")
	}

	// Connect the challenge to the network
	err = network.Connect(challenge)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the challenge is connected to the network
	inspectedContainer, err := client.InspectContainerWithOptions(
		docker.InspectContainerOptions{
			ID: challenge.GetContainerID(),
		})
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := inspectedContainer.NetworkSettings.Networks[network.GetNetworkID()]; ok {
		t.Fatal("challenge is not connected to the network")
	}

	// Disconnect the challenge from the network
	err = client.DisconnectNetwork(network.GetNetworkID(), docker.NetworkConnectionOptions{
		Container: challenge.GetContainerID(),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func newProvisionNetworkOptions() (*ProvisionNetworkOptions, error) {
	subnet, err := getIPBank().GetUnusedSubnet()
	if err != nil {
		return nil, err
	}

	return &ProvisionNetworkOptions{
		Subnet: subnet,
	}, nil
}
