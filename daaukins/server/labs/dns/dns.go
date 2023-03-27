package dns

import (
	"fmt"
	"os"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	"github.com/andreaswachs/bachelors-project/daaukins/server/utils"
	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
	docker "github.com/fsouza/go-dockerclient"
)

var (
	coreFile = `. {
	file zonefile
	errors         # show errors
	log            # enable query logs
}`

	zoneFile = `$ORIGIN .
@   3600 IN SOA sns.dns.icann.org. noc.dns.icann.org. (
	2017042745 ; serial
	7200       ; refresh (2 hours)
	3600       ; retry (1 hour)
	1209600    ; expire (2 weeks)
	3600       ; minimum (1 hour)
	)
`
)

type DNSService struct {
	container      *docker.Container
	filesToCleanup []string
}

type ZoneFileEntry struct {
	Hostname string
	Ip       string
}

type ProvisionDNSOptions struct {
	ZoneFileEntries []ZoneFileEntry
}

func Provision(options *ProvisionDNSOptions) (*DNSService, error) {
	zoneFile := zoneFile
	for _, entry := range options.ZoneFileEntries {
		zoneFile += fmt.Sprintln(toEntry(entry))
	}

	// Create temporary files for coreFile and zoneFile
	tmpCoreFile, err := os.CreateTemp("", "Corefile")
	if err != nil {
		return nil, err
	}

	tmpZoneFile, err := os.CreateTemp("", "zonefile")
	if err != nil {
		return nil, err
	}

	// Write coreFile and zoneFile to temporary files
	if _, err := tmpCoreFile.WriteString(coreFile); err != nil {
		return nil, err
	}

	if _, err := tmpZoneFile.WriteString(zoneFile); err != nil {
		return nil, err
	}

	// Close temporary files
	if err := tmpCoreFile.Close(); err != nil {
		return nil, err
	}

	if err := tmpZoneFile.Close(); err != nil {
		return nil, err
	}

	// Create container
	container, err := virtual.DockerClient().CreateContainer(docker.CreateContainerOptions{
		Name: fmt.Sprintf("daaukins-dns-%s", utils.RandomName()),
		Config: &docker.Config{
			Image:  config.GetDockerConfig().Dns.Image,
			Memory: 128 * 1024 * 1024,
			Labels: map[string]string{
				"daaukins":         "true",
				"daaukins-service": "dns",
			},
		},
		HostConfig: &docker.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:/Corefile", tmpCoreFile.Name()),
				fmt.Sprintf("%s:/zonefile", tmpZoneFile.Name()),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return &DNSService{
		container:      container,
		filesToCleanup: []string{tmpCoreFile.Name(), tmpZoneFile.Name()},
	}, nil
}

func (s *DNSService) Start() error {
	return virtual.DockerClient().StartContainer(s.container.ID, nil)
}

func (s *DNSService) Stop() error {
	return virtual.DockerClient().StopContainer(s.container.ID, 0)
}

func (s *DNSService) GetContainer() *docker.Container {
	return s.container
}

func (s *DNSService) Cleanup() error {
	return utils.DeleteFiles(s.filesToCleanup)
}

func toEntry(entry ZoneFileEntry) string {
	return fmt.Sprintf("%s IN A %s", entry.Hostname, entry.Ip)
}
