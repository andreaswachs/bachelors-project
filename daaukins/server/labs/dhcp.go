package labs

import (
	"fmt"
	"os"
	"strings"

	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
	docker "github.com/fsouza/go-dockerclient"
)

var (
	dhcpConfig = `option domain-name-servers $dns;
subnet $subnet netmask 255.255.255.0 {
		range $subnetBeginRange $subnetEndRange;
		option subnet-mask 255.255.255.0;
		option broadcast-address $subnetBroadcastAddress;
		option routers $subnetRouterAddress;
}`
)

type provisionDHCPOptions struct {
	DNSAddr string
	Subnet  string
}

func provisionDHCP(options *provisionDHCPOptions) (*networkService, error) {
	if err := validateProvisionDHCPOptions(options); err != nil {
		return nil, err
	}

	dhcpConfig := strings.ReplaceAll(dhcpConfig, "$dns", options.DNSAddr)
	dhcpConfig = strings.ReplaceAll(dhcpConfig, "$subnetBeginRange", getBeginRangeFromSubnet(options.Subnet))
	dhcpConfig = strings.ReplaceAll(dhcpConfig, "$subnetEndRange", getEndRangeFromSubnet(options.Subnet))
	dhcpConfig = strings.ReplaceAll(dhcpConfig, "$subnetBroadcastAddress", getBroadcastAddressFromSubnet(options.Subnet))
	dhcpConfig = strings.ReplaceAll(dhcpConfig, "$subnetRouterAddress", getBroadcastAddressFromSubnet(options.Subnet))
	dhcpConfig = strings.ReplaceAll(dhcpConfig, "$subnet", options.Subnet)

	// Write dhcpConfig to a temporary file
	dhcpConfigFile, err := os.CreateTemp("", "dhcpd.conf")
	if err != nil {
		return nil, err
	}

	container, err := virtual.DockerClient().CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:  "networkboot/dhcpd:1.2.0",
			Memory: 128 * 1024 * 1024,
		},
		HostConfig: &docker.HostConfig{
			DNS: []string{options.DNSAddr},
			Binds: []string{
				fmt.Sprintf("%s:/data/dhcpd.conf", dhcpConfigFile.Name()),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return &networkService{container: container}, nil
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

func getBeginRangeFromSubnet(subnet string) string {
	octets := strings.Split(subnet, ".")

	return fmt.Sprintf("%s.%s.%s.4", octets[0], octets[1], octets[2])
}

func getEndRangeFromSubnet(subnet string) string {
	octets := strings.Split(subnet, ".")

	return fmt.Sprintf("%s.%s.%s.254", octets[0], octets[1], octets[2])
}

func getBroadcastAddressFromSubnet(subnet string) string {
	octets := strings.Split(subnet, ".")

	return fmt.Sprintf("%s.%s.%s.255", octets[0], octets[1], octets[2])
}
