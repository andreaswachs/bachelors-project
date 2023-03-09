package store

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var (
	store *Store

	ErrorChallengeNotFound       = fmt.Errorf("challenge not found")
	ErrorChallengeNameEmpty      = fmt.Errorf("challenge name is empty")
	ErrorChallengeIdEmpty        = fmt.Errorf("challenge id is empty")
	ErrorChallengeImageEmpty     = fmt.Errorf("challenge image is empty")
	ErrorChallengeMemoryEmpty    = fmt.Errorf("challenge memory is empty")
	ErrorChallengeNegativeMemory = fmt.Errorf("challenge memory is negative")
)

type ChallengeTemplate struct {
	Name   string `yaml:"name"`
	Id     string `yaml:"id"`
	Image  string `yaml:"image"`
	Memory int64  `yaml:"memory"`
}

type Store struct {
	challenges map[string]ChallengeTemplate
}

type storeDTO struct {
	Challenges []ChallengeTemplate `yaml:"challenges"`
}

func init() {
	store = &Store{
		challenges: make(map[string]ChallengeTemplate),
	}

	if err := Load("store.yaml"); err != nil {
		log.Err(err).Msg("failed to load store from default path \"store.yaml\"")
	}
}

func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	dto, err := loadStore(data)
	if err != nil {
		return err
	}

	if err = validateStoreDTO(dto); err != nil {
		return err
	}

	transferChallenges(dto)

	return nil
}

func GetChallenge(name string) (ChallengeTemplate, error) {
	if _, ok := store.challenges[name]; !ok {
		return ChallengeTemplate{}, ErrorChallengeNotFound
	}

	return store.challenges[name], nil
}

func ChallengeExists(name string) bool {
	if _, ok := store.challenges[name]; !ok {
		return false
	}

	return true
}

func loadStore(data []byte) (*storeDTO, error) {
	var dto storeDTO

	err := yaml.Unmarshal(data, &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

func transferChallenges(dto *storeDTO) {
	for _, c := range dto.Challenges {
		store.challenges[c.Name] = c
	}
}

func validateStoreDTO(dto *storeDTO) error {
	for _, c := range dto.Challenges {
		if c.Name == "" {
			return ErrorChallengeNameEmpty
		}

		if c.Id == "" {
			return ErrorChallengeIdEmpty
		}

		if c.Image == "" {
			return ErrorChallengeImageEmpty
		}

		if c.Memory == 0 {
			return ErrorChallengeMemoryEmpty
		}

		if c.Memory < 0 {
			return ErrorChallengeNegativeMemory
		}
	}

	return nil
}
