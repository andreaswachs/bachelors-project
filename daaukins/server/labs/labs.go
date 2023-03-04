package labs

import (
	"fmt"
	"os"

	"github.com/andreaswachs/bachelors-project/daaukins/server/challenge"
	"github.com/andreaswachs/bachelors-project/daaukins/server/networks"
	"github.com/andreaswachs/bachelors-project/daaukins/server/store"
	"gopkg.in/yaml.v3"
)

var (
	ErrorLabNameMissing     = fmt.Errorf("lab name is missing")
	ErrorLabNoChallenges    = fmt.Errorf("lab has no challenges")
	ErrorChallengeNameEmpty = fmt.Errorf("challenge name is empty")
	ErrorChallengeIDEmpty   = fmt.Errorf("challenge id is empty")
	ErrorChallengeNoDNS     = fmt.Errorf("challenge has no dns servers")
)

// Lab is the front facing data structure that is used to interact with labs
type Lab struct {
	labs map[string]*lab
}

// lab is the data bearing struct for a lab
type lab struct {
	name       string
	challenges []*challenge.Challenge
	network    *networks.Network
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

	network, err := networks.Provision(nil, networks.ProvisionNetworkOptions{})
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

		storedChallenge, err := store.WithStore().GetChallenge(labChallenge.Challenge)
		if err != nil {
			return lab{}, err
		}

		// TODO: insert docker client, DNSServers
		newChallenge, err := challenge.Provision(nil, &challenge.ProvisionChallengeOptions{
			Image:       storedChallenge.Image,
			DNSServers:  nil,
			DNSSettings: labChallenge.Dns,
		})
		if err != nil {
			return lab{}, err
		}

		challenges = append(challenges, newChallenge)
	}

	return lab{
		name:       labDTO.Name,
		challenges: challenges,
		network:    network,
	}, nil
}

func (l *Lab) Start() error {
	// TODO
	return nil
}

func load(path string) (labDTO, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return labDTO{}, err
	}

	var lab labDTO
	err = yaml.Unmarshal(data, &lab)
	if err != nil {
		return labDTO{}, err
	}

	return lab, nil
}
