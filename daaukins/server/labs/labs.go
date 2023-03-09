package labs

import (
	"fmt"
	"os"

	"github.com/andreaswachs/bachelors-project/daaukins/server/challenge"
	"github.com/andreaswachs/bachelors-project/daaukins/server/networks"
	"github.com/andreaswachs/bachelors-project/daaukins/server/store"
	docker "github.com/fsouza/go-dockerclient"
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

// lab is the data bearing struct for a lab
type lab struct {
	name        string
	challenges  []*challenge.Challenge
	dhcpService *networkService
	dnsService  *networkService
	network     *networks.Network
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

type networkService struct {
	container      *docker.Container
	filesToCleanup []string
}

func init() {
	// Ensure labs is initiated to an empty slice
	labs = make(map[string]*lab, 0)
}

// WithName returns a lab with a given name. If the lab does not exist, an error is returned
func WithName(name string) (*lab, error) {
	if labs[name] == nil {
		return nil, fmt.Errorf("%w: %s", ErrorLabDoesntExist, name)
	}

	return labs[name], nil
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

	zoneFileEntries := make([]zoneFileEntry, 0)

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

		for _, dns := range labChallenge.Dns {
			zoneFileEntries = append(zoneFileEntries, zoneFileEntry{
				hostname: dns,
				ip:       "", // TODO
			})
		}

		challenges = append(challenges, newChallenge)
	}

	thisLab := lab{
		name:       labDTO.Name,
		challenges: challenges,
		network:    network,
	}

	labs[labDTO.Name] = &thisLab

	return thisLab, nil
}

// Start starts the lab by starting all the challenges and connecting them to the isolated network
func (l *lab) Start() error {
	if err := l.network.Create(); err != nil {
		return err
	}

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
	}

	return nil
}

// Remove removes the lab by removing all the challenges and then the isolated network
func (l *lab) Remove() error {
	for _, challenge := range l.challenges {
		if err := challenge.Remove(); err != nil {
			return err
		}
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
