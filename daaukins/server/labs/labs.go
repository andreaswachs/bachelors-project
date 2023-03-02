package labs

import (
	"fmt"
	"os"

	"github.com/andreaswachs/bachelors-project/daaukins/server/challenge"
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

func NewLab(path string) (*Lab, error) {
	lab := &Lab{
		labs: make(map[string]*lab),
	}

	if err := lab.load(path); err != nil {
		return nil, err
	}

	return lab, nil
}

func (l *Lab) load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	dto, err := loadLab(data)
	if err != nil {
		return err
	}

	if err = validateLabDTO(dto); err != nil {
		return err
	}

	l.transferLab(dto)

	return nil
}

func loadLab(data []byte) (*labDTO, error) {
	dto := &labDTO{}
	if err := yaml.Unmarshal(data, dto); err != nil {
		return nil, err
	}

	return dto, nil
}

func validateLabDTO(dto *labDTO) error {
	if dto.Name == "" {
		return ErrorLabNameMissing
	}

	if len(dto.Challenges) == 0 {
		return ErrorLabNoChallenges
	}

	for _, challenge := range dto.Challenges {
		if err := validateChallenge(&challenge); err != nil {
			return err
		}
	}

	return nil
}

func validateChallenge(c *labChallenge) error {
	if c.Name == "" {
		return ErrorChallengeNameEmpty
	}

	if c.Challenge == "" {
		return ErrorChallengeIDEmpty
	}

	if len(c.Dns) == 0 {
		return ErrorChallengeNoDNS
	}

	return nil
}

func transferLab(dto *labDTO) (*lab, error) {
	lab := &lab{
		name:       dto.Name,
		challenges: make([]*challenge.Challenge, 0),
	}

	for _, dtoChallenge := range dto.Challenges {
		c, err := transferChallenge(&dtoChallenge)
		if err != nil {
			return nil, err
		}
		lab.challenges = append(lab.challenges, c)
	}

	return lab, nil
}

func transferChallenge(dto *labChallenge) (*challenge.Challenge, error) {
	storedChallenge, err := store.WithStore().GetChallenge(dto.Name)
	if err != nil {
		return nil, err
	}

	return challenge.Provision(nil, &challenge.ProvisionChallengeOptions{
		DNSServers: "TODO", // this points to the CoreDNS service
		Image:      storedChallenge.Image,
	})

}
