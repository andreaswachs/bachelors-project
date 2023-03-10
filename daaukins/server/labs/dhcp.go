package labs

import (
	"fmt"
	"os"
	"strings"

	"github.com/andreaswachs/bachelors-project/daaukins/server/utils"
	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"

	docker "github.com/fsouza/go-dockerclient"
)

type provisionDHCPOptions struct {
	DNSAddr     string
	Subnet      string
	NetworkMode string
}

func format(subnet string, lastOctet int) string {
	octets := strings.Split(subnet, ".")

	return fmt.Sprintf("%s.%d", strings.Join(octets[:len(octets)-1], "."), lastOctet)
}

func provisionDHCP(options *provisionDHCPOptions) (*networkService, error) {
	if err := validateProvisionDHCPOptions(options); err != nil {
		return nil, err
	}

	confStr := fmt.Sprintf(
		`option domain-name-servers %s;

	subnet %s netmask 255.255.255.0 {
		range %s %s;
		option subnet-mask 255.255.255.0;
		option broadcast-address %s;
		option routers %s;
	}`, options.DNSAddr,
		options.Subnet[:len(options.Subnet)-3], // Remove the '/24' from the ip
		format(options.Subnet, 4),
		format(options.Subnet, 254),
		format(options.Subnet, 255),
		format(options.Subnet, 1))

	// Write dhcpConfig to a temporary file
	dhcpConfigFile, err := os.CreateTemp("", "dhcpd.conf")
	if err != nil {
		return nil, err
	}
	dhcpConfigFile.Write([]byte(confStr))

	container, err := virtual.DockerClient().CreateContainer(docker.CreateContainerOptions{
		Name: utils.RandomName(),
		Config: &docker.Config{
			Image:  "networkboot/dhcpd:1.2.0",
			Memory: 128 * 1024 * 1024,
			Labels: map[string]string{
				"daaukins": "true",
			},
			Cmd: []string{
				"eth0",
			},
		},
		HostConfig: &docker.HostConfig{
			DNS: []string{options.DNSAddr},
			Binds: []string{
				fmt.Sprintf("%s:/data/dhcpd.conf", dhcpConfigFile.Name()),
			},
			NetworkMode: options.NetworkMode,
		},
	})

	if err != nil {
		return nil, err
	}

	return &networkService{
		container:      container,
		filesToCleanup: []string{dhcpConfigFile.Name()},
	}, nil
}

func validateProvisionDHCPOptions(options *provisionDHCPOptions) error {
	if options.DNSAddr == "" {
		return fmt.Errorf("dns address is missing")
	}

	if options.Subnet == "" {
		return fmt.Errorf("subnet is missing")
	}

	return nil
}
