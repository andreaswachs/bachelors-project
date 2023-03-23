// Labs are a collection of challenges that are connected to each other over an isolated network.
// The frontend allows users to interact with the services deployed to the network.
// A lab contains a single frontend. This is due to the nature of this being my bachelors project
// and this source code being a PoC.
package labs

import (
	"fmt"
	"os"

	"github.com/andreaswachs/bachelors-project/daaukins/server/challenge"
	"github.com/andreaswachs/bachelors-project/daaukins/server/frontend"
	"github.com/andreaswachs/bachelors-project/daaukins/server/labs/dhcp"
	"github.com/andreaswachs/bachelors-project/daaukins/server/labs/dns"
	"github.com/andreaswachs/bachelors-project/daaukins/server/networks"
	"github.com/andreaswachs/bachelors-project/daaukins/server/store"
	"github.com/andreaswachs/bachelors-project/daaukins/server/utils"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var (
	labs map[string]*lab

	ErrorLabNameMissing     = fmt.Errorf("lab name is missing")
	ErrorLabNoChallenges    = fmt.Errorf("lab has no challenges")
	ErrorChallengeNameEmpty = fmt.Errorf("challenge name is empty")
	ErrorChallengeIDEmpty   = fmt.Errorf("challenge id is empty")
	ErrorChallengeNoDNS     = fmt.Errorf("challenge has no dns servers")
	ErrorLabDoesntExist     = fmt.Errorf("lab does not exist")
)

type networkService interface {
	Start() error
	Stop() error
	GetContainer() *docker.Container
	Cleanup() error
}

// lab is the data bearing struct for a lab
type lab struct {
	name        string
	challenges  []*challenge.Challenge
	dhcpService networkService
	dnsService  networkService
	network     *networks.Network
	isStarted   bool
	frontend    frontend.T
}

type labChallenge struct {
	Name      string   `yaml:"name"`
	Challenge string   `yaml:"challenge"`
	Dns       []string `yaml:"dns"`
}

type labDTO struct {
	Name       string         `yaml:"name"`
	Challenges []labChallenge `yaml:"challenges"`
}

// type networkService struct {
// 	container      *docker.Container
// 	filesToCleanup []string
// }

func init() {
	// Ensure labs is initiated to an empty slice
	labs = make(map[string]*lab, 0)
}

// GetByName returns a lab with a given name. If the lab does not exist, an error is returned
func GetByName(name string) (*lab, error) {
	if labs[name] == nil {
		return nil, fmt.Errorf("%w: %s", ErrorLabDoesntExist, name)
	}

	return labs[name], nil
}

// GetAll returns all the labs that have been provisioned
func GetAll() []*lab {
	labsBuffer := make([]*lab, 0)

	for _, lab := range labs {
		if lab != nil {
			labsBuffer = append(labsBuffer, lab)
		}
	}

	return labsBuffer
}

func GetAllStarted() []*lab {
	labs := make([]*lab, 0)

	for _, lab := range labs {
		if lab.isStarted {
			labs = append(labs, lab)
		}
	}

	return labs
}

func HasCapacity(path string) (bool, error) {
	labDTO, err := load(path)
	if err != nil {
		return false, nil
	}

	totalMemoryRequired := 2 * 1024 // 2GB for the frontend

	for _, labChallenge := range labDTO.Challenges {
		storedChallenge, err := store.GetChallenge(labChallenge.Challenge)
		if err != nil {
			return false, err
		}

		totalMemoryRequired += storedChallenge.Memory
	}

	availableMemory, err := utils.GetAvailableMemory()
	if err != nil {
		return false, err
	}

	return availableMemory >= totalMemoryRequired, nil
}

func GetCapacity() (int, error) {
	availableMemory, err := utils.GetAvailableMemory()
	if err != nil {
		return 0, err
	}

	return availableMemory, nil
}

func Provision(path string) (lab, error) {
	labDTO, err := load(path)
	if err != nil {
		return lab{}, err
	}

	if labDTO.Name == "" {
		return lab{}, ErrorLabNameMissing
	}

	if len(labDTO.Challenges) == 0 {
		return lab{}, ErrorLabNoChallenges
	}

	network, err := networks.Provision(networks.ProvisionNetworkOptions{})
	if err != nil {
		return lab{}, err
	}

	// Challenges that are provisioned but not started
	challenges := make([]*challenge.Challenge, 0)
	for _, labChallenge := range labDTO.Challenges {
		if labChallenge.Name == "" {
			return lab{}, ErrorChallengeNameEmpty
		}

		if labChallenge.Challenge == "" {
			return lab{}, ErrorChallengeIDEmpty
		}

		if len(labChallenge.Dns) == 0 {
			return lab{}, ErrorChallengeNoDNS
		}

		storedChallenge, err := store.GetChallenge(labChallenge.Challenge)
		if err != nil {
			return lab{}, err
		}

		newChallenge, err := challenge.Provision(&challenge.ProvisionChallengeOptions{
			Image:       storedChallenge.Image,
			DNSServers:  []string{network.GetDNSAddr()},
			DNSSettings: labChallenge.Dns,
		})
		if err != nil {
			return lab{}, err
		}

		challenges = append(challenges, newChallenge)
	}

	frontend, err := frontend.Provision(&frontend.ProvisionFrontendOptions{
		Network: network.GetName(),
		DNS:     []string{network.GetDNSAddr()},
		Port:    4000,
	})
	if err != nil {
		return lab{}, err
	}

	thisLab := lab{
		name:       labDTO.Name,
		challenges: challenges,
		network:    network,
		frontend:   frontend,
	}

	labs[labDTO.Name] = &thisLab

	return thisLab, nil
}

func RemoveAll() error {
	for _, lab := range labs {
		if err := lab.Remove(); err != nil {
			return err
		}
	}

	return nil
}

// Start starts the lab by starting all the challenges and connecting them to the isolated network
func (l *lab) Start() error {
	log.Debug().Msgf("Starting lab %s, data: %v", l.name, l)

	if err := l.network.Create(); err != nil {
		return err
	}

	zoneFileEntries := make([]dns.ZoneFileEntry, 0)

	for _, challenge := range l.challenges {

		// start the challenge, stop if an error occurs and remove all challenges
		if err := challenge.Start(); err != nil {
			for _, challenge := range l.challenges {
				challenge.Remove()
			}

			return err
		}

		// Connect the newly started challenge to the isolated network
		if err := l.network.Connect(challenge); err != nil {

			// If a challenge could not be connected, lets just remove it and continue
			challenge.Remove()
			return err
		}

		for _, hostname := range challenge.GetDNS() {
			zoneFileEntries = append(zoneFileEntries, dns.ZoneFileEntry{
				Hostname: hostname,
				Ip:       challenge.GetIP(),
			})
		}
	}

	// Provision and start the DNS service
	dnsService, err := dns.Provision(&dns.ProvisionDNSOptions{
		ZoneFileEntries: zoneFileEntries,
	})
	if err != nil {
		return err
	}

	if err = l.network.ConnectDNS(dnsService.GetContainer()); err != nil {
		return err
	}

	dhcpService, err := dhcp.Provision(&dhcp.ProvisionDHCPOptions{
		DNSAddr:     l.network.GetDNSAddr(),
		Subnet:      l.network.GetSubnet(),
		NetworkMode: l.network.GetName(),
	})
	if err != nil {
		return err
	}

	if err = l.network.ConnectDHCP(dhcpService.GetContainer()); err != nil {
		return err
	}

	if err = dhcpService.Start(); err != nil {
		return err
	}

	if err = dnsService.Start(); err != nil {
		return err
	}

	if err = l.frontend.Start(); err != nil {
		return err
	}

	cip, err := networks.ConnectToBridge(l.frontend.GetContainer())
	if err != nil {
		return err
	}

	log.Debug().Msgf("Frontend connected to bridge network, IP: %s", cip)

	l.dhcpService = dhcpService
	l.dnsService = dnsService

	l.isStarted = true
	labs[l.name] = l

	return nil
}

// Remove removes the lab by removing all the challenges and then the isolated network
func (l *lab) Remove() error {
	// Remove all the challenge containers
	for _, challenge := range l.challenges {
		if err := challenge.Remove(); err != nil {
			return err
		}
	}

	// Remove the DHCP and DNS service
	if err := l.dhcpService.Stop(); err != nil {
		return err
	}

	if err := l.dnsService.Stop(); err != nil {
		return err
	}

	if err := l.network.Remove(); err != nil {
		return err
	}

	delete(labs, l.name)

	return nil
}

func (l *lab) GetName() string {
	return l.name
}

func (l *lab) GetChallenges() []*challenge.Challenge {
	return l.challenges
}

func load(path string) (labDTO, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return labDTO{}, err
	}

	var lab labDTO
	if err = yaml.Unmarshal(data, &lab); err != nil {
		return labDTO{}, err
	}

	return lab, nil
}
