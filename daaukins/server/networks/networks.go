package networks

import (
	"fmt"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/google/uuid"

	"github.com/andreaswachs/bachelors-project/daaukins/server/challenge"
)

type Network struct {
	network *docker.Network
	client  *docker.Client
	subnet  string
}

type ProvisionNetworkOptions struct {
	Subnet string
}

func Provision(client *docker.Client, conf ProvisionNetworkOptions) (*Network, error) {
	subnet, err := getIPBank().GetUnusedSubnet()
	if err != nil {
		return nil, err
	}

	return &Network{
		client: client,
		subnet: subnet}, nil
}

func (n *Network) Create() error {
	if n.network != nil {
		return fmt.Errorf("network is already created")
	}

	name := fmt.Sprintf("daaukins-%s", uuid.New().String())

	network, err := n.client.CreateNetwork(docker.CreateNetworkOptions{
		Name:   name,
		Driver: "macvlan",
		IPAM: &docker.IPAMOptions{
			Config: []docker.IPAMConfig{{
				Subnet: n.subnet,
			}}},
		Labels: map[string]string{
			"daaukins": "network",
		},
	})

	n.network = network

	if err != nil {
		return err
	}

	return nil
}

func (n *Network) Remove() error {
	err := n.client.RemoveNetwork(n.network.ID)
	if err != nil {
		return err
	}

	return nil
}

func (n *Network) Connect(challenge *challenge.Challenge) error {
	if challenge == nil {
		return fmt.Errorf("challenge is nil")
	}

	containerIP, err := getIPBank().GetFreeIP(n.subnet)
	if err != nil {
		return err
	}

	err = n.client.ConnectNetwork(n.network.ID, docker.NetworkConnectionOptions{
		Container: challenge.GetContainerID(),
		EndpointConfig: &docker.EndpointConfig{
			IPAMConfig: &docker.EndpointIPAMConfig{
				IPv4Address: containerIP,
			},
			IPAddress: containerIP,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (n *Network) Disconnect(challenge *challenge.Challenge) error {
	if challenge == nil {
		return fmt.Errorf("challenge is nil")
	}

	err := n.client.DisconnectNetwork(n.network.ID, docker.NetworkConnectionOptions{
		Container: challenge.GetContainerID(),
	})

	if err != nil {
		return err
	}

	return nil
}

// GetNetworkID returns the ID of the isolated docker network
func (n *Network) GetNetworkID() string {
	return n.network.ID
}

// GetDNSAddr returns the IP address of the DNS server in the isolated network
func (n *Network) GetDNSAddr() string {
	subnetOctets := strings.Split(n.subnet, ".")

	return fmt.Sprintf("%s.%s.%s.3", subnetOctets[0], subnetOctets[1], subnetOctets[2])
}
