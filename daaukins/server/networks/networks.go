package networks

import (
	"fmt"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/google/uuid"

	"github.com/andreaswachs/bachelors-project/daaukins/server/challenge"
	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
)

type Network struct {
	network *docker.Network
	subnet  string
	name    string
}

type ProvisionNetworkOptions struct {
	Subnet string
}

func Provision(conf ProvisionNetworkOptions) (*Network, error) {
	subnet, err := ipPool().GetUnusedSubnet()
	if err != nil {
		return nil, err
	}

	return &Network{
		subnet: subnet}, nil
}

func (n *Network) Create() error {
	if n.network != nil {
		return fmt.Errorf("network is already created")
	}

	name := fmt.Sprintf("daaukins-%s", uuid.New().String())
	n.name = name

	network, err := virtual.DockerClient().CreateNetwork(docker.CreateNetworkOptions{
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
	err := virtual.DockerClient().RemoveNetwork(n.network.ID)
	if err != nil {
		return err
	}

	return nil
}

func (n *Network) Connect(challenge *challenge.Challenge) error {
	containerIP, err := ipPool().GetFreeIP(n.subnet)
	if err != nil {
		return err
	}

	challenge.SetIP(containerIP)

	return n.connectContainer(challenge.GetContainerID(), containerIP)
}

func (n *Network) ConnectDNS(container *docker.Container) error {
	return n.connectContainer(container.ID, n.GetDNSAddr())
}

func (n *Network) ConnectDHCP(container *docker.Container) error {
	return n.connectContainer(container.ID, n.GetDHCPAddr())
}

func (n *Network) Disconnect(challenge *challenge.Challenge) error {
	if challenge == nil {
		return fmt.Errorf("challenge is nil")
	}

	err := virtual.DockerClient().DisconnectNetwork(n.network.ID, docker.NetworkConnectionOptions{
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
	return addrInSubnet(n.subnet, "3")
}

// GetDHCPAddr returns the IP address of the DHCP server in the isolated network
func (n *Network) GetDHCPAddr() string {
	return addrInSubnet(n.subnet, "2")

}

// GetSubnet returns the subnet of the isolated network
func (n *Network) GetSubnet() string {
	return n.subnet
}

func (n *Network) GetName() string {
	return n.name
}

func (n *Network) connectContainer(containerID string, containerIP string) error {
	err := virtual.DockerClient().ConnectNetwork(n.network.ID, docker.NetworkConnectionOptions{
		Container: containerID,
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

func addrInSubnet(subnet, octet string) string {
	subnetOctets := strings.Split(subnet, ".")

	return fmt.Sprintf("%s.%s.%s.%s", subnetOctets[0], subnetOctets[1], subnetOctets[2], octet)
}
