package labs

import (
	"fmt"
	"os"

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

type zoneFileEntry struct {
	hostname string
	ip       string
}

type provisionDNSOptions struct {
	zoneFileEntries []zoneFileEntry
}

func provisionDNS(options *provisionDNSOptions) (*networkService, error) {
	zoneFile := zoneFile
	for _, entry := range options.zoneFileEntries {
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
		Config: &docker.Config{
			Image:  "coredns/coredns:1.10.0",
			Memory: 128 * 1024 * 1024,
			Labels: map[string]string{
				"daaukins": "true",
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

	return &networkService{
		container: container,
		filesToCleanup: []string{
			tmpCoreFile.Name(),
			tmpZoneFile.Name(),
		},
	}, nil
}

func toEntry(entry zoneFileEntry) string {
	return fmt.Sprintf("%s IN A %s", entry.hostname, entry.ip)
}
